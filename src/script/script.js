
const boardEl = document.getElementById("board");
const resetBtn = document.getElementById("resetBtn");
let currentPlayer = 1;
let winner = 0;


const message = document.getElementById('message');

async function fetchBoard() {
  const res = await fetch("/api/board");
  const state = await res.json();
  renderBoard(state);
  updateMessage(state);
}

async function play(col) {
  if (winner !== 0) return;
  await fetch("/api/play", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ column: col }), // clé en minuscule !
  });
  fetchBoard();
}

async function resetGame() {
  await fetch("/api/reset", { method: "POST" });
  fetchBoard();
}

let lastMove = null;

function renderBoard(state) {
  boardEl.innerHTML = "";
  currentPlayer = state.currentPlayer;
  winner = state.winner;

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

  state.board.forEach((row, r) => {
    row.forEach((cell, c) => {
      const cellEl = document.createElement("div");
      cellEl.classList.add("cell");
      cellEl.addEventListener("click", () => play(c));

      if (cell !== 0) {
        const token = document.createElement("div");
        token.classList.add("token", cell === 1 ? "player1" : "player2");
        // Animation de chute si c'est le dernier pion joué
        if (lastMove && lastMove.row === r && lastMove.col === c) {
          // Animation de chute réaliste : durée et distance selon la ligne
          token.classList.add("fall-real");
          token.style.setProperty('--fall-dist', `${(r) * 68}px`); // 60px cell + 8px gap
          token.style.setProperty('--fall-dur', `${0.12 + r*0.07}s`);
        }
        cellEl.appendChild(token);
      } else {
        cellEl.style.opacity = '0.5';
      }

      boardEl.appendChild(cellEl);
    });
  });
  // Sauvegarde du plateau pour la prochaine animation
  window.oldBoard = state.board.map(row => row.slice());
}

function updateMessage(state) {
  if (state.winner === 1) {
    message.innerHTML = 'Joueur 1 (<span class="jaune">jaune</span>) a gagné !';
  } else if (state.winner === 2) {
    message.innerHTML = 'Joueur 2 (<span class="rouge">rouge</span>) a gagné !';
  } else {
    const color = state.currentPlayer === 1 ? '<span class="jaune">jaune</span>' : '<span class="rouge">rouge</span>';
    message.innerHTML = `À ${color} de jouer.`;
  }
}

if (resetBtn) resetBtn.onclick = resetGame;

fetchBoard();

