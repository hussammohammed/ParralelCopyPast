package main

import (
	"fmt"

	fileTransfer "github.com/hussammohammed/ParallelCopyPast/filetransfer"
)

func main() {
	handler := fileTransfer.FileTransferHandler{}
	sourcePath := "/media/hossam/Work/hossam/cources/copy past test"
	destinationPath := "/media/hossam/Work"
	err := handler.CopyPast(sourcePath, destinationPath)
	if err != nil {
		fmt.Println(err.Error())
	}
}
