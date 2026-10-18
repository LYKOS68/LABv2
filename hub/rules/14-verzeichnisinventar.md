# 14 — Verzeichnisinventar (löst `03-verzeichnisinventar.md` ab)

Stand 2026-10-18, gegen den realen Bestand geprüft (`git ls-files`, `wc -l`, `grep -c "^func Test"`). `03` bleibt als Stand vom 2026-10-17 unverändert stehen.

Für alle Zeilen gilt, wenn nicht anders vermerkt: **Eigentümer** Roman (entscheidet über Inhalt und Umfang), **Schreiber** KI (Roman schreibt keinen Code, `04-arbeitsvertraege.md`).

## Programm — `cmd/lab-studio/` (ein Paket `main`)

| Datei | Zeilen | Zweck | Prüfverfahren |
|---|---|---|---|
| `main.go` | 46 | `main()` + `newMux` (alle Routen an einer Stelle) | `go build`, Tests nutzen `newMux` |
| `handlers.go` | 262 | Startseite, Konto: Registrieren, Bestätigen, Anmelden, Passwort vergessen/zurücksetzen, Abmelden, Löschen | `main_test.go` (Ende-zu-Ende über HTTP) |
| `studio.go` | 325 | LAB-STUDIO: Seite + JSON-Schnittstelle, Sitzungs- und CSRF-Prüfung | `studio_test.go` |
| `connect.go` | 11 | LAB-CONNECT-Seite (Oberflächen-Prototyp) | `connect_test.go` |
| `store.go` | 487 | **einzige Datei mit SQL**: Konten, Sitzungen, Einmal-Links, Songs, Studio-Felder/-Ansicht | Tests inkl. Neustart und Konto-Trennung |
| `csrf.go` | 96 | CSRF-Token | Always-Pass-Test (403 ohne/mit falschem Token) |
| `ratelimit.go` | 78 | Begrenzung der Versuche pro IP | `main_test.go` |
| `password.go` | 64 | Argon2id | Rundlauf + Korrupt-Hash-Tests |
| `mail.go` | 66 | Mail-Versand (SMTP aus Umgebungsvariablen, sonst Konsole) | Ende-zu-Ende-Tests lesen Links aus dem Log |
| `validate.go` | 26 | E-Mail-/Passwort-Regeln | Tests mit ungültigen Eingaben |
| `templates.go` | 209 | HTML/CSS von Startseite und Konto-Seiten | Tests + Screenshot-Sichtprüfung |
| `main_test.go` | 476 | 19 Tests: Konto, Links, Startseite, Navigation, Hilfs-„Test-Browser" | `go test ./...` |
| `studio_test.go` | 350 | 14 Tests LAB-STUDIO | `go test ./...` |
| `connect_test.go` | 47 | 4 Tests LAB-CONNECT + Schrift | `go test ./...` |

### `cmd/lab-studio/web/` — per `go:embed` ins Programm eingebaut

| Datei | Zweck |
|---|---|
| `studio.html`, `connect.html` | Seitengerüste LAB-STUDIO, LAB-CONNECT |
| `static/brand.css` | Markenfarben + Manrope für **alle** Seiten |
| `static/studio.css`, `static/studio.js` | LAB-STUDIO Aussehen / Verhalten |
| `static/connect.css`, `static/connect.js` | LAB-CONNECT Aussehen / Verhalten (Beispieldaten, speichert nichts) |
| `static/img/krone.jpg` | Kopf-Icon Startseite ← `hub/brand/originale/krone-logo-second.jpeg` |
| `static/img/schriftzug.png` | Schriftzug Startseite ← `rap-lab-schriftzug.png` |
| `static/img/leder-hintergrund.jpg` | blasser Hintergrund Startseite ← `unterschrift-ledercover2.jpeg` |
| `static/img/wordmark-header.jpg` | Kopf-Logo LAB-CONNECT, Ausschnitt ← `unterschrift-ledercover2.jpeg` |
| `static/fonts/manrope-latin-{400,500,600,700}-normal.woff2`, `manrope-OFL.txt` | Schrift + Lizenz (aus `@fontsource/manrope` 5.3.0, identisch zum Original-Prototyp) |

## Laufzeitdaten (nicht auf GitHub)

| Pfad | Schreiber | Zweck |
|---|---|---|
| `data/raplab.sqlite` (+ `-wal`, `-shm`) | **der laufende Server** | Konten, Links, Songs, Studio-Felder. Entsteht beim ersten Start |
| `logs/` | — | **LÜCKE**: keine Schreibanbindung, Server schreibt nur ins Konsolenfenster |

## Werkzeuge und Metadaten

| Datei | Schreiber | Zweck |
|---|---|---|
| `go.mod`, `go.sum` | `go`-Befehl, nie von Hand | Abhängigkeiten (einzige direkte: `modernc.org/sqlite`, `golang.org/x/crypto`) |
| `.github/workflows/test.yml` | KI | GitHub prüft bei jedem Push: `gofmt`, `go vet`, `go test` |
| `.vscode/tasks.json` | KI | VS-Code-Aufgabe „RAP-LAB: Server starten" |
| `.claude/settings.json` | Claude Code | Einstellungen des KI-Werkzeugs |
| `.gitignore` | KI | hält `data/`, `logs/`, `*.exe` aus Git heraus |
| `README.md` | KI | Einstieg: was die Seite kann, Start, Stand, Offenes |

## `hub/` — alles, was kein Programm ist

| Pfad | Zweck |
|---|---|
| `hub/README.md` | Einstieg in die zweite Wurzel |
| `hub/agent-docs/00-orientierung.md` | Betriebsregeln + aktueller Stand für jede KI |
| `hub/chatroom.jsonl`, `chatroom.schema.json` | Nachrichten zwischen Agents und Roman (append-only, derzeit leer) |
| `hub/drafts/06-funktions-ausblick-draft.md` | Kleinanzeigen-Mechaniken, **kein Bauauftrag** |
| `hub/brand/originale/` | 11 Marken-Originale in voller Auflösung — **nie verändern** |
| `hub/brand/referenzen/` | 7 Design-Vorlagen (Startseite, LAB-CONNECT-Screenshots, Original-Prototyp-Vorschauen) |

### `hub/rules/` — nummeriert, wird nicht überschrieben, nur abgelöst

| Nr. | Inhalt | Status |
|---|---|---|
| `01` | Startseiten-Spezifikation | LÜCKEN 1–4 abgelöst durch `12` |
| `02` | Sicherheit: SQLite, Argon2id, CSRF, Rate-Limiting | umgesetzt |
| `03` | Verzeichnisinventar 2026-10-17 | abgelöst durch `14` |
| `04` | Arbeitsverträge: wer bedient was | gültig |
| `05` | Zerlegung von `main.go` | umgesetzt |
| `07` | E-Mail-Bestätigung | umgesetzt |
| `08` | Passwort vergessen | umgesetzt |
| `09` | LAB-STUDIO MVP | umgesetzt |
| `10` | Startseite im Markendesign | umgesetzt, Roman: „noch nicht fertig" |
| `11` | LAB-CONNECT Etappe 1 (Oberfläche, an Original angeglichen) | umgesetzt als Oberfläche |
| `12` | Brand-Spec | Entwurf, offene Punkte §6 |
| `13` | Zweck und Vision | Entwurf, offene Fragen §7 |
| `14` | dieses Inventar | gültig |
| `CONTRACT.md` | Pfadvertrag (zwei Wurzeln, Dateien je Achse) | gültig |

(`06` liegt bewusst in `hub/drafts/`.)

## Einstiegspunkt

Einziger: `cmd/lab-studio/main.go`, `func main()`. Start: `go run ./cmd/lab-studio` → `http://127.0.0.1:8090`.
