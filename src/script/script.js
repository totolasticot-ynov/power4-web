const rows = 6;
const cols = 7;
let currentPlayer = 'red';
let board = [];

const boardElement = document.getElementById('board');
const resetButton = document.getElementById('reset');

function initBoard() {
    board = Array.from({ length: rows }, () => Array(cols).fill(null));
    boardElement.innerHTML = '';
    for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
            const cell = document.createElement('div');
            cell.classList.add('cell');
            cell.dataset.row = r;
            cell.dataset.col = c;
            cell.addEventListener('click', handleCellClick);
            boardElement.appendChild(cell);
        }
    }
}

function handleCellClick(e) {
    const col = e.target.dataset.col;
    for (let r = rows - 1; r >= 0; r--) {
        if (!board[r][col]) {
            board[r][col] = currentPlayer;
            const cell = boardElement.querySelector(`.cell[data-row='${r}'][data-col='${col}']`);
            cell.classList.add(currentPlayer);
            if (checkWin(r, col)) {
                setTimeout(() => alert(`${currentPlayer.toUpperCase()} a gagné!`), 100);
                return;
            }
            currentPlayer = currentPlayer === 'red' ? 'yellow' : 'red';
            break;
        }
    }
}

function checkWin(r, c) {
    const directions = [
        [[0,1],[0,-1]],   // horizontal
        [[1,0],[-1,0]],   // vertical
        [[1,1],[-1,-1]],  // diagonale \
        [[1,-1],[-1,1]]   // diagonale /
    ];

    for (let dir of directions) {
        let count = 1;
        for (let [dr, dc] of dir) {
            let nr = r + dr;
            let nc = c + dc;
            while (nr >= 0 && nr < rows && nc >= 0 && nc < cols && board[nr][nc] === currentPlayer) {
                count++;
                nr += dr;
                nc += dc;
            }
        }
        if (count >= 4) return true;
    }
    return false;
}

resetButton.addEventListener('click', initBoard);

initBoard();
