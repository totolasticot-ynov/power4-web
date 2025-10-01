const boardEl = document.getElementById("board");
const messageEl = document.getElementById("message");

let currentPlayer = 1;
let winner = 0;

async function fetchBoard() {
  const res = await fetch("/api/board");
  const data = await res.json();
  renderBoard(data.board);
  currentPlayer = data.currentPlayer;
  winner = data.winner;

  if (winner) {
    messageEl.textContent = `Joueur ${winner} a gagné ! 🎉`;
  } else {
    messageEl.textContent = `Tour du joueur ${currentPlayer} (${currentPlayer === 1 ? "Rouge" : "Jaune"})`;
  }
}

function renderBoard(board) {
  boardEl.innerHTML = "";
  for (let r = 0; r < board.length; r++) {
    const row = document.createElement("tr");
    for (let c = 0; c < board[r].length; c++) {
      const cell = document.createElement("td");
      cell.dataset.col = c;
      if (board[r][c] === 1) cell.classList.add("player1");
      else if (board[r][c] === 2) cell.classList.add("player2");
      cell.addEventListener("click", handleClick);
      row.appendChild(cell);
    }
    boardEl.appendChild(row);
  }
}

async function handleClick(e) {
  if (winner) return;

  const col = parseInt(e.target.dataset.col);
  const res = await fetch("/api/play", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ column: col })
  });

  const data = await res.json();
  renderBoard(data.board);
  currentPlayer = data.currentPlayer;
  winner = data.winner;

  if (winner) {
    messageEl.textContent = `Joueur ${winner} a gagné ! 🎉`;
  } else {
    messageEl.textContent = `Tour du joueur ${currentPlayer} (${currentPlayer === 1 ? "Rouge" : "Jaune"})`;
  }
}

async function resetGame() {
  await fetch("/api/reset", { method: "POST" });
  winner = 0;
  fetchBoard();
}

window.onload = fetchBoard;
