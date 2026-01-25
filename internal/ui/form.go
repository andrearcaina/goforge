package ui

import (
	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/charmbracelet/huh"
)

func Run(defaults *spec.Config) (*spec.Config, error) {
	cfg := defaults

	// create the form
	form := createForm(cfg)
	if form == nil {
		// end early, no need to run the form (all flags are set)
		// so just return original cfg and nil error
		return cfg, nil
	}

	// run the form
	if err := form.Run(); err != nil {
		return nil, err
	}

	// return the cfg with updated values
	return cfg, nil
}

func createForm(cfg *spec.Config) *huh.Form {
	var groups []*huh.Group

	// check if any flag is empty

	if cfg.Form.SomeFlag == "" {
		groups = appendGroup(groups, "Who should we say hello to?", &cfg.Form.SomeFlag)
	}

	if cfg.Form.AnotherFlag == "" {
		groups = appendGroup(groups, "Another flag", &cfg.Form.AnotherFlag)
	}

	// if no groups were added, return nil
	if len(groups) == 0 {
		return nil
	}

	form := huh.NewForm(groups...)

	return form
}

func appendGroup(groups []*huh.Group, title string, value *string) []*huh.Group {
	return append(groups, huh.NewGroup(
		huh.NewInput().
			Title(title).
			Value(value),
	))
}
