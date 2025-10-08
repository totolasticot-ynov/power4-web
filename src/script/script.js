
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

function renderBoard(state) {
  boardEl.innerHTML = "";
  currentPlayer = state.currentPlayer;
  winner = state.winner;

  state.board.forEach((row, r) => {
    row.forEach((cell, c) => {
      const cellEl = document.createElement("div");
      cellEl.classList.add("cell");
      cellEl.addEventListener("click", () => play(c));

      if (cell !== 0) {
        const token = document.createElement("div");
        token.classList.add("token", cell === 1 ? "player1" : "player2");
        cellEl.appendChild(token);
      } else {
        cellEl.style.opacity = '0.5';
      }

      boardEl.appendChild(cellEl);
    });
  });
}

function updateMessage(state) {
  if (state.winner === 1) {
    message.textContent = 'Joueur 1 (jaune) a gagné !';
  } else if (state.winner === 2) {
    message.textContent = 'Joueur 2 (rouge) a gagné !';
  } else {
    message.textContent = `À ${state.currentPlayer === 1 ? 'jaune' : 'rouge'} de jouer.`;
  }
}

if (resetBtn) resetBtn.onclick = resetGame;

fetchBoard();

