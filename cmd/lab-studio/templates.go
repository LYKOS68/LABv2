package main

import (
	"html/template"
	"log"
	"net/http"
)

// Startseite und Konto-Seiten: layout.html ist der Rahmen, home.html und auth.html liefern
// die Inhalte als {{define "content-..."}}-Bloecke.
var tpl = template.Must(template.ParseFS(webFiles, "web/layout.html", "web/home.html", "web/auth.html"))

type pageData struct {
	LoggedIn   bool
	Email      string
	Error      string
	CSRFToken  string
	ShowResend bool
	Token      string
}

func render(w http.ResponseWriter, r *http.Request, store *Store, contentName string, data pageData) {
	data.CSRFToken = csrfTokenFor(w, r, store)
	clone, err := tpl.Clone()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	clone = template.Must(clone.Parse(`{{define "content"}}{{template "` + contentName + `" .}}{{end}}`))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := clone.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("render error: %v", err)
	}
}
