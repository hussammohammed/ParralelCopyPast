package fileTransfer

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileTransferHandler struct{}
type fileChunk struct {
	fileName string
	num      int
	data     []byte
	//msg  string
}
type Result struct {
	chunkNum   int
	message    string
	error      error
	statusCode int
}

const (
	chunkSize = 5 * 1024 * 1024 // 5MB is the chunk size
)

func (fth *FileTransferHandler) CopyPast(sourcePath string, destinationPath string) error {
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	if fileInfo.IsDir() {
		directoryName := getLastFolderName(sourcePath)
		err := CopyPastDirectory(directoryName, sourcePath, destinationPath)
		return err
	} else {
		err := sendFileAsChunks(sourcePath, destinationPath)
		return err
	}
}
func CopyPastDirectory(directoryName string, baseSourcePath string, destinationPath string) error {
	files, err := os.ReadDir(baseSourcePath)
	if err != nil {
		log.Fatal(err)
	}
	//create new folder at distination path with same name as source path
	destinationPath = filepath.Join(destinationPath, directoryName)
	os.Mkdir(destinationPath, fs.ModePerm)
	begin := make(chan interface{})
	var wg sync.WaitGroup
	// loop over all files and folders inside main folder
	for _, file := range files {
		if err != nil {
			log.Fatal(err)
		}
		curFileName := file.Name()
		sourcePath := baseSourcePath
		// if current file is a folder then call CopyPastDirectory again to copy nest folders and files
		if file.IsDir() {
			sourcePath = filepath.Join(baseSourcePath, curFileName)
			CopyPastDirectory(curFileName, sourcePath, destinationPath)
		} else { // copy/past files inparallel (separated goroute) then wait all goroutens to finish
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-begin // Here the goroutine waits until it is told it can continue.
				fileSourcePath := filepath.Join(sourcePath, fmt.Sprintf("/%v", curFileName))
				fileDestinationPath := filepath.Join(destinationPath, fmt.Sprintf("/%v", curFileName))
				sendFileAsChunks(fileSourcePath, fileDestinationPath)
			}()

		}
	}
	close(begin) // Here we close the channel, thus unblocking all the goroutines simultaneously
	wg.Wait()
	return nil
}
func getLastFolderName(path string) string {
	// Clean the path to ensure consistent separators
	cleanPath := filepath.Clean(path)
	// Split the cleaned path into individual components
	components := strings.Split(cleanPath, string(filepath.Separator))
	// Get the last component, which represents the last folder name
	lastComponent := components[len(components)-1]
	return lastComponent
}

// this functions split file into bunch of chucks and send those chuncks via channel to another function which listen to channel and combine chuncks again to new file and destination path
func sendFileAsChunks(sourcePath string, destinationPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		log.Println("can not open file")
		return err
	}
	// get file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("can not read file info")
	}
	defer file.Close()

	fileSize := fileInfo.Size()
	totalChunksCount := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		totalChunksCount++
	}
	done := make(chan interface{})
	chunks := make(chan fileChunk, totalChunksCount) // Use buffered channel with a capacity of chunks count
	defer close(done)

	buffer := make([]byte, chunkSize)
	result := receiveChunksChannel(done, chunks, destinationPath)
	for i := 0; i < int(totalChunksCount); i++ {
		//chunkData, err := file.Read(buffer)
		readBytes, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("end of the file")
				break
			}
			return err
		}
		if readBytes > 0 {
			// ensures that chunk contains only the actual data read from the file, without any uninitialized or unused bytes.
			chunk := make([]byte, readBytes)
			copy(chunk, buffer[:readBytes])
			chunks <- fileChunk{fileName: fileInfo.Name(), num: i, data: chunk}
		}

	}
	// after sending all file chunks, then close chunks channel to allow receiveChunksChannel function to return
	close(chunks)

	for {
		select {
		case r, ok := <-result:
			if !ok {
				log.Println("all file chunks have been sent")
				return nil
			}
			log.Printf("result of chunk #%v is: %v\n", r.chunkNum, r.statusCode)
			// if r.statusCode == http.StatusInternalServerError {
			// 	return
			// }
		}
	}
}

func receiveChunksChannel(done <-chan interface{}, chunks <-chan fileChunk, destinationPath string) <-chan Result {
	combinedFile, err := os.Create(destinationPath)
	if err != nil {
		log.Fatal(err)
	}

	terminated := make(chan Result)
	go func() {
		defer log.Println("all file chunks have been received")
		defer combinedFile.Close()
		defer close(terminated)
		for {
			select {
			case <-done:
				return

			case chunk, ok := <-chunks:
				if !ok {
					return
				}
				//log.Printf("i received chunk %v #%v size: %v\n", chunk.fileName, chunk.num, len(chunk.data))
				_, err := combinedFile.Write(chunk.data)
				if err != nil {
					terminated <- Result{chunkNum: chunk.num, statusCode: http.StatusInternalServerError, message: fmt.Sprintf("can not write chunk #%v", chunk.num)}
					log.Fatal(err)
				}
			}
		}
	}()
	return terminated
}
