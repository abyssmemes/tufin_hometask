package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Entry point of the application
func main() {
	var rootCmd = &cobra.Command{Use: "tufin"}

	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(statusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
