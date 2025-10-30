<?php
// Retourne le top N joueurs par score en JSON
// Usage: GET returns JSON array [{"username":"foo","score":42}, ...]
header('Content-Type: application/json; charset=utf-8');
include '../../includes/db_connect.php';

$limit = 5;
$sql = "SELECT username, score FROM users ORDER BY score DESC, id ASC LIMIT ?";
$stmt = $conn->prepare($sql);
if (!$stmt) {
    echo json_encode([]);
    exit;
}
$stmt->bind_param('i', $limit);
$stmt->execute();
$res = $stmt->get_result();
$out = [];
while ($row = $res->fetch_assoc()) {
    $out[] = ['username' => $row['username'], 'score' => (int)$row['score']];
}
echo json_encode($out);
$stmt->close();
$conn->close();
exit;
?>