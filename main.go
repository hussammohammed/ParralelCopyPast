package main

import (
	"fmt"
	"path/filepath"

	fileTransfer "github.com/hussammohammed/ParallelCopyPast/filetransfer"
)

func main() {

	fmt.Println("init project")
	handler := fileTransfer.FileTransferHandler{}
	sourcePath := filepath.Join("/media/hossam/Work/hossam/cources/microservice", "NET Microservices – Full Course.mp4")
	destinationPath := filepath.Join("/media/hossam/Work", "NET Microservices – Full Course.mp4")
	err := handler.CopyPastFile(sourcePath, destinationPath)
	if err != nil {
		fmt.Println(err.Error())
	}
}
