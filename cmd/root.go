/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "web_spider",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	inputFile string
	ws        bool
	inputDir  string
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	rootCmd.PersistentFlags().StringVarP(&inputDir, "input_dir", "d", "", "Input dir path")
	rootCmd.PersistentFlags().BoolVarP(&ws, "ws", "", false, "use websocket")
}
