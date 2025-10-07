package menu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Menu() error {
	// --- Logique du jeu Puissance 4 ---
	type GameState struct {
		Board         [][]int `json:"board"`
		CurrentPlayer int     `json:"currentPlayer"`
		Winner        int     `json:"winner"`
	}

	var (
		rows, cols = 6, 7
		state      = GameState{
			Board:         make([][]int, rows),
			CurrentPlayer: 1,
			Winner:        0,
		}
	)
	for i := range state.Board {
		state.Board[i] = make([]int, cols)
	}

	// Vérifie s'il y a un gagnant
	checkWinner := func(board [][]int) int {
		// Horizontal, vertical, diagonal
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				player := board[r][c]
				if player == 0 {
					continue
				}
				// Horizontal
				if c+3 < cols && player == board[r][c+1] && player == board[r][c+2] && player == board[r][c+3] {
					return player
				}
				// Vertical
				if r+3 < rows && player == board[r+1][c] && player == board[r+2][c] && player == board[r+3][c] {
					return player
				}
				// Diagonal droite
				if r+3 < rows && c+3 < cols && player == board[r+1][c+1] && player == board[r+2][c+2] && player == board[r+3][c+3] {
					return player
				}
				// Diagonal gauche
				if r+3 < rows && c-3 >= 0 && player == board[r+1][c-1] && player == board[r+2][c-2] && player == board[r+3][c-3] {
					return player
				}
			}
		}
		return 0
	}

	// API: état du jeu
	http.HandleFunc("/api/board", func(w http.ResponseWriter, r *http.Request) {
		state.Winner = checkWinner(state.Board)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(state)
	})

	// API: jouer un coup
	http.HandleFunc("/api/play", func(w http.ResponseWriter, r *http.Request) {
		if state.Winner != 0 {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(state)
			return
		}
		var req struct{ Column int }
		_ = json.NewDecoder(r.Body).Decode(&req)
		col := req.Column
		for r := rows - 1; r >= 0; r-- {
			if state.Board[r][col] == 0 {
				state.Board[r][col] = state.CurrentPlayer
				state.CurrentPlayer = 3 - state.CurrentPlayer
				break
			}
		}
		state.Winner = checkWinner(state.Board)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(state)
	})

	// API: reset
	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		for r := range state.Board {
			for c := range state.Board[r] {
				state.Board[r][c] = 0
			}
		}
		state.CurrentPlayer = 1
		state.Winner = 0
		w.WriteHeader(http.StatusOK)
	})
	fs := http.FileServer(http.Dir("./src"))
	http.Handle("/src/", http.StripPrefix("/src/", fs))

	// Handler pour la page du jeu, sur /jeu
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
	<!DOCTYPE html>
	<html lang="fr">
	<head>
		<meta charset="UTF-8">
		<title>Accueil</title>
		<link href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;600&display=swap" rel="stylesheet">
		<style>
			body {
				font-family: 'Poppins', sans-serif;
				background: linear-gradient(135deg, #74ABE2, #5563DE);
				color: white;
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				height: 100vh;
				margin: 0;
				text-align: center;
			}
			h1 { font-size: 2.2rem; margin-bottom: 1rem; }
			a.button { background: #ffffff; color: #5563DE; padding: 12px 25px; font-weight:600; text-decoration:none; border-radius:8px }
		</style>
	</head>
	<body>
		<h1>Bienvenue sur mon serveur Go 🚀</h1>
		<h1>Jouez au Puissance 4 contre un ami !</h1>
		<a href="/jeu" class="button">🎮 Jouer au Puissance 4</a>
	</body>
	</html>`
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/jeu", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<title>Puissance 4</title>
			<link rel="stylesheet" href="/src/style/style.css">
		</head>
		<body>
			<h1>Puissance 4</h1>
			<p id="message">Chargement...</p>
			<div id="board" class="board"></div>
			<button id="resetBtn">Recommencer</button>
			<script src="/src/script/script.js"></script>
		</body>
		</html>`
		fmt.Fprint(w, html)
	})
	fmt.Println("Serveur démarré sur http://localhost:8080 🚀")
	return http.ListenAndServe(":8080", nil)
}
