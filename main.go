package main

import (
	"log"

	"github.com/totolasticot-ynov/power4-web/src/menu"
)

func main() {
	// Lance le serveur Go qui sert l'API du jeu et les templates
	if err := menu.Menu(); err != nil {
		log.Fatalf("Erreur lors de l’exécution du serveur : %v", err)
	}
}
