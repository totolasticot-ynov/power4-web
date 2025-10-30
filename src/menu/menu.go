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
	// small helper to set CORS headers and handle preflight
	setCORS := func(w http.ResponseWriter, r *http.Request) bool {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-P4-Mode")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return true
		}
		return false
	}

	// API: configurer la taille et la difficulté
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) {
			return
		}
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

	// API: register (tout accepté, aucune vérification)
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
		w.WriteHeader(http.StatusOK)
	})

	// API: login (tout accepté, aucune vérification)
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
		w.WriteHeader(http.StatusOK)
	})
	// Handlers statiques
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("src"))))

	// API: état du jeu
	http.HandleFunc("/api/board", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
		state.Winner, state.WinCells = checkWinner(state.Board)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(state)
	})

	// API: jouer un coup
	http.HandleFunc("/api/play", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
		if state.Winner != 0 {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(state)
			return
		}
		var req struct{ Column int }
		_ = json.NewDecoder(r.Body).Decode(&req)
		col := req.Column
		// simple validation: colonne dans les bornes
		if col < 0 || col >= cols {
			http.Error(w, "invalid column", http.StatusBadRequest)
			return
		}
		// Place le pion du joueur humain
		for rr := rows - 1; rr >= 0; rr-- {
			if state.Board[rr][col] == 0 {
				state.Board[rr][col] = state.CurrentPlayer
				state.CurrentPlayer = 3 - state.CurrentPlayer
				break
			}
		}
		state.Winner, state.WinCells = checkWinner(state.Board)
		// Si mode solo et pas de gagnant, le bot joue immédiatement (synchronement)
		mode := r.Header.Get("X-P4-Mode")
		if mode == "solo" && state.Winner == 0 && state.CurrentPlayer == 2 {
			// Bot amélioré : gagne si possible, sinon bloque, sinon aléatoire
			validCols := []int{}
			for c := 0; c < cols; c++ {
				if state.Board[0][c] == 0 {
					validCols = append(validCols, c)
				}
			}
			// 1. Cherche à gagner
			for _, c := range validCols {
				row := -1
				for r := rows - 1; r >= 0; r-- {
					if state.Board[r][c] == 0 {
						row = r
						break
					}
				}
				if row == -1 { continue }
				state.Board[row][c] = 2
				win, _ := checkWinner(state.Board)
				state.Board[row][c] = 0
				if win == 2 {
					state.Board[row][c] = 2
					state.CurrentPlayer = 1
					state.Winner, state.WinCells = checkWinner(state.Board)
					goto bot_end
				}
			}
			// 2. Bloque l'adversaire
			for _, c := range validCols {
				row := -1
				for r := rows - 1; r >= 0; r-- {
					if state.Board[r][c] == 0 {
						row = r
						break
					}
				}
				if row == -1 { continue }
				state.Board[row][c] = 1
				win, _ := checkWinner(state.Board)
				state.Board[row][c] = 0
				if win == 1 {
					state.Board[row][c] = 2
					state.CurrentPlayer = 1
					state.Winner, state.WinCells = checkWinner(state.Board)
					goto bot_end
				}
			}
			// 3. Sinon, joue aléatoirement
			if len(validCols) > 0 {
				botCol := validCols[rand.Intn(len(validCols))]
				for r := rows - 1; r >= 0; r-- {
					if state.Board[r][botCol] == 0 {
						state.Board[r][botCol] = 2
						state.CurrentPlayer = 1
						break
					}
				}
				state.Winner, state.WinCells = checkWinner(state.Board)
			}
		bot_end:
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(state)
	})

	// API: reset
	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
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
		renderTemplate(w, "templates/login/login.html", nil)
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

	log.Println("Serveur Go démarré sur http://localhost:8080 🚀")
	log.Println("Page de connexion (Go server): http://localhost:8080/login")
	log.Println("Si tu utilises Apache/XAMPP, page de connexion statique: http://127.0.0.1/power4-web/templates/login/login.html")
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
