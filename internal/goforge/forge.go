package goforge

import (
	"errors"
	"fmt"

	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/andrearcaina/goforge/internal/ui"
)

func Forge(initial *spec.Config) error {
	// try to run the ui
	cfg, err := ui.Run(initial)
	if err != nil {
		return errors.New("ui failed")
	}

	// try to generate the file
	if err := Generate(cfg); err != nil {
		return errors.New("generating file failed")
	}

	fmt.Println("File generated successfully")

	// success if none fails
	return nil
}
