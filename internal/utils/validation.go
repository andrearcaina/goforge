package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/andrearcaina/goforge/internal/config"
)

var projectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func Validate(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if strings.TrimSpace(cfg.OutputPath) == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if err := ValidateProjectName(cfg.Form.Name); err != nil {
		return err
	}

	switch cfg.Form.ServerTypeFlag {
	case config.REST, config.GRPC, config.GraphQL:
	default:
		return fmt.Errorf("unsupported server type %q (must be rest, grpc, or graphql)", cfg.Form.ServerTypeFlag)
	}

	if cfg.Form.DockerFlag && !cfg.Form.DatabaseFlag {
		return fmt.Errorf("docker generation requires database generation")
	}

	return nil
}

func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if !projectNamePattern.MatchString(name) {
		return fmt.Errorf("project name can only contain letters, numbers, underscores, or dashes")
	}
	if filepath.IsAbs(name) {
		return fmt.Errorf("project name must not be an absolute path")
	}
	return nil
}
