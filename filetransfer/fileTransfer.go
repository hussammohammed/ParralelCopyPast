package fileTransfer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type FileTransferHandler struct{}
type fileChunk struct {
	num  int
	data []byte
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

func (fth *FileTransferHandler) CopyPastFile(sourcePath string, destinationPath string) error {
	readAndSendChunks(sourcePath, destinationPath)
	return nil
}

func readAndSendChunks(sourcePath string, destinationPath string) error {
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
	log.Printf("file size is %v", fileSize)
	totalChunksCount := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		totalChunksCount++
	}
	log.Printf("total chunks are %v\n", totalChunksCount)
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
			chunks <- fileChunk{num: i, data: chunk}
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
	fmt.Println(len(chunks))
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
				log.Printf("i received chunk #%v size: %v\n", chunk.num, len(chunk.data))
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
