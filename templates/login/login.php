<?php
include '../../includes/db_connect.php';

header('Content-Type: text/plain; charset=utf-8');

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = isset($_POST['username']) ? trim($_POST['username']) : '';
    $password = isset($_POST['password']) ? trim($_POST['password']) : '';

    if ($username === '' || $password === '') {
        echo 'invalid';
        exit;
    }

    // Get the stored hash for this username
    $stmt = $conn->prepare("SELECT id, password FROM users WHERE username = ?");
    if (!$stmt) {
        echo 'invalid';
        exit;
    }
    $stmt->bind_param('s', $username);
    $stmt->execute();
    $stmt->store_result();

    if ($stmt->num_rows === 0) {
        // no such user
        echo 'invalid';
        $stmt->close();
        $conn->close();
        exit;
    }

    $stmt->bind_result($id, $storedHash);
    $stmt->fetch();

    $ok = false;

    // If stored value is a proper hash, use password_verify
    if (password_verify($password, $storedHash)) {
        $ok = true;
    } else {
        // Backwards compatibility: if storedHash looks like a plain password (legacy), allow login
        // and re-hash the password for future safety.
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
