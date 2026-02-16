package ui

import (
	"fmt"

	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/charmbracelet/huh"
)

func Run(cfg *spec.Config) error {
	// create the form
	form := createForm(cfg)
	if form == nil {
		// end early, no need to run the form (all flags are set)
		// so just return original cfg and nil error
		return nil
	}

	// run the form
	if err := form.Run(); err != nil {
		return err
	}

	// return nil, cfg is updated by reference if Run() succeeds
	return nil
}

func createForm(cfg *spec.Config) *huh.Form {
	var groups []*huh.Group

	// check if any flag is empty

	if cfg.Form.Name == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("What should be the name of the project/go module?").
				Value(&cfg.Form.Name).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("project name cannot be empty")
					}
					return nil
				}),
		))
	}

	if cfg.Form.ServerTypeFlag == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[spec.ServerTypeFlag]().
				Title("What should be the type of server?").
				Options(
					huh.NewOption("REST", spec.REST),
					huh.NewOption("gRPC", spec.GRPC),
					huh.NewOption("GraphQL", spec.GraphQL),
				).
				Value(&cfg.Form.ServerTypeFlag),
		))
	}

	if !cfg.Form.DatabaseFlag {
		groups = append(groups, huh.NewGroup(
			huh.NewConfirm().
				Title("Should I generate database files?").
				Value(&cfg.Form.DatabaseFlag),
		))
	}

	if !cfg.Form.MakefileFlag {
		groups = append(groups, huh.NewGroup(
			huh.NewConfirm().
				Title("Should I generate Makefile?").
				Value(&cfg.Form.MakefileFlag),
		))
	}

	// if no groups were added, return nil
	if len(groups) == 0 {
		return nil
	}

	form := huh.NewForm(groups...)

	return form
}
