package ui

import (
	"testing"

	"github.com/andrearcaina/goforge/internal/config"
)

func TestCreateFormSkipsExplicitFalseFlags(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		DatabaseFlagSet: true,
		MakefileFlagSet: true,
		DockerFlagSet:   true,
		Form: config.Form{
			Name:           "example-service",
			ServerTypeFlag: config.REST,
			DatabaseFlag:   false,
			MakefileFlag:   false,
			DockerFlag:     false,
		},
	}

	if form := createForm(cfg); form != nil {
		t.Fatal("createForm() returned a form even though every option was explicitly supplied")
	}
}

func TestCreateFormDoesNotPromptForDockerWithoutDatabase(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		DatabaseFlagSet: true,
		MakefileFlagSet: true,
		Form: config.Form{
			Name:           "example-service",
			ServerTypeFlag: config.REST,
			DatabaseFlag:   false,
			MakefileFlag:   false,
		},
	}

	if form := createForm(cfg); form != nil {
		t.Fatal("createForm() prompted for Docker even though database generation was explicitly disabled")
	}
}
