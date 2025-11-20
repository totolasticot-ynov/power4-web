package menu // Déclare le package "menu"

import (
	"encoding/json"   // Gestion JSON
	"html/template"   // Rendu des templates HTML
	"log"             // Logs serveur
	"math/rand"       // Génération aléatoire
	"net/http"        // Serveur HTTP
)

type GameState struct {
	Board         [][]int  `json:"board"`         // Plateau de jeu
	CurrentPlayer int      `json:"currentPlayer"` // Joueur courant (1 ou 2)
	Winner        int      `json:"winner"`        // Gagnant (0 si aucun)
	WinCells      [][2]int `json:"winCells"`      // Cases gagnantes
}

// Définition des paramètres du jeu et état initial
var (
	rows, cols, winLen = 6, 7, 3 // Taille du plateau et longueur à aligner
	state              = GameState{
		Board:         make([][]int, 6), // Init des lignes
		CurrentPlayer: 1,                // Le joueur 1 commence
		Winner:        0,                // Aucun gagnant
		WinCells:      nil,              // Pas de cellules gagnantes
	}
)

func init() {
	for i := range state.Board {
		state.Board[i] = make([]int, cols) // Initialise chaque ligne
	}
}

// init crée un plateau vide rempli de 0

func Menu() error {
	// Gestion simple du CORS pour API + requêtes OPTIONS
	setCORS := func(w http.ResponseWriter, r *http.Request) bool {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Autorise tout le monde
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-P4-Mode")
		if r.Method == http.MethodOptions { // Préflight
			w.WriteHeader(http.StatusOK)
			return true
		}
		return false
	}

	// API : configuration du jeu (rows, cols, win)
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
		var req struct {
			Rows int `json:"rows"` // Nouvelles lignes
			Cols int `json:"cols"` // Nouvelles colonnes
			Win  int `json:"win"`  // Longueur gagnante
		}
		_ = json.NewDecoder(r.Body).Decode(&req) // Decode JSON
		if req.Rows > 0 && req.Cols > 0 && req.Win > 0 {
			rows, cols, winLen = req.Rows, req.Cols, req.Win // Met à jour les tailles
			state.Board = make([][]int, rows)                // Recrée le plateau
			for i := range state.Board {
				state.Board[i] = make([]int, cols)
			}
			state.CurrentPlayer = 1 // Reset joueur
			state.Winner = 0        // Reset gagnant
			state.WinCells = nil    // Reset cellules gagnantes
		}
		w.WriteHeader(http.StatusOK)
	})

	// API : register (aucune vérification)
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// API : login (aucune vérification)
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Fichiers statiques (assets)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("src"))))

	// API : retourne l’état du jeu
	http.HandleFunc("/api/board", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
		state.Winner, state.WinCells = checkWinner(state.Board) // Vérifie gagnant
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(state) // Renvoie JSON
	})

	// API : jouer un coup
	http.HandleFunc("/api/play", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }

		// Si la partie est finie
		if state.Winner != 0 {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(state)
			return
		}

		var req struct{ Column int }          // Colonne jouée
		_ = json.NewDecoder(r.Body).Decode(&req)
		col := req.Column

		// Vérifie la validité de la colonne
		if col < 0 || col >= cols {
			http.Error(w, "invalid column", http.StatusBadRequest)
			return
		}

		// Dépose le jeton du joueur courant
		for rr := rows - 1; rr >= 0; rr-- {
			if state.Board[rr][col] == 0 {           // Cherche case vide
				state.Board[rr][col] = state.CurrentPlayer
				state.CurrentPlayer = 3 - state.CurrentPlayer // Alterne joueur
				break
			}
		}

		state.Winner, state.WinCells = checkWinner(state.Board) // Vérifie victoire

		// Mode solo : le bot joue immédiatement
		mode := r.Header.Get("X-P4-Mode")
		if mode == "solo" && state.Winner == 0 && state.CurrentPlayer == 2 {

			// Collecte des colonnes jouables
			validCols := []int{}
			for c := 0; c < cols; c++ {
				if state.Board[0][c] == 0 {
					validCols = append(validCols, c)
				}
			}

			// 1. Bot cherche un coup gagnant
			for _, c := range validCols {
				row := -1
				for r := rows - 1; r >= 0; r-- {
					if state.Board[r][c] == 0 { row = r; break }
				}
				if row == -1 { continue }

				state.Board[row][c] = 2 // Simule coup bot
				win, _ := checkWinner(state.Board)
				state.Board[row][c] = 0 // Annule simulation
				if win == 2 { // Si gagnant => joue
					state.Board[row][c] = 2
					state.CurrentPlayer = 1
					state.Winner, state.WinCells = checkWinner(state.Board)
					goto bot_end // Sortie directe
				}
			}

			// 2. Bot bloque le joueur 1
			for _, c := range validCols {
				row := -1
				for r := rows - 1; r >= 0; r-- {
					if state.Board[r][c] == 0 { row = r; break }
				}
				if row == -1 { continue }

				state.Board[row][c] = 1 // Simule coup adverse
				win, _ := checkWinner(state.Board)
				state.Board[row][c] = 0 // Annule
				if win == 1 { // Bloque si victoire adverse
					state.Board[row][c] = 2
					state.CurrentPlayer = 1
					state.Winner, state.WinCells = checkWinner(state.Board)
					goto bot_end
				}
			}

			// 3. Sinon coup aléatoire
			if len(validCols) > 0 {
				botCol := validCols[rand.Intn(len(validCols))] // Choix random
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

		// Retour JSON du nouvel état
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(state)
	})

	// API : reset du plateau
	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		if setCORS(w, r) { return }
		for r := range state.Board {
			for c := range state.Board[r] {
				state.Board[r][c] = 0 // Vide chaque case
			}
		}
		state.CurrentPlayer = 1 // Reset joueur
		state.Winner = 0        // Reset gagnant
		w.WriteHeader(http.StatusOK)
	})

	// Route login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/login/login.html", nil)
	})

	// Route menu
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/menu/menu.html", nil)
	})

	// Route du jeu
	http.HandleFunc("/jeu", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "templates/index/index.html", nil)
	})

	// Route du jeu classique (6x7, aligner 4)
	http.HandleFunc("/jeu/classique", func(w http.ResponseWriter, r *http.Request) {
		rows, cols, winLen = 6, 7, 4 // Valeurs classiques
		state.Board = make([][]int, rows) // Recrée plateau
		for i := range state.Board {
			state.Board[i] = make([]int, cols)
		}
		state.CurrentPlayer = 1
		state.Winner = 0
		state.WinCells = nil
		renderTemplate(w, "templates/index/index.html", nil)
	})

	// Logs de démarrage
	log.Println("Serveur Go démarré sur http://localhost:8080 🚀")
	log.Println("Page de connexion (Go server): http://localhost:8080/login")
	log.Println("Si tu utilises Apache/XAMPP, page de connexion statique: http://127.0.0.1/power4-web/templates/login/login.html")

	// Lance le serveur HTTP
	return http.ListenAndServe(":8080", nil)
}

func checkWinner(board [][]int) (int, [][2]int) {
	// Parcourt tout le plateau pour détecter une ligne gagnante
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			player := board[r][c] // Case courante
			if player == 0 { continue }

			// Vérif horizontale
			if c+winLen-1 < cols {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r][c+k] != player { win = false; break }
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ { cells[k] = [2]int{r, c+k} }
					return player, cells
				}
			}

			// Vérif verticale
			if r+winLen-1 < rows {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r+k][c] != player { win = false; break }
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ { cells[k] = [2]int{r+k, c} }
					return player, cells
				}
			}

			// Diagonale ↘
			if r+winLen-1 < rows && c+winLen-1 < cols {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r+k][c+k] != player { win = false; break }
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ { cells[k] = [2]int{r+k, c+k} }
					return player, cells
				}
			}

			// Diagonale ↙
			if r+winLen-1 < rows && c-winLen+1 >= 0 {
				win := true
				for k := 1; k < winLen; k++ {
					if board[r+k][c-k] != player { win = false; break }
				}
				if win {
					cells := make([][2]int, winLen)
					for k := 0; k < winLen; k++ { cells[k] = [2]int{r+k, c-k} }
					return player, cells
				}
			}
		}
	}
	// Pas de gagnant
	return 0, nil
}

func renderTemplate(w http.ResponseWriter, path string, data any) {
	tmpl, err := template.ParseFiles(path) // Charge le fichier HTML
	if err != nil {
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, data) // Rend la page
}
