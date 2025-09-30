package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bienvenue sur mon serveur Go !")
	})

	fmt.Println("Serveur démarré sur le port 8080...")
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

	// Dossier qui contient tes fichiers HTML, CSS, JS, images...
	fs := http.FileServer(http.Dir("./static"))

	// Tout ce qui commence par /static ira chercher dans ./static
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Page d’accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// Lancer le serveur
	http.ListenAndServe(":8080", nil)
}
