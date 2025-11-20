<?php
// Authentification serveur simple : vérifie username/password
// Gère hash modernes + anciens mots de passe en clair

include '../../includes/db_connect.php'; // Connexion à la base
header('Content-Type: text/plain; charset=utf-8'); // Réponse en texte brut

if ($_SERVER['REQUEST_METHOD'] === 'POST') { // Vérifie que c'est un POST
    $username = isset($_POST['username']) ? trim($_POST['username']) : ''; // Récupère username
    $password = isset($_POST['password']) ? trim($_POST['password']) : ''; // Récupère password

    if ($username === '' || $password === '') { // Vérifie champs vides
        echo 'invalid';
        exit;
    }

    // Prépare la récupération du hash stocké
    $stmt = $conn->prepare("SELECT id, password FROM users WHERE username = ?");
    if (!$stmt) { // Erreur SQL
        echo 'invalid'; // Pas de détail pour la sécurité
        exit;
    }
    $stmt->bind_param('s', $username); // Lie la valeur
    $stmt->execute(); // Exécute la requête
    $stmt->store_result(); // Stocke le résultat

    if ($stmt->num_rows === 0) { // Aucun utilisateur trouvé
        echo 'invalid';
        $stmt->close();
        $conn->close();
        exit;
    }

    $stmt->bind_result($id, $storedHash); // Associe les colonnes
    $stmt->fetch(); // Récupère la ligne

    $ok = false; // Flag d'authentification

    // Vérifie via password_verify si hash sécurisé
    if (password_verify($password, $storedHash)) {
        $ok = true;
    } else {
        // Ancien cas : mot de passe stocké en clair
        if ($storedHash === $password) {
            $ok = true; // Connexion acceptée
            $newHash = password_hash($password, PASSWORD_DEFAULT); // Nouveau hash sécurisé
            $up = $conn->prepare("UPDATE users SET password = ? WHERE id = ?"); // Mise à jour
            if ($up) {
                $up->bind_param('si', $newHash, $id); // Lie hash + ID
                $up->execute(); // Exécute l'update
                $up->close(); // Ferme la requête
            }
        }
    }

    if ($ok) { // Authentification réussie
        echo 'success';
    } else { // Sinon
        echo 'invalid';
    }

    $stmt->close(); // Ferme la requête SELECT
    $conn->close(); // Ferme la connexion
    exit; // Termine
}

?>
