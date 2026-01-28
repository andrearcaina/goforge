package goforge

import (
	"errors"
	"fmt"

	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/andrearcaina/goforge/internal/ui"
	"github.com/charmbracelet/huh/spinner"
)

func Forge(cfg *spec.Config) error {
	if cfg.Default {
		cfg.Form = spec.Form{
			Name:           "example-server",
			ServerTypeFlag: "rest",
		}
	} else {
		if err := ui.Run(cfg); err != nil {
			return errors.New("failed to run ui form")
		}
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
