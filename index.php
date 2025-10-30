<?php
// Lightweight landing page so the project root doesn't auto-redirect.
// This avoids fast redirect loops with the XAMPP dashboard and lets you click the login link.
?><!doctype html>
<html lang="fr">
<head>
	<meta charset="utf-8">
	<title>power4-web — Accueil</title>
	<style>body{font-family:Segoe UI,Roboto,Arial;background:#f4f7ff;color:#222;display:flex;align-items:center;justify-content:center;height:100vh;margin:0} .card{background:#fff;padding:28px;border-radius:12px;box-shadow:0 6px 24px rgba(16,24,100,0.08);max-width:480px;text-align:center} a.button{display:inline-block;padding:12px 20px;background:#5563DE;color:#fff;border-radius:8px;text-decoration:none;font-weight:600;margin-top:12px}</style>
</head>
<body>
	<div class="card">
		<h1>power4-web</h1>
		<p>Bienvenue — clique sur le bouton pour ouvrir la page de connexion.</p>
		<p><a class="button" href="/power4-web/templates/login/login.html">Ouvrir la page de connexion</a></p>
		<p style="margin-top:12px;color:#666;font-size:0.9em">Si tu veux vider l'état de connexion stocké localement, ouvre la console et fais : <code>localStorage.removeItem('p4_user')</code></p>
	</div>
</body>
</html>

