const boardEl = document.getElementById("board");
const resetBtn = document.getElementById("resetBtn");
let currentPlayer = 1;
let winner = 0;

async function fetchBoard() {
  const res = await fetch("/api/board");
  const state = await res.json();
  renderBoard(state);
}

async function play(col) {
  if (winner !== 0) return;
  await fetch("/api/play", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ Column: col }),
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
      }

      boardEl.appendChild(cellEl);
    });
  });

  if (winner !== 0) {
    document.querySelector("h1").textContent = `🎉 Joueur ${winner} a gagné !`;
    // Animation sur les pions gagnants (bonus simple)
    document.querySelectorAll(".token").forEach(el => el.classList.add("winner"));
  } else {
    document.querySelector("h1").textContent = `Tour du joueur ${currentPlayer}`;
  }
}

resetBtn.addEventListener("click", resetGame);
fetchBoard();
