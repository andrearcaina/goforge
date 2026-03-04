package ui

import (
	"fmt"
	"os"

	"github.com/andrearcaina/goforge/internal/config"
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

func OutputSuccess(cfg *config.Config) {
	fmt.Println(successStyle.Render("🔥 Generated successfully!"))
	fmt.Println("Make sure to run the following commands:")

	commands := buildCommands(cfg)
	printCommands(commands)

	os.Exit(0)
}

func buildCommands(cfg *config.Config) []string {
	makefile := cfg.Form.MakefileFlag
	db := cfg.Form.DatabaseFlag
	docker := cfg.Form.DockerFlag
	isREST := cfg.Form.ServerTypeFlag == config.REST
	isGRPC := cfg.Form.ServerTypeFlag == config.GRPC

	commands := []string{
		fmt.Sprintf("cd %s", cfg.OutputPath),
	}

	// docker and db commands if both are set
	if docker && db {
		commands = append(commands,
			"make docker-run  # make sure to have docker installed and running",
		)
	}

	// otherwise if only docker is set, then just print a message about the docker-compose file
	if docker && !db {
		commands = append(commands,
			"the docker-compose file only contains a db service, so you can ignore it since you didn't select the database option",
		)
	}

	// makefile and db commands
	if makefile && db {
		commands = append(commands,
			"make migrate-up  # make sure you have goose installed and your DB is running",
		)
	}

	// regardless if db is set or not, then just run make run if makefile is set
	if makefile {
		commands = append(commands, "make run")
	}

	// non makefile commands if makefile flag not set
	if !makefile {
		commands = append(commands, "go mod tidy")

		if db {
			commands = append(commands,
				"goose up -dir ./db/migrations postgres \"your_db_connection_string\"  # install goose and DB is running",
				"sqlc generate ./...                # install sqlc if you haven't already",
			)
		}

		if isREST {
			commands = append(commands,
				"swag init -g ./cmd/server/main.go  # install swag if you haven't already",
			)
		} else if isGRPC {
			commands = append(commands,
				"protoc --go_out=. --go_opt=paths=source_relative \\",
				"    --go-grpc_out=. --go-grpc_opt=paths=source_relative \\",
				"    internal/pb/user/user.proto    # install protoc and go plugins if you haven't already",
			)
		}

		commands = append(commands, "go run ./cmd/server/main.go")
	}

	// additional command for REST server to view API docs
	if isREST {
		commands = append(commands,
			"visit http://localhost:8000/swagger/index.html  # view API docs",
		)
	}

	return commands
}

func printCommands(commands []string) {
	for _, cmd := range commands {
		fmt.Println(cmdStyle.Render(cmd))
	}
}
