package goforge

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/andrearcaina/goforge/internal/templates"
)

func Generate(cfg *spec.Config) error {
	if cfg.Form.ServerTypeFlag == "rest" {
		return generateRESTServer(cfg)
	}

	return errors.New("unknown server type")
}

func generateRESTServer(cfg *spec.Config) error {
	dirs := []string{
		"cmd/server",
		"internal/api",
		"internal/config",
	}

	if cfg.Form.DatabaseFlag {
		dirs = append(dirs, "internal/db")
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(cfg.OutputPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	if err := generateFile("rest/go.mod.tmpl", filepath.Join(cfg.OutputPath, "go.mod"), cfg); err != nil {
		return err
	}

	if err := generateFile("config.go.tmpl", filepath.Join(cfg.OutputPath, "internal/config/config.go"), cfg); err != nil {
		return err
	}
	if err := generateFile("main.go.tmpl", filepath.Join(cfg.OutputPath, "cmd/server/main.go"), cfg); err != nil {
		return err
	}

	if err := generateFile("rest/server.go.tmpl", filepath.Join(cfg.OutputPath, "internal/api/server.go"), cfg); err != nil {
		return err
	}

	if err := generateFile("rest/handler.go.tmpl", filepath.Join(cfg.OutputPath, "internal/api/handler.go"), cfg); err != nil {
		return err
	}

	if err := generateFile("rest/service.go.tmpl", filepath.Join(cfg.OutputPath, "internal/api/service.go"), cfg); err != nil {
		return err
	}

	if err := generateFile("rest/types.go.tmpl", filepath.Join(cfg.OutputPath, "internal/api/types.go"), cfg); err != nil {
		return err
	}

	if cfg.Form.DatabaseFlag {
		if err := os.MkdirAll(filepath.Join(cfg.OutputPath, "internal/db/generated"), 0755); err != nil {
			return err
		}

		if err := generateFile("db/db.go.tmpl", filepath.Join(cfg.OutputPath, "internal/db/db.go"), cfg); err != nil {
			return err
		}

		if err := generateFile("sqlc.yaml.tmpl", filepath.Join(cfg.OutputPath, "sqlc.yaml"), cfg); err != nil {
			return err
		}

		timestamp := time.Now().Format("20060102150405")
		migrationFile := fmt.Sprintf("db/migrations/%s_init.sql", timestamp)
		if err := generateFile("db/migrations/init.sql.tmpl", filepath.Join(cfg.OutputPath, "internal", migrationFile), cfg); err != nil {
			return err
		}

		if err := generateFile("db/queries/user.sql.tmpl", filepath.Join(cfg.OutputPath, "internal/db/queries/user.sql"), cfg); err != nil {
			return err
		}
	}

	return nil
}

func generateFile(tmplPath, outputPath string, data interface{}) error {
	tmpl, err := template.ParseFS(templates.FS, tmplPath)
	if err != nil {
		return err
	}

	// ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// create the file
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := tmpl.Execute(out, data); err != nil {
		return err
	}

	return nil
}
