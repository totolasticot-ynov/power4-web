// Script principal du jeu (client)
 // Gère l'UI du plateau + appels API + score

// Eléments du DOM principaux
const boardEl = document.getElementById("board"); // Conteneur du plateau
const resetBtn = document.getElementById("resetBtn"); // Bouton reset
let currentPlayer = 1; // Joueur courant
let winner = 0;        // 0 = pas de gagnant
const message = document.getElementById('message'); // Zone d’affichage message

// API base: prend API_BASE global si existe sinon calcule par défaut
const API_BASE_LOCAL =
  (typeof API_BASE !== 'undefined')
    ? API_BASE
    : ((location.port === '8080') ? '' : 'http://localhost:8080');

// Récupère config stockée
let p4_rows = parseInt(localStorage.getItem('p4_rows') || '6'); // Nb lignes
let p4_cols = parseInt(localStorage.getItem('p4_cols') || '7'); // Nb colonnes
let p4_win = parseInt(localStorage.getItem('p4_win') || '3');   // Alignement gagnant

// Envoie config au backend (sync serveur)
async function sendConfig() {
  try {
    const res = await fetch(API_BASE_LOCAL + '/api/config', { // API config
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ rows: p4_rows, cols: p4_cols, win: p4_win })
    });
    return res.ok; // Retourne succès ou non
  } catch (err) {
    console.error('sendConfig error', err); // Log erreur
    return false; // Indique échec
  }
}

// Récupère l’état du plateau + MAJ UI
async function fetchBoard() {
  try {
    const res = await fetch(API_BASE_LOCAL + "/api/board"); // API board
    if (!res.ok) throw new Error('API returned ' + res.status); // Vérif status
    const state = await res.json(); // Parse JSON état
    renderBoard(state);             // Affiche plateau
    updateMessage(state);           // Met à jour message
  } catch (err) {
    console.error('fetchBoard error', err); // Log erreur
    if (message)
      message.innerHTML =
        "Impossible de contacter le backend du jeu."; // Message erreur UI
  }
}

// Joue un coup dans une colonne col
async function play(col) {
  if (winner !== 0) return; // Stop si partie finie
  const mode = localStorage.getItem('p4_mode') || 'multi'; // Mode solo/multi

  // Vérifie validité colonne
  if (typeof col !== 'number' || isNaN(col) || col < 0 || col >= p4_cols) {
    console.warn('Invalid column:', col); // Log warning
    if (message) message.innerHTML = 'Colonne invalide.'; // Message UI
    return;
  }

  try {
    // Envoie du coup au backend
    const res = await fetch(API_BASE_LOCAL + "/api/play", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-P4-Mode": mode // Indique mode au serveur
      },
      body: JSON.stringify({ column: col })
    });

    if (!res.ok) {
      const txt = await res.text().catch(() => ''); // Récup texte erreur
      throw new Error('play API ' + res.status + ' ' + txt);
    }

    const newState = await res.json().catch(() => null); // Parse réponse
    if (newState) {
      renderBoard(newState);   // MAJ plateau
      updateMessage(newState); // MAJ message
    } else fetchBoard();       // Fallback : recharger plateau
  } catch (err) {
    console.error('play error', err); // Log erreur
    if (message)
      message.innerHTML = 'Erreur: impossible d\'envoyer le coup.'; // UI erreur
  }
}

// Réinitialise partie côté serveur + nettoyage client
async function resetGame() {
  try {
    await fetch(API_BASE_LOCAL + "/api/reset", { method: "POST" }); // Reset backend
    scoreUpdated = false; // Reset flag score
    const oldScoreEl = document.getElementById('scoreChange'); // Élément score
    if (oldScoreEl && oldScoreEl.parentNode) oldScoreEl.remove(); // Supprime affichage score
    fetchBoard(); // Recharge plateau
  } catch (err) {
    console.error('reset error', err); // Log erreur
    if (message) message.innerHTML = 'Erreur: impossible de réinitialiser.'; // UI erreur
  }
}

let lastMove = null;      // Sauvegarde dernier coup
let scoreUpdated = false; // Évite double mise à jour score

// Déduction de la base projet (cas Apache)
const PROJECT_BASE =
  (location.href.indexOf('/power4-web') !== -1)
    ? '/power4-web'
    : '';
const SCORE_ENDPOINT = PROJECT_BASE + '/templates/login/score.php'; // API score

// Rafraîchit leaderboard si présent
function refreshLeaderboard() {
  try {
    const lb = document.getElementById('leaderboardList'); // UL leaderboard
    if (!lb) return; // Pas de leaderboard → stop
    const lbEndpoint = PROJECT_BASE + '/templates/login/leaderboard.php';

    fetch(lbEndpoint)
      .then(r => r.json())     // Parse JSON classement
      .then(list => {
        lb.innerHTML = '';     // Vide liste

        if (!Array.isArray(list) || list.length === 0) {
          const li = document.createElement('li'); // Ajout "aucun joueur"
          li.textContent = 'Aucun joueur';
          lb.appendChild(li);
          return;
        }

        const me = localStorage.getItem('p4_user'); // Nom user courant

        // Affiche chaque joueur
        list.forEach(item => {
          const li = document.createElement('li');
          li.textContent = item.username + ' — ' + item.score;

          if (me && item.username === me) {
            li.style.fontWeight = '700'; // Mets en gras l’utilisateur
            li.style.color = '#2b8a3e';  // Mets en vert
          }
          lb.appendChild(li);
        });
      })
      .catch(err => console.warn('Leaderboard error', err)); // Log erreur
  } catch (e) {
    console.warn('refreshLeaderboard error', e); // Log erreur
  }
}

// Construit l’affichage du plateau
function renderBoard(state) {
  boardEl.innerHTML = "";              // Vide plateau
  currentPlayer = state.currentPlayer; // MAJ joueur
  winner = state.winner;               // MAJ gagnant

  // Liste des cases gagnantes
  const winCells = Array.isArray(state.winCells)
    ? state.winCells.map(([r, c]) => `${r},${c}`)
    : [];

  // Si dernier coup supprimé → reset
  if (lastMove && state.board[lastMove.row][lastMove.col] === 0) {
    lastMove = null;
  }

  // Recherche du dernier coup (diff oldBoard)
  if (window.oldBoard) {
    for (let r = 0; r < state.board.length; r++) {
      for (let c = 0; c < state.board[r].length; c++) {
        if (window.oldBoard[r][c] !== state.board[r][c] &&
            state.board[r][c] !== 0) {
          lastMove = { row: r, col: c, player: state.board[r][c] }; // Sauvegarde
        }
      }
    }
  }

  // Affichage grille selon config
  for (let r = 0; r < p4_rows; r++) {
    for (let c = 0; c < p4_cols; c++) {
      const cell = (state.board[r] || [])[c] || 0; // Valeur case
      const cellEl = document.createElement("div"); // Div case
      cellEl.classList.add("cell");                 // Classe CSS

      // Clic = joue dans cette colonne
      cellEl.addEventListener("click", () => play(c));

      // Si pion présent
      if (cell !== 0) {
        const token = document.createElement("div");     // Jeton
        token.classList.add("token", cell === 1 ? "p1" : "p2"); // Couleur

        // Animation si dernier coup
        if (lastMove && lastMove.row === r && lastMove.col === c) {
          token.classList.add("fall-real"); // Classe animation
          token.style.setProperty('--fall-dist', `${r * 68}px`); // Distance
          token.style.setProperty('--fall-dur', `${0.12 + r * 0.07}s`); // Durée
        }

        // Surbrillance si case gagnante
        if (winCells.includes(`${r},${c}`)) {
          token.classList.add("win-token");
        }

        cellEl.appendChild(token); // Place jeton dans case
      } else {
        cellEl.style.opacity = '0.5'; // Case vide gris
      }

      boardEl.appendChild(cellEl); // Ajoute case au plateau
    }
  }

  // Ajuste taille grille CSS
  boardEl.style.gridTemplateColumns = `repeat(${p4_cols}, 1fr)`;
  boardEl.style.gridTemplateRows = `repeat(${p4_rows}, 1fr)`;

  // Sauvegarde board actuel pour détection dernier coup
  window.oldBoard = state.board.map(row => row.slice());
}

// Met à jour message d’état
function updateMessage(state) {
  if (state.winner === 1) {
    message.innerHTML = 'Joueur 1 (<span class="jaune">jaune</span>) a gagné !';
  } else if (state.winner === 2) {
    message.innerHTML = 'Joueur 2 (<span class="rouge">rouge</span>) a gagné !';
  } else if (isDraw(state.board)) {
    message.innerHTML = '<span style="color:#5563DE;font-weight:bold;">Match nul !</span>';
  } else {
    const color =
      state.currentPlayer === 1
        ? '<span class="jaune">jaune</span>'
        : '<span class="rouge">rouge</span>';
    message.innerHTML = `À ${color} de jouer.`; // Message tour
  }

  // Gestion score solo
  try {
    const mode = localStorage.getItem('p4_mode') || 'multi'; // Mode
    const user = localStorage.getItem('p4_user') || null;    // Nom joueur

    // Score uniquement en solo + partie finie + pas encore MAJ
    if (mode === 'solo' && user && state.winner !== 0 && !scoreUpdated) {
      const delta = (state.winner === 1) ? 1 : -1; // +1 si victoire, -1 si défaite

      // Affichage du delta sous le message
      let scoreEl = document.getElementById('scoreChange');
      if (!scoreEl) {
        scoreEl = document.createElement('div');
        scoreEl.id = 'scoreChange';
        scoreEl.style.marginTop = '8px';
        scoreEl.style.fontWeight = '700';
        scoreEl.style.fontSize = '0.95em';
        scoreEl.style.color =
          state.winner === 1 ? '#2b8a3e' : '#c0392b'; // Vert/rouge
        message.parentNode.insertBefore(scoreEl, message.nextSibling);
      }

      scoreEl.textContent =
        (delta > 0 ? '+' + delta : delta) + ' point' + (Math.abs(delta) > 1 ? 's' : '');

      // Envoi du score au serveur PHP
      fetch(SCORE_ENDPOINT, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ username: user, delta: String(delta) })
      })
        .then(res => res.text())
        .then(txt => {
          const newScore = parseInt(txt, 10); // Score mis à jour
          if (!isNaN(newScore)) {
            const sd = document.getElementById('scoreDisplay');
            if (sd) sd.textContent = 'Score: ' + newScore; // MAJ UI score
          }
          scoreUpdated = true; // Empêche double update
          refreshLeaderboard(); // Rafraîchit classement
        })
        .catch(err => {
          console.warn('score update failed', err); // Log erreur
          if (scoreEl) scoreEl.textContent += ' (non sauvegardé)'; // Indique échec
        });
    }
  } catch (e) {
    console.warn('score update error', e); // Log erreur
  }
}

// Vérifie si plateau plein
function isDraw(board) {
  for (let r = 0; r < board.length; r++) {
    for (let c = 0; c < board[r].length; c++) {
      if (board[r][c] === 0) return false; // Si case vide -> pas nul
    }
  }
  return true; // Plateau plein -> match nul
}

if (resetBtn) resetBtn.onclick = resetGame; // Bouton reset → resetGame()

// Envoie config au backend, puis récupère plateau
sendConfig().then(fetchBoard);
