package server

import (
	"net/http"

	"rap-lab/cmd/lab-studio/web"
)

// LAB-CONNECT ist vorerst ein reiner Oberflaechen-Prototyp (hub/rules/11-lab-connect-mvp-plan.md):
// die Seite laeuft komplett im Browser mit Beispielanzeigen und schickt nichts an den Server.
func connectPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web.Files, "connect.html")
	}
}
