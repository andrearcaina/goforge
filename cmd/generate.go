/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/andrearcaina/goforge/internal/config"
	"github.com/andrearcaina/goforge/internal/goforge"
	"github.com/spf13/cobra"
)

var cfg config.Config

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
	generateCmd.Flags().BoolVarP(&cfg.Form.MakefileFlag, "makefile", "m", false, "Generate makefile (if flag is set, set to true)")
	generateCmd.Flags().BoolVarP(&cfg.Form.DockerFlag, "docker", "D", false, "Generate Docker Compose (for DB) and .env file (if flag is set, set to true)")

	// normalize server type flag to lowercase
	generateCmd.PreRun = func(cmd *cobra.Command, args []string) {
		// first convert to string, then lowercase the string, then convert back to ServerTypeFlag
		cfg.Form.ServerTypeFlag = config.ServerTypeFlag(strings.ToLower(string(cfg.Form.ServerTypeFlag)))

		// validate go.mod names
		if cfg.Form.Name != "" {
			if matched, _ := regexp.MatchString(`[^a-zA-Z0-9_\-]`, cfg.Form.Name); matched {
				fmt.Println("Error: project name can only contain letters, numbers, underscores, or dashes")
				os.Exit(1)
			}
		}
	}
}
