

const boardEl = document.getElementById("board");
const resetBtn = document.getElementById("resetBtn");
let currentPlayer = 1;
let winner = 0;
const message = document.getElementById('message');

// API base: prefer an existing global `API_BASE` if another template set it (avoids redeclaration),
// otherwise compute a sensible default (relative when served by Go on :8080, or localhost otherwise)
const API_BASE_LOCAL = (typeof API_BASE !== 'undefined') ? API_BASE : ((location.port === '8080') ? '' : 'http://localhost:8080');

// Récupère la config stockée (localStorage)
let p4_rows = parseInt(localStorage.getItem('p4_rows') || '6');
let p4_cols = parseInt(localStorage.getItem('p4_cols') || '7');
let p4_win = parseInt(localStorage.getItem('p4_win') || '3');

// Envoie la config au backend si besoin
async function sendConfig() {
  try {
  const res = await fetch(API_BASE_LOCAL + '/api/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ rows: p4_rows, cols: p4_cols, win: p4_win })
    });
    return res.ok;
  } catch (err) {
    console.error('sendConfig error', err);
    return false;
  }
}

async function fetchBoard() {
  try {
  const res = await fetch(API_BASE_LOCAL + "/api/board");
    if (!res.ok) throw new Error('API returned ' + res.status);
    const state = await res.json();
    renderBoard(state);
    updateMessage(state);
  } catch (err) {
    console.error('fetchBoard error', err);
    if (message) message.innerHTML = "Impossible de contacter le backend du jeu (api). Démarre le serveur Go ou vérifie qu'il écoute sur http://localhost:8080";
  }
}

async function play(col) {
  if (winner !== 0) return;
  const mode = localStorage.getItem('p4_mode') || 'multi';
  // Validate column index
  if (typeof col !== 'number' || isNaN(col) || col < 0 || col >= p4_cols) {
    console.warn('play() called with invalid column:', col);
    if (message) message.innerHTML = 'Colonne invalide.';
    return;
  }
  try {
    console.log('Sending play to', API_BASE_LOCAL + '/api/play', { column: col, mode });
    const res = await fetch(API_BASE_LOCAL + "/api/play", {
      method: "POST",
      headers: { "Content-Type": "application/json", "X-P4-Mode": mode },
      body: JSON.stringify({ column: col }),
    });
    if (!res.ok) {
      const txt = await res.text().catch(() => '');
      throw new Error('play API ' + res.status + ' ' + txt);
    }
    // Mise à jour immédiate du plateau pour voir notre coup
    const newState = await res.json().catch(() => null);
    if (newState) {
      renderBoard(newState);
      updateMessage(newState);
    } else {
      fetchBoard();
    }
  } catch (err) {
    console.error('play error', err);
    if (message) message.innerHTML = 'Erreur: impossible d\'envoyer le coup au serveur. ' + (err && err.message ? err.message : '');
  }
}

async function resetGame() {
  try {
  await fetch(API_BASE_LOCAL + "/api/reset", { method: "POST" });
    fetchBoard();
  } catch (err) {
    console.error('reset error', err);
    if (message) message.innerHTML = 'Erreur: impossible de réinitialiser la partie.';
  }
}

let lastMove = null;


function renderBoard(state) {
  boardEl.innerHTML = "";
  currentPlayer = state.currentPlayer;
  winner = state.winner;

  // Prépare la liste des cases gagnantes (pour surbrillance)
  const winCells = Array.isArray(state.winCells) ? state.winCells.map(([r, c]) => `${r},${c}`) : [];

  // Trouver la dernière case jouée (pour l'animation)
  if (lastMove && state.board[lastMove.row][lastMove.col] === 0) {
    lastMove = null;
  }
  // Recherche la dernière case jouée (diff entre board et oldBoard)
  if (window.oldBoard) {
    for (let r = 0; r < state.board.length; r++) {
      for (let c = 0; c < state.board[r].length; c++) {
        if (window.oldBoard[r][c] !== state.board[r][c] && state.board[r][c] !== 0) {
          lastMove = { row: r, col: c, player: state.board[r][c] };
        }
      }
    }
  }

  // Affichage dynamique selon la config
  for (let r = 0; r < p4_rows; r++) {
    for (let c = 0; c < p4_cols; c++) {
      const cell = (state.board[r] || [])[c] || 0;
      const cellEl = document.createElement("div");
      cellEl.classList.add("cell");
      cellEl.addEventListener("click", () => play(c));

      if (cell !== 0) {
        const token = document.createElement("div");
        token.classList.add("token", cell === 1 ? "p1" : "p2");
        // Animation de chute classique
        if (lastMove && lastMove.row === r && lastMove.col === c) {
          token.classList.add("fall-real");
          token.style.setProperty('--fall-dist', `${(r) * 68}px`);
          token.style.setProperty('--fall-dur', `${0.12 + r*0.07}s`);
        }
        if (winCells.includes(`${r},${c}`)) {
          token.classList.add("win-token");
        }
        cellEl.appendChild(token);
      } else {
        cellEl.style.opacity = '0.5';
      }
      boardEl.appendChild(cellEl);
    }
  }
  // CSS grid dynamique
  boardEl.style.gridTemplateColumns = `repeat(${p4_cols}, 1fr)`;
  boardEl.style.gridTemplateRows = `repeat(${p4_rows}, 1fr)`;
  window.oldBoard = state.board.map(row => row.slice());
}

function updateMessage(state) {
  if (state.winner === 1) {
    message.innerHTML = 'Joueur 1 (<span class="jaune">jaune</span>) a gagné !';
  } else if (state.winner === 2) {
    message.innerHTML = 'Joueur 2 (<span class="rouge">rouge</span>) a gagné !';
  } else if (isDraw(state.board)) {
    message.innerHTML = '<span style="color:#5563DE;font-weight:bold;">Match nul&nbsp;!</span>';
  } else {
    const color = state.currentPlayer === 1 ? '<span class="jaune">jaune</span>' : '<span class="rouge">rouge</span>';
    message.innerHTML = `À ${color} de jouer.`;
  }
}

function isDraw(board) {
  for (let r = 0; r < board.length; r++) {
    for (let c = 0; c < board[r].length; c++) {
      if (board[r][c] === 0) return false;
    }
  }
  return true;
}

if (resetBtn) resetBtn.onclick = resetGame;

// Envoie la config au backend puis charge le plateau
sendConfig().then(fetchBoard);

