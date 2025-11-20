<?php
// Connexion MySQL simple utilisée par les endpoints PHP
// - Modifie ces valeurs si ta configuration MySQL diffère
$servername = "localhost";
$username = "root";
$password = "";
$dbname = "power4";

// Création de la connexion mysqli
$conn = new mysqli($servername, $username, $password, $dbname);

// Vérification basique de la connexion
if ($conn->connect_error) {
    die("Échec de la connexion : " . $conn->connect_error);
}

// Forcer l'encodage UTF-8 (utf8mb4 pour gérer les emojis et caractères spéciaux)
$conn->set_charset("utf8mb4");

?>
