-- Schema for power4 DB
CREATE DATABASE IF NOT EXISTS `power4` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `power4`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(100) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `score` INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Example user (password in clear, but better to insert hashed)
-- INSERT INTO users (username, password) VALUES ('testuser', 'testpass');

-- To insert a hashed password in SQL directly (example using PHP's password_hash output):
-- INSERT INTO users (username, password) VALUES ('testuser', '$2y$10$...');
