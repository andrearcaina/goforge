package goforge

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/andrearcaina/goforge/internal/config"
	"github.com/andrearcaina/goforge/internal/templates"
	"github.com/andrearcaina/goforge/internal/utils"
)

func Generate(cfg *config.Config) error {
	if err := utils.Validate(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	outputPath := filepath.Clean(cfg.OutputPath)
	parentPath := filepath.Dir(outputPath)
	if err := os.MkdirAll(parentPath, 0755); err != nil {
		return fmt.Errorf("create output parent directory: %w", err)
	}

	if info, err := os.Lstat(outputPath); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output path %q must not be a symbolic link", outputPath)
		}
		if !info.IsDir() {
			return fmt.Errorf("output path %q is not a directory", outputPath)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect output path: %w", err)
	}

	stagingPath, err := os.MkdirTemp(parentPath, ".goforge-stage-*")
	if err != nil {
		return fmt.Errorf("create staging directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(stagingPath)
	}()

	stagingCfg := *cfg
	stagingCfg.OutputPath = stagingPath
	if err := generateProject(&stagingCfg); err != nil {
		return fmt.Errorf("render project: %w", err)
	}

	if err := commitGeneratedProject(stagingPath, outputPath, cfg.Force); err != nil {
		return err
	}

	return nil
}

func generateProject(cfg *config.Config) error {
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

func commitGeneratedProject(stagingPath, outputPath string, force bool) error {
	if _, err := os.Lstat(outputPath); os.IsNotExist(err) {
		if err := os.Rename(stagingPath, outputPath); err != nil {
			return fmt.Errorf("commit generated project: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("inspect output path: %w", err)
	}

	var dirs []string
	var files []string
	err := filepath.WalkDir(stagingPath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == stagingPath {
			return nil
		}

		relativePath, err := filepath.Rel(stagingPath, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			dirs = append(dirs, relativePath)
		} else {
			files = append(files, relativePath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("inspect staged project: %w", err)
	}

	for _, relativePath := range dirs {
		targetPath := filepath.Join(outputPath, relativePath)
		if info, err := os.Lstat(targetPath); err == nil && !info.IsDir() {
			return fmt.Errorf("cannot create directory %q because a file already exists there", targetPath)
		} else if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("inspect destination %q: %w", targetPath, err)
		}
	}

	for _, relativePath := range files {
		targetPath := filepath.Join(outputPath, relativePath)
		if info, err := os.Lstat(targetPath); err == nil {
			if info.IsDir() {
				return fmt.Errorf("cannot generate file %q because a directory already exists there", targetPath)
			}
			if !force {
				return fmt.Errorf("refusing to overwrite existing file %q; use --force to overwrite generated files", targetPath)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect destination %q: %w", targetPath, err)
		}
	}

	for _, relativePath := range dirs {
		targetPath := filepath.Join(outputPath, relativePath)
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return fmt.Errorf("create destination directory %q: %w", targetPath, err)
		}
	}

	for _, relativePath := range files {
		sourcePath := filepath.Join(stagingPath, relativePath)
		targetPath := filepath.Join(outputPath, relativePath)
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("create destination directory for %q: %w", targetPath, err)
		}
		if force {
			if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("replace existing file %q: %w", targetPath, err)
			}
		}
		if err := os.Rename(sourcePath, targetPath); err != nil {
			return fmt.Errorf("commit generated file %q: %w", targetPath, err)
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

	if err := tmpl.Execute(out, data); err != nil {
		_ = out.Close()
		return err
	}

	return out.Close()
}
