<?php
// Inscription d'un utilisateur
// - Reçoit POST {username, password}
// - Vérifie que le nom n'existe pas, hash le mot de passe et l'insère.
// - Retourne 'success', 'exists' ou 'error'.

include '../../includes/db_connect.php';

header('Content-Type: text/plain; charset=utf-8');

// Activer DEBUG uniquement en dev pour afficher les erreurs SQL
$DEBUG = true;

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = isset($_POST['username']) ? trim($_POST['username']) : '';
    $password = isset($_POST['password']) ? trim($_POST['password']) : '';

    if ($username === '' || $password === '') {
        echo 'error';
        exit;
    }

    // Helper de log côté serveur (ne contient jamais de mots de passe)
    function log_msg($message) {
        $logDir = dirname(__DIR__, 2) . '/logs';
        if (!is_dir($logDir)) {
            @mkdir($logDir, 0777, true);
        }
        $logFile = $logDir . '/register.log';
        $time = date('Y-m-d H:i:s');
        @file_put_contents($logFile, "[$time] " . $message . PHP_EOL, FILE_APPEND | LOCK_EX);
    }

    // Vérifie si le nom d'utilisateur existe déjà
    $check = $conn->prepare("SELECT id FROM users WHERE username = ?");
    if (!$check) {
        $err = $conn->error;
        log_msg("prepare SELECT failed for username='" . addslashes($username) . "' -- error: $err");
        if ($DEBUG) { echo 'error: ' . $err; } else { echo 'error'; }
        exit;
    }
    $check->bind_param('s', $username);
    $check->execute();
    $check->store_result();
    if ($check->num_rows > 0) {
        echo 'exists';
        $check->close();
        $conn->close();
        exit;
    }
    $check->close();

    // Hash du mot de passe puis insertion
    $hash = password_hash($password, PASSWORD_DEFAULT);
    $ins = $conn->prepare("INSERT INTO users (username, password) VALUES (?, ?)");
    if (!$ins) {
        $err = $conn->error;
        log_msg("prepare INSERT failed for username='" . addslashes($username) . "' -- error: $err");
        if ($DEBUG) { echo 'error: ' . $err; } else { echo 'error'; }
        exit;
    }
    $ins->bind_param('ss', $username, $hash);
    if ($ins->execute()) {
        echo 'success';
    } else {
        $err = $ins->error;
        log_msg("execute INSERT failed for username='" . addslashes($username) . "' -- error: $err");
        if ($DEBUG) { echo 'error: ' . $err; } else { echo 'error'; }
    }
    $ins->close();
    $conn->close();
    exit;
}

?>
