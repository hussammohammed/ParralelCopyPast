package fileTransfer

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type FileTransferHandler struct{}
type fileChunk struct {
	num int
	//data []byte
	//msg  string
}
type Result struct {
	chunkNum   int
	message    string
	error      error
	statusCode int
}

const (
	chunkSize = 1024 * 1024 // 1MB is the chunk size
)

func (fth *FileTransferHandler) CopyPastFile(sourcePath string, destinationPath string) error {
	fth.splitFileToChunks(sourcePath, destinationPath)
	return nil
}

func (fth *FileTransferHandler) splitFileToChunks(sourcePath string, destinationPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		log.Println("can not open file")
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("can not read file info")
	}
	fileSize := fileInfo.Size()
	log.Printf("file size is %v", fileSize)
	sendChunksChannelAsync(sourcePath, destinationPath, fileSize)
	return nil
}

func sendChunksChannelAsync(sourcePath string, destinationPath string, fileSize int64) {
	totalChunksCount := fileSize / chunkSize
	log.Printf("total chunks are %v\n", totalChunksCount)
	done := make(chan interface{})
	chunks := make(chan fileChunk, totalChunksCount) // Use buffered channel with a capacity of chunks count
	defer close(done)
	for i := 0; i < 10; i++ {
		chunks <- fileChunk{num: i}
	}
	close(chunks)
	for result := range receiveChunksChannel(done, chunks, destinationPath) {
		log.Printf("result of chunk #%v is: %v\n", result.chunkNum, result.statusCode)
	}
	//<-terminated
	// close(chunks)
	log.Println("all file chunks have been sent")
}

func receiveChunksChannel(done <-chan interface{}, chunks <-chan fileChunk, destinationPath string) <-chan Result {
	fmt.Println(len(chunks))
	terminated := make(chan Result)
	go func() {
		defer log.Println("all file chunks have been received")
		defer close(terminated)
		for {
			select {
			case chunk, ok := <-chunks:
				if !ok {
					return
				}
				log.Printf("i received chunk #%v\n", chunk.num)
				// this just for initial testing of communication between channels
				if chunk.num == 1 {
					terminated <- Result{chunkNum: chunk.num, statusCode: http.StatusInternalServerError, message: fmt.Sprintf("can not write chunk #%v", chunk.num)}
				} else {
					terminated <- Result{chunkNum: chunk.num, statusCode: http.StatusOK}
				}

			case <-done:
				return
			}
		}
	}()
	return terminated
}