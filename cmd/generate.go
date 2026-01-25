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
		if err := goforge.Generate(cfg); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// first flag is just for testing (will be removed later)
	generateCmd.Flags().StringVar(&cfg.SomeFlag, "some-flag", "World", "A flag to specify who to say hello to")
	generateCmd.Flags().StringVarP(&cfg.OutputPath, "path", "p", ".", "The directory to write the generated file to")
}
