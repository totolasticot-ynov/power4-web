<?php
// Script d'inscription d'un utilisateur
// Reçoit username/password, vérifie doublon, hash le mot de passe et insère

include '../../includes/db_connect.php'; // Connexion à la base
header('Content-Type: text/plain; charset=utf-8'); // Réponse en texte brut

$DEBUG = true; // Active les messages d’erreur en développement

if ($_SERVER['REQUEST_METHOD'] === 'POST') { // Vérifie que la requête est POST
    $username = isset($_POST['username']) ? trim($_POST['username']) : ''; // Récupère username
    $password = isset($_POST['password']) ? trim($_POST['password']) : ''; // Récupère password

    if ($username === '' || $password === '') { // Vérifie champs vides
        echo 'error';
        exit;
    }

    // Fonction de log serveur (sans mots de passe)
    function log_msg($message) {
        $logDir = dirname(__DIR__, 2) . '/logs'; // Dossier logs
        if (!is_dir($logDir)) { @mkdir($logDir, 0777, true); } // Crée le dossier si absent
        $logFile = $logDir . '/register.log'; // Fichier de log
        $time = date('Y-m-d H:i:s'); // Timestamp
        @file_put_contents($logFile, "[$time] " . $message . PHP_EOL, FILE_APPEND | LOCK_EX); // Écrit log
    }

    // Vérifie si le username existe déjà
    $check = $conn->prepare("SELECT id FROM users WHERE username = ?");
    if (!$check) { // Erreur de préparation
        $err = $conn->error;
        log_msg("prepare SELECT failed for username='" . addslashes($username) . "' -- error: $err");
        if ($DEBUG) { echo 'error: ' . $err; } else { echo 'error'; }
        exit;
    }
    $check->bind_param('s', $username); // Lie la valeur
    $check->execute(); // Exécute la requête
    $check->store_result(); // Stocke le résultat
    if ($check->num_rows > 0) { // Username déjà utilisé
        echo 'exists';
        $check->close();
        $conn->close();
        exit;
    }
    $check->close(); // Ferme la requête

    // Hash du mot de passe et insertion
    $hash = password_hash($password, PASSWORD_DEFAULT); // Hash sécurisé
    $ins = $conn->prepare("INSERT INTO users (username, password) VALUES (?, ?)");
    if (!$ins) { // Erreur SQL
        $err = $conn->error;
        log_msg("prepare INSERT failed for username='" . addslashes($username) . "' -- error: $err");
        if ($DEBUG) { echo 'error: ' . $err; } else { echo 'error'; }
        exit;
    }
    $ins->bind_param('ss', $username, $hash); // Lie username + hash
    if ($ins->execute()) { // Tente l’insertion
        echo 'success'; // OK
    } else {
        $err = $ins->error; // Erreur d’exécution
        log_msg("execute INSERT failed for username='" . addslashes($username) . "' -- error: $err");
        if ($DEBUG) { echo 'error: ' . $err; } else { echo 'error'; }
    }

    $ins->close(); // Ferme la requête INSERT
    $conn->close(); // Ferme la connexion
    exit; // Termine le script
}

?>
