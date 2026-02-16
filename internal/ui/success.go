package ui

import (
	"fmt"

	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/charmbracelet/lipgloss"
)

func OutputSuccess(cfg *spec.Config) {
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(1, 0)

	cmdStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		PaddingLeft(2)

	fmt.Println(successStyle.Render("🔥 Generated successfully!"))
	fmt.Println("Make sure to run the following commands:")

	// base commands to run
	commands := []string{
		fmt.Sprintf("cd %s", cfg.OutputPath),
		"go mod tidy",
	}

	// add commands based on selected options
	if cfg.Form.ServerTypeFlag == spec.REST {
		commands = append(commands, "swag init -g ./cmd/server/main.go  # install swag if you haven't already")
	}

	// if database is selected, add sqlc generate and server run commands
	if cfg.Form.DatabaseFlag {
		commands = append(commands,
			"sqlc generate ./...                # install sqlc if you haven't already",
			"goose up -dir ./db/migrations postgres \"your_db_connection_string\"  # install goose and DB is running",
		)
	}

	if cfg.Form.MakefileFlag {
		commands = append(commands, "make")
	} else {
		commands = append(commands, "go run ./cmd/server/main.go")
	}

	// if REST server is selected, add command to view API docs
	if cfg.Form.ServerTypeFlag == spec.REST {
		commands = append(commands,
			"visit http://localhost:8000/swagger/index.html  # view API docs",
		)
	}

	for _, cmd := range commands {
		fmt.Println(cmdStyle.Render(cmd))
	}
}
