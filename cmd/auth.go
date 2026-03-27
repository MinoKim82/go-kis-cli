package cmd

import (
	"fmt"
	"os"

	"github.com/MinoKim82/go-kis-cli/pkg/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Manage API tokens for Korea Investment & Securities REST API.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Issue and cache an OAuth token",
	Long:  `Issues a new OAuth token using the configured APPKEY and APPSECRET, and caches it locally (~/.kis-cli-token.json).`,
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IssueToken()
		if err != nil {
			fmt.Printf("Failed to issue token: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully logged in.")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
}
