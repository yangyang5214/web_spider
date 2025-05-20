package cmd

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
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
		chrome, chromeCancel, err := pkg.NewChromePool(logger, ws)
		if err != nil {
			panic(err)
		}
		if !ws {
			chromeCancel()
		}

		sd, err := sciencedirect.NewScienceDirect(chrome, logger)
		if err != nil {
			panic(err)
		}

		// 设置信号处理
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			sd.Close() // 收到终止信号时关闭
			os.Exit(0)
		}()

		err = sd.List()
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
		sd, err := sciencedirect.NewScienceDirect(nil, log.DefaultLogger)
		if err != nil {
			panic(err)
		}
		err = sd.Detail()
		if err != nil {
			panic(err)
		}
	},
}
