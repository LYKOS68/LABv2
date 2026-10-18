package main

import "net/http"

// LAB-CONNECT ist vorerst ein reiner Oberflaechen-Prototyp (hub/rules/11-lab-connect-mvp-plan.md):
// die Seite laeuft komplett im Browser mit Beispielanzeigen und schickt nichts an den Server.
func connectPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, webFiles, "web/connect.html")
	}
}
