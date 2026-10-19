package server

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

// Das Inserieren-Feld hat einen Schliessen-Knopf und die Preisarten des Original-Prototyps.
func TestConnectListingFormBasics(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	page := web.get("/connect")
	expectContains(t, page, "data-close-sheet")
	expectContains(t, page, `name="amount"`)
	for _, priceType := range []string{"Festpreis", "Ab-Preis", "Verhandelbar", "Auf Anfrage", "Zusammenarbeit"} {
		expectContains(t, page, ">"+priceType+"</option>")
	}
	expectMissing(t, page, "PLATZHALTER")
}

// Manrope gilt fuer die ganze Seite und wird lokal ausgeliefert, samt Lizenz.
func TestManropeFontIsServedLocally(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	expectContains(t, web.get("/static/brand.css"), `font-family: "Manrope"`)
	expectContains(t, web.get("/static/fonts/manrope-OFL.txt"), "SIL Open Font License")
	for _, weight := range []string{"400", "500", "600", "700"} {
		resp, err := web.client.Get(web.base + "/static/fonts/manrope-latin-" + weight + "-normal.woff2")
		if err != nil {
			t.Fatalf("GET font %s: %v", weight, err)
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("font %s: status %d", weight, resp.StatusCode)
		}
	}
}
