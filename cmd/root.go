package cmd

import (
	"fmt"
	"os"

	"github.com/MinoKim82/go-kis-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kis",
	Short: "Korea Investment & Securities (한국투자증권) CLI",
	Long:  `A command line interface for Korea Investment & Securities REST API to query quotes, check balances, and execute trades.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.kis-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&config.EnvName, "env", "mock", "environment to use: mock or prod (default is mock)")
}
