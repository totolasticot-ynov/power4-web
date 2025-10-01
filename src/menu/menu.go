package menu

import (
	"fmt"
	"net/http"
)

func Menu() error {
	fs := http.FileServer(http.Dir("/src"))
	http.Handle("/src/", http.StripPrefix("/src/", fs))

	// Handler pour la page du jeu, sur /jeu
	http.HandleFunc("/jeu", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./src/index.html")
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

	return http.ListenAndServe(":8080", nil)
}
