package cmd

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"web_spider/pkg"
	"web_spider/pkg/sciencedirect"
)

// sciencedirectCmd represents the sciencedirect command
var sciencedirectCmd = &cobra.Command{
	Use:   "sciencedirect",
	Short: "www.sciencedirect.com",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: sciencedirect [list|detail]")
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
		logger := log.DefaultLogger
		chrome, cacnel, err := pkg.NewChromePool(logger)
		if err != nil {
			panic(err)
		}
		defer cacnel()
		err = sciencedirect.NewScienceDirect(chrome, logger).List()
		if err != nil {
			panic(err)
		}
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
