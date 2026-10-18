package main

import "testing"

func TestNavLinksToConnect(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	expectContains(t, web.get("/"), `href="/connect">LAB-CONNECT<`)
}

// LAB-CONNECT ist ein Oberflaechen-Prototyp: ohne Anmeldung erreichbar, laedt seine Dateien.
func TestConnectPageLoadsWithoutLogin(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	page := web.get("/connect")
	expectContains(t, page, "Was fehlt deinem Projekt?")
	expectContains(t, page, `<script src="/static/connect.js">`)
	expectContains(t, web.get("/static/connect.js"), "function submitListing")
	expectContains(t, web.get("/static/connect.css"), ".sheet")
}

// Das Inserieren-Feld hat einen Schliessen-Knopf und eine Preis-Einheit als Platzhalter.
func TestConnectListingFormBasics(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	page := web.get("/connect")
	expectContains(t, page, "data-close-sheet")
	expectContains(t, page, `name="price"`)
	expectContains(t, page, "<option>PLATZHALTER</option>")
}
