package main

import (
	"log"
	"power4-web/src/menu"
)

func main() {
	// On lance le menu, et on gère une éventuelle erreur
	if err := menu.Menu(); err != nil {
		log.Fatalf("Erreur lors de l’exécution du menu : %v", err)
	}
}
