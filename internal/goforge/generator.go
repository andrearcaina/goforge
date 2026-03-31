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
	if err := generateBaseFiles(cfg); err != nil {
		return err
	}

	switch cfg.Form.ServerTypeFlag {
	case config.REST:
		files := map[string]string{
			"rest/go.mod.tmpl":     "go.mod",
			"rest/server.go.tmpl":  "internal/api/server.go",
			"rest/handler.go.tmpl": "internal/api/handler.go",
			"rest/service.go.tmpl": "internal/api/service.go",
		}
		if err := generateRequiredFiles(cfg, files); err != nil {
			return err
		}
	case config.GRPC:
		files := map[string]string{
			"grpc/go.mod.tmpl":      "go.mod",
			"grpc/server.go.tmpl":   "internal/api/server.go",
			"grpc/handler.go.tmpl":  "internal/api/handler.go",
			"proto/user.proto.tmpl": "internal/pb/user/user.proto",
		}
		if err := generateRequiredFiles(cfg, files); err != nil {
			return err
		}
	case config.GraphQL:
		dirs := []string{
			"graph",
			"graph/generated",
			"graph/model",
			"graph/resolvers",
			"graph/schemas",
		}
		if err := generateDirs(cfg, dirs); err != nil {
			return err
		}

		files := map[string]string{
			"graphql/go.mod.tmpl":                "go.mod",
			"graphql/server.go.tmpl":             "internal/api/server.go",
			"graphql/gqlgen.yml.tmpl":            "gqlgen.yml",
			"graphql/schemas/user.graphqls.tmpl": "graph/schemas/user.graphqls",
		}
		if err := generateRequiredFiles(cfg, files); err != nil {
			return err
		}
	}

	if cfg.Form.DatabaseFlag {
		files := map[string]string{
			"base/sqlc.yaml.tmpl": "sqlc.yaml",
			"db/init.sql.tmpl":    fmt.Sprintf("internal/db/migrations/%s_init.sql", time.Now().Format("20060102150405")),
			"db/user.sql.tmpl":    "internal/db/queries/user.sql",
			"db/db.go.tmpl":       "internal/db/db.go",
		}
		if err := generateRequiredFiles(cfg, files); err != nil {
			return err
		}
	}

	return nil
}

func generateDirs(cfg *config.Config, dirs []string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(cfg.OutputPath, dir), 0755); err != nil {
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

	if err := generateDirs(cfg, dirs); err != nil {
		return err
	}

	// files to generate
	files := map[string]string{}

	switch cfg.Form.ServerTypeFlag {
	case config.REST:
		files["rest/main.go.tmpl"] = "cmd/server/main.go"
	case config.GRPC:
		files["grpc/main.go.tmpl"] = "cmd/server/main.go"
	case config.GraphQL:
		files["graphql/main.go.tmpl"] = "cmd/server/main.go"
	}

	files["base/config.go.tmpl"] = "internal/config/config.go"
	files["base/logger.go.tmpl"] = "internal/logger/logger.go"

	if cfg.Form.MakefileFlag {
		files["base/Makefile.tmpl"] = "Makefile"
	}

	if cfg.Form.DockerFlag {
		files["base/compose.yml.tmpl"] = "docker-compose.yml"
		files["base/.env.example.tmpl"] = ".env.example"
	}

	if err := generateRequiredFiles(cfg, files); err != nil {
		return err
	}

	return nil
}

func generateRequiredFiles(cfg *config.Config, files map[string]string) error {
	for tmpl, out := range files {
		if err := generateSpecificFile([]string{tmpl}, filepath.Join(cfg.OutputPath, out), cfg); err != nil {
			return err
		}
	}
	return nil
}

func generateSpecificFile(tmplPath []string, outputPath string, data interface{}) error {
	tmpl, err := template.ParseFS(templates.FS, tmplPath...)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

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
