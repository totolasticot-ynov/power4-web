<?php
// Contrôle d'authentification simple côté serveur.
// - Reçoit POST {username, password} depuis le formulaire.
// - Vérifie l'utilisateur en base et supporte deux formats stockés :
//   1) hash sécurisé (password_hash) => on utilise password_verify
//   2) ancien mot de passe en clair (legacy) => on accepte la connexion et on ré-hash pour sécuriser
// Retourne en texte brut : 'success' ou 'invalid'.

include '../../includes/db_connect.php';

header('Content-Type: text/plain; charset=utf-8');

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = isset($_POST['username']) ? trim($_POST['username']) : '';
    $password = isset($_POST['password']) ? trim($_POST['password']) : '';

    // Vérification minimale des champs
    if ($username === '' || $password === '') {
        echo 'invalid';
        exit;
    }

    // Récupère l'ID et la valeur stockée (hash ou ancien mot de passe)
    $stmt = $conn->prepare("SELECT id, password FROM users WHERE username = ?");
    if (!$stmt) {
        // Erreur de préparation : on renvoie 'invalid' pour éviter de divulguer des détails
        echo 'invalid';
        exit;
    }
    $stmt->bind_param('s', $username);
    $stmt->execute();
    $stmt->store_result();

    if ($stmt->num_rows === 0) {
        // utilisateur introuvable
        echo 'invalid';
        $stmt->close();
        $conn->close();
        exit;
    }

    $stmt->bind_result($id, $storedHash);
    $stmt->fetch();

    $ok = false;

    // Si la valeur en base est un hash (password_hash), password_verify renverra true
    if (password_verify($password, $storedHash)) {
        $ok = true;
    } else {
        // Compatibilité : si l'ancien mot de passe était stocké en clair,
        // on l'accepte (migration) puis on remplace par un hash sécurisé.
        if ($storedHash === $password) {
            $ok = true;
            $newHash = password_hash($password, PASSWORD_DEFAULT);
            $up = $conn->prepare("UPDATE users SET password = ? WHERE id = ?");
            if ($up) {
                $up->bind_param('si', $newHash, $id);
                $up->execute();
                $up->close();
            }
        }
    }

    if ($ok) {
        echo 'success';
    } else {
        echo 'invalid';
    }

    $stmt->close();
    $conn->close();
    exit;
}

?>
