package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur mon API Go !")
	})

	fmt.Println("Le port 8080 est utilisé pour lancer l'API Go !")
	http.ListenAndServe(":8080", nil)
}
