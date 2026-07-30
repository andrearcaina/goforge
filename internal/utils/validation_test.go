package utils

import (
	"strings"
	"testing"

	"github.com/andrearcaina/goforge/internal/config"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr string
	}{
		{
			name: "valid REST configuration",
			cfg: &config.Config{
				OutputPath: "service",
				Form: config.Form{
					Name:           "example-service",
					ServerTypeFlag: config.REST,
				},
			},
		},
		{
			name:    "nil configuration",
			wantErr: "configuration cannot be nil",
		},
		{
			name: "empty output path",
			cfg: &config.Config{
				OutputPath: " ",
				Form: config.Form{
					Name:           "example",
					ServerTypeFlag: config.REST,
				},
			},
			wantErr: "output path cannot be empty",
		},
		{
			name: "empty project name",
			cfg: &config.Config{
				OutputPath: "service",
				Form: config.Form{
					ServerTypeFlag: config.REST,
				},
			},
			wantErr: "project name cannot be empty",
		},
		{
			name: "invalid project name",
			cfg: &config.Config{
				OutputPath: "service",
				Form: config.Form{
					Name:           "invalid name",
					ServerTypeFlag: config.REST,
				},
			},
			wantErr: "project name can only contain",
		},
		{
			name: "unsupported server",
			cfg: &config.Config{
				OutputPath: "service",
				Form: config.Form{
					Name:           "example",
					ServerTypeFlag: "invalid",
				},
			},
			wantErr: "unsupported server type",
		},
		{
			name: "docker without database",
			cfg: &config.Config{
				OutputPath: "service",
				Form: config.Form{
					Name:           "example",
					ServerTypeFlag: config.REST,
					DockerFlag:     true,
				},
			},
			wantErr: "docker generation requires database generation",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(test.cfg)
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() returned unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("Validate() error = %v, want error containing %q", err, test.wantErr)
			}
		})
	}
}
