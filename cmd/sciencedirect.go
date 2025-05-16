package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sciencedirectCmd represents the sciencedirect command
var sciencedirectCmd = &cobra.Command{
	Use:   "sciencedirect",
	Short: "www.sciencedirect.com",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sciencedirect called")
	},
}

func init() {
	rootCmd.AddCommand(sciencedirectCmd)

	sciencedirectCmd.AddCommand(listCmd)
	sciencedirectCmd.AddCommand(detailCmd)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "search list page",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

// listCmd represents the list command
var detailCmd = &cobra.Command{
	Use:   "detail",
	Short: "single page",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}
