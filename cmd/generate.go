/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/andrearcaina/goforge/internal/goforge"
	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/spf13/cobra"
)

var cfg spec.Config

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a hello world message",
	Long: `Generates a simple hello world message using the provided flag.
For example:

goforge generate --some-flag "Developer"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := goforge.Forge(&cfg); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// general flags
	generateCmd.Flags().StringVarP(&cfg.OutputPath, "path", "p", ".", "The directory to write the generated file to")
	generateCmd.Flags().BoolVarP(&cfg.Default, "default", "d", false, "Use default configuration values")

	// flags for form fields
	generateCmd.Flags().StringVarP(&cfg.Form.Name, "name", "n", "", "The name for the go.mod module")
	generateCmd.Flags().StringVarP((*string)(&cfg.Form.ServerTypeFlag), "server", "s", "", "Type of server to generate (rest/grpc/graphql)")
	generateCmd.Flags().BoolVar(&cfg.Form.DatabaseFlag, "database", false, "Generate database files (if flag is set, set to true)")

	// normalize server type flag to lowercase
	generateCmd.PreRun = func(cmd *cobra.Command, args []string) {
		// first convert to string, then lowercase the string, then convert back to ServerTypeFlag
		cfg.Form.ServerTypeFlag = spec.ServerTypeFlag(strings.ToLower(string(cfg.Form.ServerTypeFlag)))
	}
}
