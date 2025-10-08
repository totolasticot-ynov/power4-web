package menu

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// MenuHandler affiche la page du menu principal
func MenuHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "menu.html")
}

// JouerHandler redirige vers la page du jeu
func JouerHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ReglesHandler affiche une page avec les règles
func ReglesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "regles.html")
}

// QuitterHandler affiche un message de sortie
func QuitterHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "quitter.html")
}

// Fonction utilitaire pour rendre un template
func renderTemplate(w http.ResponseWriter, filename string) {
	tmplPath := filepath.Join("templates", filename)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
