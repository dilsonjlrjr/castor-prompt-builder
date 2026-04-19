package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dilsonrabelo/castor-prompt-builder/internal/parser"
	"github.com/dilsonrabelo/castor-prompt-builder/internal/tui"
)

func main() {
	roles, err := parser.LoadAllRoles("roles")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao carregar roles:", err)
		os.Exit(1)
	}
	if len(roles) == 0 {
		fmt.Fprintln(os.Stderr, "Nenhum role encontrado em roles/")
		os.Exit(1)
	}

	models, err := parser.LoadAllModels("models")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao carregar models:", err)
		os.Exit(1)
	}
	if len(models) == 0 {
		fmt.Fprintln(os.Stderr, "Nenhum model encontrado em models/")
		os.Exit(1)
	}

	m := tui.New(models, roles)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
