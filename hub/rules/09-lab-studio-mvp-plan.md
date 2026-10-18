# 09 — Plan: LAB-STUDIO als MVP in die Runtime übernehmen

**Status:** umgesetzt, Browser-Test durch Roman offen. `go vet` OK, `go test ./...` 32/32 PASS (davon 14 Studio: Zugang ohne Anmeldung → Login/401; Schreiben ohne/falsches CSRF-Token → 403, nichts gespeichert; neues Konto leer; Song speichern/laden; Feld-Änderung lässt ID/Typ unverändert; Konto B kann Songs/Felder von A weder lesen noch ändern noch darauf zeigen; ungültige Eingaben → 400, nichts gespeichert; Zoom begrenzt; Daten überleben Neustart; Konto löschen löscht Studio-Daten; Seite lädt ihre Dateien). Gegenprobe: Konto-Bindung (`owner = ?`) aus den 4 ID-Abfragen entfernt → beide Konto-Trennungs-Tests FAIL. `node --check studio.js` OK.

**Auslöser:** Roman-Entscheidung (Session 2026-10-18): Der LAB-STUDIO-Prototyp (`R:\dw\RAP-LAB-windows-version\RAP-LAB-windows-version\main.go`, 685 Zeilen, eigener Server auf Port 3000) wird refaktoriert und als MVP in `cmd/lab-studio` eingebaut. Romans eigene Lyrics waren nur Entwicklungsmaterial — **jeder Nutzer startet mit leerem Studio**, auch Roman.

## Was der Prototyp kann (Ist-Stand, gelesen)

- Endlose Arbeitsfläche: verschieben mit Maus/Finger, zoomen mit Mausrad/Zwei-Finger/Knöpfen (25 %–400 %), Punkt-Raster im Hintergrund
- Linke Leiste: `SONGS`, `+ FELD` (Notiz), `ZOOM +`, `ZOOM −`, `ANSICHT` (zurücksetzen), Zoom-Anzeige
- Felder auf der Fläche, an der Kopfzeile verschiebbar:
  - **Notiz**: Textfeld, speichert automatisch 0,5 s nach dem Tippen
  - **Songs**: Kacheln (Bildplatzhalter + Name) → Klick öffnet den Text im Editor mit `Zurück`, `Speichern` und Strg+S
- Ansicht und Felder überleben Neuladen und Neustart (`data/workspace.json`), Songs liegen als `.txt`-Dateien
- Schreibzugriffe laufen im Browser nacheinander (keine überholten Speichervorgänge); Server ändert eine Kopie und übernimmt nur bei Erfolg
- Farben: Schwarz `#000000`, Weiß `#FFFFFF`, Neongrün `#ACE182`, Gold `#C2B198`, Braun `#543824`

## Was sich beim Einbau ändern muss (und warum)

| Prototyp | MVP in der Runtime | Grund |
|---|---|---|
| ein gemeinsamer Song-Ordner, kein Konto | Songs, Felder und Ansicht **pro Konto** in SQLite | Website mit fremden Nutzern |
| Songs nur aus Seed-Dateien, kein Anlegen | **neu:** `+ SONG` legt einen leeren Song mit Titel an | leeres Studio ohne Seed wäre sonst unbenutzbar |
| Seed aus `__source-docs/lykos-lyrics` | kein Seed | Romans Entscheidung |
| eigener Server, Port 3000, nur Loopback | Route `/studio` im bestehenden Server (8090), nur angemeldet | eine Runtime |
| Schutz nur über `Origin`-Prüfung | Sitzung + CSRF-Token als Header `X-CSRF-Token` bei jedem Schreibzugriff | bestehendes Gesetz `mutation_ausgeführt ⟹ token_gültig` |
| alles in einer `main.go` (Prototyp-Regel E-3) | Aufteilung nach `CONTRACT.md` (s. Dateien) | LABv2-Vertrag gilt, nicht der des Prototyps |
| JavaScript als dichte Einzeiler im Go-String | eigene Datei, lesbar: kurze Funktionen, eine Anweisung pro Zeile | Romans Lesbarkeitsregel |
| Prüfskripte `00_worm` (Python, AST auf `main.go`, Port 3000) | normale Go-Tests in `*_test.go` | Skripte prüfen eine Datei/einen Port, die es so nicht mehr gibt |

## Dateien (neu/geändert)

| Datei | Achse | Inhalt |
|---|---|---|
| `store.go` | Persistenz | Tabellen `songs`, `studio_fields`, `studio_view` + Zugriffe — bleibt einzige Datei mit SQL |
| `studio.go` (neu) | HTTP-Routing Studio | JSON-API unter `/studio/api/…`, Seite `/studio` |
| `web/studio.html`, `web/studio.css`, `web/studio.js` (neu, unter `cmd/lab-studio/`) | Präsentation Studio | per `go:embed` ins Programm eingebaut — erster echter Inhalt für `web/` (CONTRACT Regel 1); `go:embed` erlaubt nur Unterordner des Pakets, daher nicht Top-Level |
| `main.go` | Wiring | Studio-Routen in `newMux` |
| `templates.go` | Präsentation | Nav-Link `LAB-STUDIO` → `/studio` |

## Kriterien

| Punkt | Inhalt |
|---|---|
| Ziel | Angemeldeter Nutzer öffnet `/studio`, legt Songs und Notizen an, schreibt Texte, verschiebt die Fläche — alles ist nach Neuladen und Server-Neustart wieder da |
| Mindestkriterien | alle Prototyp-Funktionen oben + `+ SONG`; neues Konto sieht leere Fläche ohne Songs; nicht angemeldet → `/studio` leitet auf `/login`, API antwortet 401 |
| Qualitätskriterien | jede SQL-Abfrage ist an das Konto gebunden; Text max. 1 MB; Zoom serverseitig auf 0,25–4 begrenzt; unbekannter Feldtyp → 400; Ändern = ganz oder gar nicht (Transaktion) |
| Prüfverfahren | Go-Tests: (1) Konto B kann Song/Feld von Konto A weder lesen noch ändern (404) — **Always-Pass-Falle**: Test muss rot werden, wenn die Konto-Bindung aus der Abfrage entfernt wird; (2) ohne Sitzung 401, ohne/falsches CSRF-Token 403, nichts gespeichert; (3) ungültiger Typ/Größe/Body → 400, nichts gespeichert; (4) Song speichern → neuer Store auf derselben DB-Datei liest denselben Text; (5) neues Konto: leere Liste. Danach Live-Test im Browser durch Roman |
| Fehlerbedingungen | fremde Songs sichtbar; Schreiben ohne Token möglich; Änderung nach Neustart weg; Romans Lyrics tauchen in einem Konto auf |
| Verbote | kein JS-Framework, kein Build-Schritt, keine neue Go-Dependency; nichts aus der Prototyp-Liste „Nicht im ersten Abschnitt" (Reimbuch, Released/Drafts, Relatable-Bars, 1-2-1-2, Multi-User-Canvas) |

## Entscheidungen (Roman, 2026-10-18)

1. **Farben:** Die Prototyp-Farben gelten für die **ganze** Seite. Sie liegen einmal in `web/static/brand.css` und werden von Startseite/Anmeldung und Studio gemeinsam genutzt (zweite reale Verwendung → gemeinsame Datei, §6.3 a).
2. **Löschen/Umbenennen** von Songs und Notizen: **nicht** im MVP, kommt später.
3. **Speichern:** wie im Prototyp — Songtext per Knopf/Strg+S, Notizen automatisch.
4. `+ SONG` sitzt als kleines Eingabefeld „Neuer Song" im SONGS-Fenster über den Kacheln; die linke Leiste bleibt exakt wie im Prototyp.

BAUSTEIN: LAB-STUDIO (Tabellen `songs`/`studio_fields`/`studio_view`, `studio.go`, `web/`) | AUSLÖSER: Roman-Entscheidung außerhalb §6.3 | BELEG: Prototyp `main.go` (685 Z.) läuft getrennt, ohne Konten | INVARIANTE: `song/feld gelesen oder geändert ⟹ song.owner == angemeldeter_nutzer`
