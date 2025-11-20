<?php
// Indique que la réponse sera du JSON
header('Content-Type: application/json; charset=utf-8');
// Inclut la connexion à la base de données
include '../../includes/db_connect.php';

// Valeur par défaut du nombre de résultats
$n = 5;
// Récupère n depuis GET si présent
if (isset($_GET['n'])) { $n = intval($_GET['n']); }
// Récupère n depuis POST si présent
if (isset($_POST['n'])) { $n = intval($_POST['n']); }
// Réassigne 5 si n est invalide
if ($n <= 0) { $n = 5; }

// Requête pour récupérer les meilleurs scores
$sql = "SELECT username, score FROM users ORDER BY score DESC, id ASC LIMIT ?";
// Prépare la requête sécurisée
$stmt = $conn->prepare($sql);
// En cas d'erreur de préparation, renvoie un tableau vide
if (!$stmt) {
    echo json_encode([]);
    $conn->close();
    exit;
}
// Lie le paramètre LIMIT à la requête
$stmt->bind_param('i', $n);
// Exécute la requête
$stmt->execute();
// Récupère le résultat
$res = $stmt->get_result();
// Initialise le tableau de sortie
$out = [];
// Parcourt chaque ligne du résultat
while ($row = $res->fetch_assoc()) {
    // Ajoute username et score sous forme typée
    $out[] = ['username' => $row['username'], 'score' => (int)$row['score']];
}

// Renvoie le tableau final en JSON
echo json_encode($out);
// Ferme la requête
$stmt->close();
// Ferme la connexion
$conn->close();
// Termine le script
exit;
?>
