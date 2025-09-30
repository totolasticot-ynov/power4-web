package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handler pour la page du jeu, sur /jeu
	http.HandleFunc("/jeu", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

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
			<a href="/jeu" style="padding:10px 20px; background:#28a745; color:#fff; text-decoration:none; border-radius:5px; display:inline-block;">Jouer au Puissance 4</a>
		</body>
		</html>`
		fmt.Fprint(w, html)
	})

	fmt.Println("Serveur démarré sur le port 8080...")
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
