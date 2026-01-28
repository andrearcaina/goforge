/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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

	generateCmd.Flags().StringVarP(&cfg.OutputPath, "path", "p", ".", "The directory to write the generated file to")
	generateCmd.Flags().BoolVarP(&cfg.Default, "default", "d", false, "Use default configuration values")

	generateCmd.Flags().StringVarP(&cfg.Form.Name, "name", "n", "", "The name for the go.mod module")
	generateCmd.Flags().StringVarP((*string)(&cfg.Form.ServerTypeFlag), "server", "s", "", "Type of server to generate (rest/grpc/graphql)")
	generateCmd.Flags().BoolVar(&cfg.Form.DatabaseFlag, "database", true, "Generate database files")
}
