package menu

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

type GameState struct {
	Board         [][]int  `json:"board"`
	CurrentPlayer int      `json:"currentPlayer"`
	Winner        int      `json:"winner"`
	WinCells      [][2]int `json:"winCells"`
}

var (
	rows, cols, winLen = 6, 7, 3
	state              = GameState{
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

// Helpers pour le bot (doivent être à la fin du fichier)
// (Anciennes fonctions utilitaires supprimées car non utilisées)

func Menu() error {
	// API: configurer la taille et la difficulté
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Rows int `json:"rows"`
			Cols int `json:"cols"`
			Win  int `json:"win"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Rows > 0 && req.Cols > 0 && req.Win > 0 {
			rows, cols, winLen = req.Rows, req.Cols, req.Win
			state.Board = make([][]int, rows)
			for i := range state.Board {
				state.Board[i] = make([]int, cols)
			}
			state.CurrentPlayer = 1
			state.Winner = 0
			state.WinCells = nil
		}
		w.WriteHeader(http.StatusOK)
	})
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
		// Place le pion du joueur humain
		for row := rows - 1; row >= 0; row-- {
			if state.Board[row][col] == 0 {
				state.Board[row][col] = state.CurrentPlayer
				state.CurrentPlayer = 3 - state.CurrentPlayer
				break
			}
		}
		state.Winner, state.WinCells = checkWinner(state.Board)
		// Si mode solo et pas de gagnant, le bot joue
		mode := r.Header.Get("X-P4-Mode")
		if mode == "solo" && state.Winner == 0 && state.CurrentPlayer == 2 {
			// Bot = joue un coup aléatoire valide
			validCols := []int{}
			for c := 0; c < cols; c++ {
				if state.Board[0][c] == 0 {
					validCols = append(validCols, c)
				}
			}
			if len(validCols) > 0 {
				botCol := validCols[rand.Intn(len(validCols))]
				botRow := getRowForCol(state.Board, botCol)
				if botRow != -1 {
					playBotMove(&state, botRow, botCol)
					state.Winner, state.WinCells = checkWinner(state.Board)
				}
			}
		}
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

	// Page de login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/login.html", nil)
	})
	// Page d'accueil (menu)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/menu/menu.html", nil)
	})

	// Page du jeu
	http.HandleFunc("/jeu", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/index/index.html", nil)
	})

	// Page du jeu - mode classique (6x7, align 4)
	http.HandleFunc("/jeu/classique", func(w http.ResponseWriter, r *http.Request) {
		rows, cols, winLen = 6, 7, 4
		state.Board = make([][]int, rows)
		for i := range state.Board {
			state.Board[i] = make([]int, cols)
		}
		state.CurrentPlayer = 1
		state.Winner = 0
		state.WinCells = nil
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
			if c+winLen-1 < cols {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r][c+k] != player {
						win = false
						break
					}
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ {
						cells[k] = [2]int{r, c + k}
					}
					return player, cells
				}
			}
			// Vertical
			if r+winLen-1 < rows {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r+k][c] != player {
						win = false
						break
					}
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ {
						cells[k] = [2]int{r + k, c}
					}
					return player, cells
				}
			}
			// Diagonal droite
			if r+winLen-1 < rows && c+winLen-1 < cols {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r+k][c+k] != player {
						win = false
						break
					}
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ {
						cells[k] = [2]int{r + k, c + k}
					}
					return player, cells
				}
			}
			// Diagonal gauche
			if r+winLen-1 < rows && c-winLen+1 >= 0 {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r+k][c-k] != player {
						win = false
						break
					}
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ {
						cells[k] = [2]int{r + k, c - k}
					}
					return player, cells
				}
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
