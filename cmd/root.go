package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "Opinionated Go CLI for generating backend service boilerplate",
	Long: `goforge is an opinionated CLI tool for generating Go backend service boilerplate
quickly and consistently.

It supports multiple service styles, including REST, gRPC, and GraphQL, and generates
the required project structure, configuration, and dependencies based on the selected
protocol.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(generateCmd)
}
