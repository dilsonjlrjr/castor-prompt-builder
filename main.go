package main

import (
	"fmt"
	"log"

	"github.com/dilsonrabelo/castor-prompt-builder/internal/parser"
)

func main() {
	roles, err := parser.LoadAllRoles("roles")
	if err != nil {
		log.Fatal(err)
	}
	models, err := parser.LoadAllModels("models")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Roles carregados: %d\n", len(roles))
	for _, r := range roles {
		fmt.Printf("  - %s (%s)\n", r.Nome, r.ID)
	}

	fmt.Printf("Modelos carregados: %d\n", len(models))
	for _, m := range models {
		fmt.Printf("  - %s (%s): %s\n", m.Nome, m.ID, m.Descricao)
	}
}
