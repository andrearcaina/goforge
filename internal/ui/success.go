package ui

import (
	"fmt"
	"os"

	"github.com/andrearcaina/goforge/internal/spec"
	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(1, 0)

	cmdStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			PaddingLeft(2)
)

func OutputSuccess(cfg *spec.Config) {
	fmt.Println(successStyle.Render("🔥 Generated successfully!"))
	fmt.Println("Make sure to run the following commands:")

	// base commands to run
	commands := []string{
		fmt.Sprintf("cd %s", cfg.OutputPath),
	}

	if cfg.Form.MakefileFlag {
		if cfg.Form.DatabaseFlag {
			commands = appendDatabaseCommands(commands)
		}

		commands = append(commands, "make")

		printCommands(commands)
	}

	commands = append(commands, "go mod tidy")

	// add commands based on selected options
	if cfg.Form.ServerTypeFlag == spec.REST {
		commands = append(commands, "swag init -g ./cmd/server/main.go  # install swag if you haven't already")
	}

	// if database is selected, add sqlc generate and server run commands
	if cfg.Form.DatabaseFlag {
		commands = appendDatabaseCommands(commands)
	}

	commands = append(commands, "go run ./cmd/server/main.go")

	// if REST server is selected, add command to view API docs
	if cfg.Form.ServerTypeFlag == spec.REST {
		commands = append(commands,
			"visit http://localhost:8000/swagger/index.html  # view API docs",
		)
	}

	printCommands(commands)
}

func appendDatabaseCommands(commands []string) []string {
	return append(commands,
		"sqlc generate ./...                # install sqlc if you haven't already",
		"goose up -dir ./db/migrations postgres \"your_db_connection_string\"  # install goose and DB is running",
	)
}

func printCommands(commands []string) {
	for _, cmd := range commands {
		fmt.Println(cmdStyle.Render(cmd))
	}

	// exit without error
	os.Exit(0)
}
