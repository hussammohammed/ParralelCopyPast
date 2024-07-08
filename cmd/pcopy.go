/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	fileTransfer "github.com/hussammohammed/ParallelCopyPast/filetransfer"
	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
)

// pcopyCmd represents the pcopy command
var pcopyCmd = &cobra.Command{
	Use:   "pcopy",
	Short: "parallel copy for folders and files",
	Long:  `this command to copy/past bunch of files and folders in parallel way via providing source and destination paths`,
	Run: func(cmd *cobra.Command, args []string) {
		if source != "" && destination != "" {
			handler := fileTransfer.FileTransferHandler{}
			handler.CopyPast(source, destination)
		}
	},
}

func init() {
	pcopyCmd.Flags().StringVarP(&source, "source", "s", "", "path of source")
	if err := pcopyCmd.MarkFlagRequired("source"); err != nil {
		fmt.Println(err)
	}

	pcopyCmd.Flags().StringVarP(&destination, "destination", "d", "", "path of destination")
	if err := pcopyCmd.MarkFlagRequired("destination"); err != nil {
		fmt.Println(err)
	}
	rootCmd.AddCommand(pcopyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pcopyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pcopyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
