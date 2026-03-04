package goforge

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/andrearcaina/goforge/internal/config"
	"github.com/andrearcaina/goforge/internal/templates"
)

func Generate(cfg *config.Config) error {
	// generate base files
	if err := generateBaseFiles(cfg); err != nil {
		return err
	}

	// generate server files based on the selected server type
	switch cfg.Form.ServerTypeFlag {
	case config.REST:
		if err := generateRESTServer(cfg); err != nil {
			return err
		}
	case config.GRPC:
		if err := generateGRPCServer(cfg); err != nil {
			return err
		}
	}

	// generate database files if flag is set
	if cfg.Form.DatabaseFlag {
		if err := generateDBFiles(cfg); err != nil {
			return err
		}
	}

	return nil
}

func generateBaseFiles(cfg *config.Config) error {
	dirs := []string{
		"cmd/server",
		"internal/api",
		"internal/config",
	}

	if cfg.Form.DatabaseFlag {
		dirs = append(dirs, "internal/db", "internal/db/generated", "internal/db/migrations", "internal/db/queries")
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(cfg.OutputPath, dir), 0755); err != nil {
			return err
		}
	}

	mainTemplates := []string{"base/main.go.tmpl"}
	switch cfg.Form.ServerTypeFlag {
	case config.REST:
		mainTemplates = []string{"rest/main.go.tmpl"}
	case config.GRPC:
		mainTemplates = []string{"grpc/main.go.tmpl"}
	}

	if err := generateFile(mainTemplates, filepath.Join(cfg.OutputPath, "cmd/server/main.go"), cfg); err != nil {
		return err
	}

	if err := generateFile([]string{"base/config.go.tmpl"}, filepath.Join(cfg.OutputPath, "internal/config/config.go"), cfg); err != nil {
		return err
	}

	if err := generateFile([]string{"base/logger.go.tmpl"}, filepath.Join(cfg.OutputPath, "internal/logger/logger.go"), cfg); err != nil {
		return err
	}

	if cfg.Form.MakefileFlag {
		if err := generateFile([]string{"base/Makefile.tmpl"}, filepath.Join(cfg.OutputPath, "Makefile"), cfg); err != nil {
			return err
		}
	}

	if cfg.Form.DockerFlag {
		if err := generateFile([]string{"base/compose.yml.tmpl"}, filepath.Join(cfg.OutputPath, "docker-compose.yml"), cfg); err != nil {
			return err
		}

		if err := generateFile([]string{"base/.env.example.tmpl"}, filepath.Join(cfg.OutputPath, ".env"), cfg); err != nil {
			return err
		}
	}

	return nil
}

func generateRESTServer(cfg *config.Config) error {
	// specific go module
	if err := generateFile([]string{"rest/go.mod.tmpl"}, filepath.Join(cfg.OutputPath, "go.mod"), cfg); err != nil {
		return err
	}

	apiFiles := map[string]string{
		"rest/server.go.tmpl":  "internal/api/server.go",
		"rest/handler.go.tmpl": "internal/api/handler.go",
		"rest/service.go.tmpl": "internal/api/service.go",
		"rest/types.go.tmpl":   "internal/api/types.go",
	}

	for tmpl, out := range apiFiles {
		if err := generateFile([]string{tmpl}, filepath.Join(cfg.OutputPath, out), cfg); err != nil {
			return err
		}
	}

	return nil
}

func generateDBFiles(cfg *config.Config) error {
	dbFiles := map[string]string{
		"base/sqlc.yaml.tmpl": "sqlc.yaml",
		"db/init.sql.tmpl":    fmt.Sprintf("internal/db/migrations/%s_init.sql", time.Now().Format("20060102150405")),
		"db/user.sql.tmpl":    "internal/db/queries/user.sql",
		"db/db.go.tmpl":       "internal/db/db.go",
	}

	for tmpl, out := range dbFiles {
		if err := generateFile([]string{tmpl}, filepath.Join(cfg.OutputPath, out), cfg); err != nil {
			return err
		}
	}

	return nil
}

func generateFile(tmplPath []string, outputPath string, data interface{}) error {
	tmpl, err := template.ParseFS(templates.FS, tmplPath...)
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

func generateGRPCServer(cfg *config.Config) error {
	// specific go module
	if err := generateFile([]string{"grpc/go.mod.tmpl"}, filepath.Join(cfg.OutputPath, "go.mod"), cfg); err != nil {
		return err
	}

	apiFiles := map[string]string{
		"grpc/server.go.tmpl":   "internal/api/server.go",
		"grpc/handler.go.tmpl":  "internal/api/handler.go",
		"proto/user.proto.tmpl": "internal/pb/user/user.proto",
	}

	for tmpl, out := range apiFiles {
		if err := generateFile([]string{tmpl}, filepath.Join(cfg.OutputPath, out), cfg); err != nil {
			return err
		}
	}

	return nil
}
