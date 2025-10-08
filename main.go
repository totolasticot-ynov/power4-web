package menu

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// MenuHandler affiche la page du menu principal
func MenuHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "menu.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Erreur template menu", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
