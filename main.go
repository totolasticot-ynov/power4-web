package menu

import (
	"log"

	"github.com/totolasticot-ynov/power4-web/src/menu"
)

func main() {
	if err := menu.Menu(); err != nil {
		log.Fatalf("Erreur lors de l’exécution du serveur : %v", err)
	}
}
