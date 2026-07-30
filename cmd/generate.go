package cmd

import (
	"strings"

	"github.com/andrearcaina/goforge/internal/config"
	"github.com/andrearcaina/goforge/internal/goforge"
	"github.com/spf13/cobra"
)

var cfg config.Config

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Go backend service",
	Long: `Generate a REST, gRPC, or GraphQL Go backend service.

Missing options are collected through an interactive form.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return goforge.Forge(&cfg)
	},
}

func init() {
	// general flags
	generateCmd.Flags().StringVarP(&cfg.OutputPath, "path", "p", ".", "The directory to write the generated file to")
	generateCmd.Flags().BoolVarP(&cfg.Default, "default", "d", false, "Use default configuration values")
	generateCmd.Flags().BoolVarP(&cfg.Force, "force", "f", false, "Overwrite generated files that already exist")

	// flags for form fields
	generateCmd.Flags().StringVarP(&cfg.Form.Name, "name", "n", "", "The name for the go.mod module")
	generateCmd.Flags().StringVarP((*string)(&cfg.Form.ServerTypeFlag), "server", "s", "", "Type of server to generate (rest/grpc/graphql)")
	generateCmd.Flags().BoolVar(&cfg.Form.DatabaseFlag, "database", false, "Generate database files (if flag is set, set to true)")
	generateCmd.Flags().BoolVarP(&cfg.Form.MakefileFlag, "makefile", "m", false, "Generate makefile (if flag is set, set to true)")
	generateCmd.Flags().BoolVarP(&cfg.Form.DockerFlag, "docker", "D", false, "Generate Docker Compose (for DB) and .env file (if flag is set, set to true)")

	// normalize server type flag to lowercase
	generateCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		cfg.Form.ServerTypeFlag = config.ServerTypeFlag(strings.ToLower(strings.TrimSpace(string(cfg.Form.ServerTypeFlag))))
		cfg.Form.Name = strings.TrimSpace(cfg.Form.Name)
		cfg.OutputPath = strings.TrimSpace(cfg.OutputPath)
		cfg.DatabaseFlagSet = cmd.Flags().Changed("database")
		cfg.MakefileFlagSet = cmd.Flags().Changed("makefile")
		cfg.DockerFlagSet = cmd.Flags().Changed("docker")
		return nil
	}
}
