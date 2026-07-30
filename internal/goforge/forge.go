package goforge

import (
	"fmt"

	"github.com/andrearcaina/goforge/internal/config"
	"github.com/andrearcaina/goforge/internal/ui"
	"github.com/andrearcaina/goforge/internal/utils"
	"github.com/charmbracelet/huh/spinner"
)

func Forge(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("invalid configuration: configuration cannot be nil")
	}

	if cfg.Default {
		cfg.Form = config.Form{
			Name:           "example-server",
			ServerTypeFlag: config.REST,
			DatabaseFlag:   false,
		}
	} else {
		if err := ui.Run(cfg); err != nil {
			return fmt.Errorf("failed to run UI: %w", err)
		}
	}

	if err := utils.Validate(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	var err error
	_ = spinner.New().
		Title("Forging your project...").
		Action(func() {
			err = Generate(cfg)
		}).Run()

	if err != nil {
		return fmt.Errorf("failed to generate the file: %w", err)
	}

	ui.OutputSuccess(cfg)

	return nil
}
