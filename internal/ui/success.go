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

	if cfg.Form.DockerFlag {
		if cfg.Form.DatabaseFlag {
			commands = append(commands, "make docker-run  # make sure to have docker installed and running")
		} else {
			commands = append(commands, "the docker-compose file only contains a db service, so you can ignore it since you didn't select the database option")
		}
	}

	if cfg.Form.MakefileFlag {
		if cfg.Form.DatabaseFlag {
			commands = append(commands, "make migrate-up  # make sure you have goose installed and your DB is running")
		}

		commands = append(commands, "make run")

		printCommands(cfg.Form.ServerTypeFlag, commands)
	}

	commands = append(commands, "go mod tidy")

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

	commands = append(commands, "go run ./cmd/server/main.go")

	printCommands(cfg.Form.ServerTypeFlag, commands)
}

func printCommands(serverType spec.ServerTypeFlag, commands []string) {
	if serverType == spec.REST {
		commands = append(commands,
			"visit http://localhost:8000/swagger/index.html  # view API docs",
		)
	}

	for _, cmd := range commands {
		fmt.Println(cmdStyle.Render(cmd))
	}

	// exit without error
	os.Exit(0)
}
