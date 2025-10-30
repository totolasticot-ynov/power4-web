<?php
$servername = "localhost";
$username = "root";
$password = "";
$dbname = "power4";

$conn = new mysqli($servername, $username, $password, $dbname);

if ($conn->connect_error) {
    die("Échec de la connexion : " . $conn->connect_error);
}

$conn->set_charset("utf8mb4");

?>
