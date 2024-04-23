package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "load_tester",
	Short: "A CLI to HTTP(s) load test your applications",
	Long: `Allows you to send multiple concurrent requests to a url and view the response
	results and statistics.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
