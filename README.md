# 🎮 Power4-Web — Le Puissance 4 le plus beau du web !

Bienvenue dans **Power4-Web**, un projet qui mélange Go, PHP, HTML, JS et un peu de magie ✨  
But du jeu ? Gagner, évidemment. Et exploser le classement 🏆.

---
 **Ce projet vise à développer une application web interactive complète pour le jeu Puissance 4 🟡🔴. L'objectif est de proposer une plateforme où les utilisateurs peuvent s'affronter en ligne. Le système gère l'authentification et l'inscription des joueurs (via PHP), stocke leurs données et scores dans une base de données SQL 💾, et affiche un classement 🏆 dynamique. La logique métier principale du jeu est assurée par le code Go 💚, tandis que le frontend utilise les standards HTML, CSS et JavaScript 💻 pour garantir une expérience utilisateur fluide et agréable.**

Backend (Go)	Go 💚	1.18+	Nécessaire pour la logique du jeu (menu.go) et le serveur.
Backend (PHP)	PHP 🐘	7.4+	Indispensable pour l'authentification (login.php, register.php) et la connexion DB.
Base de Données	Serveur SQL 💾	MySQL / PostgreSQL	Un serveur de base de données (et un outil client comme DBeaver) pour initialiser db/schema.sql.
Frontend	Navigateur Web 🌐	Moderne (Chrome, Firefox)	Pour afficher le HTML, le CSS et exécuter script.js.


## 🗂️ Arborescence du projet 

```bash
POWER4
├── assets/          # 💅 Style.css 
├── db/              # 🧠 Le cerveau SQL
├── includes/        # 🔌 Connexion DB 
├── src/             # 🛠️ Code utile 
│   ├── menu/        # 📘 Menu Go
│   └── script/      # 💡 Scripts JS
├── templates/       # 📄 Pages HTML/PHP
│   ├── index/       # 🏠 Accueil
│   └── login/       # 🔐 Login / Score / Classement
├── menu/            # 📄 Menu principal HTML
├── main.go          # 🧠 Serveur Go
└── index.php        # 🚀 Entrée du site
ble)

1. 💾 Gestion de la Base de Données
Le fichier includes/db_connect.php contient la fonction nécessaire à l'établissement de la connexion à la base de données SQL. Cette fonction est un point de dépendance crucial pour tous les autres scripts PHP qui interagissent avec les données (utilisateurs, scores).

2. 👤 Fonctions d'Authentification (PHP)
Le script register.php gère l'enregistrement de nouveaux utilisateurs. Il inclut la logique pour insérer les données d'un nouvel utilisateur dans la base de données après avoir appliqué un hachage sécurisé au mot de passe.

Le script login.php contient la logique de vérification. Il est responsable de vérifier les identifiants de l'utilisateur contre les enregistrements de la base de données et, en cas de succès, de démarrer une session pour le joueur connecté.

3. 🎲 Logique du Jeu (Go)
Les fichiers Go sont dédiés à la logique du jeu Puissance 4 :

Le code contenu dans src/menu/menu.go inclut la fonction qui permet de placer le jeton d'un joueur dans une colonne spécifiée et de mettre à jour l'état interne du plateau de jeu.

Ce même module contient la fonction essentielle qui est appelée après chaque coup pour déterminer si un joueur a gagné en vérifiant l'alignement de quatre jetons (horizontal, vertical ou diagonal) à partir de la dernière position jouée.

4. 📊 Suivi des Scores (PHP/SQL)
La page leaderboard.php exécute la requête nécessaire pour récupérer les scores des joueurs, trier ces données par rang, et les afficher sous forme de classement.

Le script score.php contient la logique permettant de mettre à jour les statistiques et le score d'un joueur dans la base de données une fois qu'une partie est terminée.


🏗️ Décisions d'Architecture & Compromis Techniques
L'architecture du projet POWER4-WEB révèle une décision clé de coupler plusieurs technologies backend (Go et PHP). Cette approche permet de tirer parti de la rapidité et de l'efficacité de Go 💚 pour la logique complexe et critique du jeu (vérification de victoire, gestion du plateau), tout en utilisant la facilité d'intégration et la maturité de PHP 🐘 pour la gestion des pages web (HTML), des sessions utilisateur, et de la base de données SQL. Le compromis réside dans la complexité accrue de la gestion du déploiement, de la communication inter-processus (comment PHP interagit avec le moteur de jeu Go), et de la maintenance, par rapport à une solution entièrement basée sur un seul langage (comme Go pour tout, ou PHP avec une bibliothèque de jeu).


Go 💚	gofmt (intégré) et go vet	Formateur canonique pour le code Go. go vet est un linter qui vérifie les erreurs statiques 
et les constructions suspectes.

PHP 🐘	PHP-CS-Fixer ou PHP_CodeSniffer	Formatters/Linters qui appliquent des standards de codage (ex: PSR-1, PSR-12) aux fichiers .php.

JavaScript 💻	ESLint et Prettier	ESLint sert de linter pour imposer des règles de codage JavaScript. Prettier est utilisé comme formateur pour maintenir un style cohérent dans les fichiers .js et potentiellement .html/.css.

CSS	Stylelint ou Prettier	Linter pour les fichiers CSS, garantissant la cohérence et la validité du style.


Transformer le module Go en un service API REST (ou gRPC) complet qui gère non seulement la logique du jeu mais aussi les interactions avec la base de données (si les performances sont critiques) et les sessions de jeu.
Utiliser des variables d'environnement pour stocker les identifiants de la base de données et les clés secrètes, plutôt que de les coder en dur dans db_connect.php.



MIT License

Copyright (c) [2025] [Dhordain Thomas, Beyney Thomas]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.