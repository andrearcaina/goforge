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
		Padding(1, 0).
		Render("🔥 Generated successfully!")

	cmdStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		PaddingLeft(2)

	fmt.Println(successStyle)
	fmt.Println("Make sure to run the following commands:")
	fmt.Println(cmdStyle.Render(fmt.Sprintf("cd %s", cfg.OutputPath)))
	fmt.Println(cmdStyle.Render("go mod tidy"))
	fmt.Println(cmdStyle.Render("swag init -g ./cmd/server/main.go (need to install swag first)"))

	if cfg.Form.DatabaseFlag {
		fmt.Println(cmdStyle.Render("sqlc generate ./... (you need to install sqlc first)"))
		fmt.Println(cmdStyle.Render("go run ./cmd/server/main.go (you need to run a database first and set up the .env file)"))
	} else {
		fmt.Println(cmdStyle.Render("go run ./cmd/server/main.go"))
	}
}
