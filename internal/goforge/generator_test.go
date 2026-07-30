package goforge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andrearcaina/goforge/internal/config"
)

func TestGenerateServerTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		serverType config.ServerTypeFlag
		wantFile   string
	}{
		{name: "REST", serverType: config.REST, wantFile: "internal/api/service.go"},
		{name: "gRPC", serverType: config.GRPC, wantFile: "internal/pb/user/user.proto"},
		{name: "GraphQL", serverType: config.GraphQL, wantFile: "graph/schemas/user.graphqls"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			outputPath := filepath.Join(t.TempDir(), "service")
			cfg := validConfig(outputPath, test.serverType)
			if err := Generate(cfg); err != nil {
				t.Fatalf("Generate() returned unexpected error: %v", err)
			}

			if _, err := os.Stat(filepath.Join(outputPath, test.wantFile)); err != nil {
				t.Fatalf("expected generated file %q: %v", test.wantFile, err)
			}
		})
	}
}

func TestGenerateGraphQLMountsEndpoint(t *testing.T) {
	t.Parallel()

	outputPath := filepath.Join(t.TempDir(), "service")
	if err := Generate(validConfig(outputPath, config.GraphQL)); err != nil {
		t.Fatalf("Generate() returned unexpected error: %v", err)
	}

	server, err := os.ReadFile(filepath.Join(outputPath, "internal/api/server.go"))
	if err != nil {
		t.Fatalf("read generated GraphQL server: %v", err)
	}
	if !strings.Contains(string(server), `router.Handle("/graphql", srv)`) {
		t.Fatal("generated GraphQL server does not mount the /graphql endpoint")
	}
}

func TestGenerateRefusesOverwriteAndForceReplaces(t *testing.T) {
	t.Parallel()

	outputPath := filepath.Join(t.TempDir(), "service")
	cfg := validConfig(outputPath, config.REST)
	if err := Generate(cfg); err != nil {
		t.Fatalf("initial Generate() returned unexpected error: %v", err)
	}

	goModPath := filepath.Join(outputPath, "go.mod")
	const customContents = "user-owned contents\n"
	if err := os.WriteFile(goModPath, []byte(customContents), 0644); err != nil {
		t.Fatalf("modify generated go.mod: %v", err)
	}

	err := Generate(cfg)
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("Generate() error = %v, want overwrite refusal", err)
	}
	contents, readErr := os.ReadFile(goModPath)
	if readErr != nil {
		t.Fatalf("read preserved go.mod: %v", readErr)
	}
	if string(contents) != customContents {
		t.Fatalf("go.mod changed after refused overwrite: got %q", contents)
	}

	cfg.Force = true
	if err := Generate(cfg); err != nil {
		t.Fatalf("Generate() with force returned unexpected error: %v", err)
	}
	contents, readErr = os.ReadFile(goModPath)
	if readErr != nil {
		t.Fatalf("read replaced go.mod: %v", readErr)
	}
	if string(contents) == customContents {
		t.Fatal("Generate() with force did not replace the existing file")
	}
}

func TestGenerateLeavesUnrelatedFilesAlone(t *testing.T) {
	t.Parallel()

	outputPath := filepath.Join(t.TempDir(), "service")
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatalf("create output directory: %v", err)
	}
	unrelatedPath := filepath.Join(outputPath, "notes.txt")
	if err := os.WriteFile(unrelatedPath, []byte("keep me\n"), 0644); err != nil {
		t.Fatalf("write unrelated file: %v", err)
	}

	if err := Generate(validConfig(outputPath, config.REST)); err != nil {
		t.Fatalf("Generate() returned unexpected error: %v", err)
	}
	contents, err := os.ReadFile(unrelatedPath)
	if err != nil {
		t.Fatalf("read unrelated file: %v", err)
	}
	if string(contents) != "keep me\n" {
		t.Fatalf("unrelated file changed: got %q", contents)
	}
}

func TestGenerateValidatesBeforeCreatingOutput(t *testing.T) {
	t.Parallel()

	parentPath := t.TempDir()
	outputPath := filepath.Join(parentPath, "service")
	cfg := validConfig(outputPath, "invalid")

	err := Generate(cfg)
	if err == nil || !strings.Contains(err.Error(), "unsupported server type") {
		t.Fatalf("Generate() error = %v, want unsupported server error", err)
	}
	if _, statErr := os.Stat(outputPath); !os.IsNotExist(statErr) {
		t.Fatalf("invalid generation created output path; stat error = %v", statErr)
	}

	stagingPaths, globErr := filepath.Glob(filepath.Join(parentPath, ".goforge-stage-*"))
	if globErr != nil {
		t.Fatalf("inspect staging paths: %v", globErr)
	}
	if len(stagingPaths) != 0 {
		t.Fatalf("generation left staging directories behind: %v", stagingPaths)
	}
}

func validConfig(outputPath string, serverType config.ServerTypeFlag) *config.Config {
	return &config.Config{
		OutputPath: outputPath,
		Form: config.Form{
			Name:           "example-service",
			ServerTypeFlag: serverType,
		},
	}
}
