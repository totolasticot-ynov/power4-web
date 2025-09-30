package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<title>Accueil</title>
		</head>
		<body>
			<h1>Bienvenue sur mon serveur Go !</h1>
			<button onclick="window.location.href='/static/index.html'">Jouer au Puissance 4</button>
		</body>
		</html>`
		fmt.Fprint(w, html)
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
