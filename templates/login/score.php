<?php
// Endpoint minimal pour gérer le score d'un utilisateur
// Usage (POST x-www-form-urlencoded):
// - Lire le score : POST { username }
//   -> retourne en clair le score (ex: "42")
// - Mettre à jour : POST { username, delta }
//   -> augmente (ou diminue) le score, retourne le score mis à jour
// Note : endpoint simple - pour production, ajouter une authentification (sessions/tokens).

include '../../includes/db_connect.php';
header('Content-Type: text/plain; charset=utf-8');

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    echo 'error';
    exit;
}

$username = isset($_POST['username']) ? trim($_POST['username']) : '';
if ($username === '') {
    echo 'error';
    exit;
}

// Si pas de 'delta' : on renvoie le score actuel
if (!isset($_POST['delta'])) {
    $stmt = $conn->prepare("SELECT score FROM users WHERE username = ?");
    if (!$stmt) { echo 'error'; exit; }
    $stmt->bind_param('s', $username);
    $stmt->execute();
    $stmt->store_result();
    if ($stmt->num_rows === 0) { echo 'error'; $stmt->close(); $conn->close(); exit; }
    $stmt->bind_result($score);
    $stmt->fetch();
    echo (int)$score;
    $stmt->close();
    $conn->close();
    exit;
}

// Mise à jour du score : applique delta et empêche valeur négative
$delta = intval($_POST['delta']);
$up = $conn->prepare("UPDATE users SET score = GREATEST(0, score + ?) WHERE username = ?");
if (!$up) { echo 'error'; exit; }
$up->bind_param('is', $delta, $username);
if (!$up->execute()) { echo 'error'; $up->close(); $conn->close(); exit; }
$up->close();

// Retourne le score mis à jour
$stmt2 = $conn->prepare("SELECT score FROM users WHERE username = ?");
if (!$stmt2) { echo 'error'; exit; }
$stmt2->bind_param('s', $username);
$stmt2->execute();
$stmt2->store_result();
if ($stmt2->num_rows === 0) { echo 'error'; $stmt2->close(); $conn->close(); exit; }
$stmt2->bind_result($score2);
$stmt2->fetch();
echo (int)$score2;
$stmt2->close();
$conn->close();
exit;
?>
