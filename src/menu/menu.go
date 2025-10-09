package menu

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type GameState struct {
	Board         [][]int  `json:"board"`
	CurrentPlayer int      `json:"currentPlayer"`
	Winner        int      `json:"winner"`
	WinCells      [][2]int `json:"winCells"`
}

var (
	rows, cols = 6, 7
	state      = GameState{
		Board:         make([][]int, 6),
		CurrentPlayer: 1,
		Winner:        0,
		WinCells:      nil,
	}
)

func init() {
	for i := range state.Board {
		state.Board[i] = make([]int, cols)
	}
}

func Menu() error {
	// Handlers statiques
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("src"))))

	// API: état du jeu
	http.HandleFunc("/api/board", func(w http.ResponseWriter, r *http.Request) {
		state.Winner, state.WinCells = checkWinner(state.Board)
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
		state.Winner, state.WinCells = checkWinner(state.Board)
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

	// Page d'accueil (menu)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/menu/menu.html", nil)
	})

	// Page du jeu
	http.HandleFunc("/jeu", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/index/index.html", nil)
	})

	log.Println("Serveur démarré sur http://localhost:8080 🚀")
	return http.ListenAndServe(":8080", nil)
}

func checkWinner(board [][]int) (int, [][2]int) {
	// Renvoie (winner, [cases gagnantes])
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			player := board[r][c]
			if player == 0 {
				continue
			}
			// Horizontal
			if c+3 < cols && player == board[r][c+1] && player == board[r][c+2] && player == board[r][c+3] {
				return player, [][2]int{{r, c}, {r, c + 1}, {r, c + 2}, {r, c + 3}}
			}
			// Vertical
			if r+3 < rows && player == board[r+1][c] && player == board[r+2][c] && player == board[r+3][c] {
				return player, [][2]int{{r, c}, {r + 1, c}, {r + 2, c}, {r + 3, c}}
			}
			// Diagonal droite
			if r+3 < rows && c+3 < cols && player == board[r+1][c+1] && player == board[r+2][c+2] && player == board[r+3][c+3] {
				return player, [][2]int{{r, c}, {r + 1, c + 1}, {r + 2, c + 2}, {r + 3, c + 3}}
			}
			// Diagonal gauche
			if r+3 < rows && c-3 >= 0 && player == board[r+1][c-1] && player == board[r+2][c-2] && player == board[r+3][c-3] {
				return player, [][2]int{{r, c}, {r + 1, c - 1}, {r + 2, c - 2}, {r + 3, c - 3}}
			}
		}
	}
	return 0, nil
}

func renderTemplate(w http.ResponseWriter, path string, data any) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, data)
}
