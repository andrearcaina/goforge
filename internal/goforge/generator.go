package goforge

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/andrearcaina/goforge/internal/templates"
)

func Generate(cfg *spec.Config) error {
	// parse the template from the embedded FS
	tmpl, err := template.ParseFS(templates.FS, "hello.go.tmpl")
	if err != nil {
		return err
	}

	var out *os.File
	if cfg.OutputPath != "" {
		// ensure the directory exists
		fullPath := filepath.Join(cfg.OutputPath, "hello.go")

		// if it doesn't exist, create it
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		// create the file
		f, err := os.Create(fullPath)
		if err != nil {
			return err
		}

		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	// execute the template (write to file or stdout)
	return tmpl.Execute(out, cfg)
}
