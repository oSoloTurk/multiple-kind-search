package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Multiple Kind Search Application",
	Long:  `A search application that demonstrates multiple kind search capabilities using Elasticsearch.`,
}

func Execute() {
	fmt.Printf("Available commands: %v\n", rootCmd.Commands())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
