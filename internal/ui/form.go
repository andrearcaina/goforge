package ui

import (
	"github.com/andrearcaina/goforge/internal/config"
	"github.com/andrearcaina/goforge/internal/utils"
	"github.com/charmbracelet/huh"
)

func Run(cfg *config.Config) error {
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

func createForm(cfg *config.Config) *huh.Form {
	var groups []*huh.Group

	// check if any flag is empty

	if cfg.Form.Name == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("What should be the name of the project/go module?").
				Value(&cfg.Form.Name).
				Validate(utils.ValidateProjectName),
		))
	}

	if cfg.Form.ServerTypeFlag == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[config.ServerTypeFlag]().
				Title("What should be the type of server?").
				Options(
					huh.NewOption("REST", config.REST),
					huh.NewOption("gRPC", config.GRPC),
					huh.NewOption("GraphQL", config.GraphQL),
				).
				Value(&cfg.Form.ServerTypeFlag),
		))
	}

	if !cfg.DatabaseFlagSet {
		groups = append(groups, huh.NewGroup(
			huh.NewConfirm().
				Title("Should I generate database files?").
				Value(&cfg.Form.DatabaseFlag),
		))
	}

	if !cfg.MakefileFlagSet {
		groups = append(groups, huh.NewGroup(
			huh.NewConfirm().
				Title("Should I generate Makefile? (recommended)").
				Value(&cfg.Form.MakefileFlag),
		))
	}

	if !cfg.DockerFlagSet && (!cfg.DatabaseFlagSet || cfg.Form.DatabaseFlag) {
		dockerGroup := huh.NewGroup(
			huh.NewConfirm().
				Title("Should I generate Docker compose file and .env file?").
				Value(&cfg.Form.DockerFlag),
		).WithHideFunc(func() bool {
			return !cfg.Form.DatabaseFlag
		})
		groups = append(groups, dockerGroup)
	}

	// if no groups were added, return nil
	if len(groups) == 0 {
		return nil
	}

	form := huh.NewForm(groups...)

	return form
}
