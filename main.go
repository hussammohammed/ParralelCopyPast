package main

import (
	"fmt"
	"path/filepath"

	fileTransfer "github.com/hussammohammed/ParallelCopyPast/filetransfer"
)

func main() {

	fmt.Println("init project")
	handler := fileTransfer.FileTransferHandler{}
	sourcePath := filepath.Join("/media/hossam/Work/hossam/cources/go", "Golang Web Server and RSS Scraper _ Full Tutorial.mp4")
	destinationPath := "/media/hossam/Work"
	err := handler.CopyPastFile(sourcePath, destinationPath)
	if err != nil {
		fmt.Println(err.Error())
	}
}
