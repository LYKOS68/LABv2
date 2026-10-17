# 03 — Verzeichnisinventar: Zweck, Verwendung, Verantwortung

Vollständiges, wahrheitsgetreues Inventar der RAP-LAB-Verzeichnisstruktur, Stand 2026-10-17. Jede Zeile ist gegen den realen Dateibestand geprüft (`find`, `wc -l`, `git log`), nicht aus dem Gedächtnis geschrieben.

**Spaltenbedeutung:**
- **Eigentümer** — wer über Inhalt/Scope entscheidet (§3.1 zusammenarbeit-gesetze.md: immer Roman)
- **Schreiber** — wer die Datei tatsächlich erzeugt/ändert
- **Leser** — wer/was die Datei zur Entwicklungs- oder Laufzeit liest
- **Auslöser** — was einen Schreibvorgang auslöst
- **Zweck** — wofür die Datei da ist
- **Prüfverfahren** — wie man feststellt, dass der Inhalt korrekt ist

## Code

**Stand nach Funktionsachsen-Zerlegung** (s. `05-main-go-zerlegung.md`): `main.go` ist auf reines Wiring reduziert (25 Zeilen), die Logik liegt in sieben nach Achse getrennten Dateien, alle im selben Package `main`.

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `cmd/lab-studio/main.go` | Roman | KI | `go build`/`run`, Roman (Review) | neue Route/neuer Start-Parameter | Wiring: Store/Limiter erzeugen, Routen registrieren, Server starten | `go build`, manueller Start |
| `cmd/lab-studio/handlers.go` | Roman | KI | `go build`, Roman (Review) | neue Route | HTTP-Handler für alle fünf Routen + `currentUser` | Smoke-Test (curl-Sequenz) |
| `cmd/lab-studio/store.go` | Roman | KI | `go build`/`test`, Roman (Review) | neue Datenanforderung | SQLite-Persistenz, einzige Datei mit SQL | Neustart-Test, `go test` |
| `cmd/lab-studio/password.go` | Roman | KI | `go test`, Roman (Review) | Sicherheitsanforderung | Argon2id-Hashing | `go test ./... -v` (Rundlauf + Korrupt-Hash) |
| `cmd/lab-studio/csrf.go` | Roman | KI | `go build`, Roman (Review) | CSRF-Lücke (Plan Schritt 3) | CSRF-Token-Lebenszyklus + Validierung | Always-Pass-Test (403 ohne/mit falschem Token) |
| `cmd/lab-studio/ratelimit.go` | Roman | KI | `go test`, Roman (Review) | Brute-Force-Lücke (Plan Schritt 4) | In-Memory-Rate-Limiter | `go test ./... -v` + 6-Versuche-Live-Test |
| `cmd/lab-studio/validate.go` | Roman | KI | `go build`, Roman (Review) | Eingabevalidierung | E-Mail-/Passwort-Regeln | manueller INVALID-Test |
| `cmd/lab-studio/templates.go` | Roman | KI | `go build`, Roman (Review) | neues Formularfeld/neue Sektion | HTML/CSS, Rendering | Browser-/curl-Abgleich gegen `01-spezifikation.md` |
| `cmd/lab-studio/main_test.go` | Roman | KI | `go test` | neue Funktion mit Prüfbedarf | automatisierte Tests: Argon2id-Rundlauf, Korrupt-Hash, Rate-Limiter-Grenzen | `go test ./... -v` → PASS |

## Modul-Metadaten

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `go.mod` | Roman | **`go`-CLI selbst** (`go get`/`go mod init`), nie von Hand editiert | `go build`/`run`/`test` | neue Dependency | Modulname, Go-Version, direkte Abhängigkeiten | `go build ./...` schlägt fehl, wenn inkonsistent |
| `go.sum` | Roman | `go`-CLI selbst | `go build`/`run`/`test` | wie `go.mod` | Checksummen aller (auch transitiven) Abhängigkeiten | dito |

## Governance / Dokumentation

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `README.md` | Roman | KI | Roman, GitHub-Weboberfläche | jeder abgeschlossene Schritt | Einstiegspunkt: Startbefehl, geprüfter Stand, offene LÜCKEN | Startbefehl tatsächlich ausführen — stimmt er mit der realen Route überein? |
| `hub/rules/01-spezifikation.md` | Roman | KI | Roman, KI (Frontend-Referenz) | Referenzbild-Analyse (erste Kette, Schritt 1/5) | Homepage-Spezifikation, LÜCKEN-Liste zu Assets/Farben/Fonts | Abgleich gegen `content-home`-Template in `main.go` |
| `hub/rules/CONTRACT.md` | Roman | KI | Roman, KI (vor jeder neuen Datei) | Strukturentscheidung (zweite Kette) | verbindlicher Pfadvertrag für alle Verzeichnisse | jede neue Datei gegen diese Tabelle prüfen |
| `hub/rules/02-sicherheitserweiterungen-plan.md` | Roman | KI | Roman, KI | Sicherheitsaudit nach Manus-Fragment-Review | 4-Schritte-Plan (Persistenz/Argon2id/CSRF/Rate-Limit) + Nachweis je Schritt | genannte Commit-Hashes in `git log` auffindbar? |
| `hub/rules/03-verzeichnisinventar.md` (diese Datei) | Roman | KI | Roman, KI | diese Kette, Schritt 1/5 | vollständiges Datei-für-Datei-Inventar | Abgleich gegen `find .` |

## Startsteuerung / Werkzeugkonfiguration

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `.vscode/tasks.json` | Roman | KI | VS Code | Startkomfort ohne Terminal-Tippen | kapselt `go run ./cmd/lab-studio` als Task | Terminal → Run Task → Server startet, Browser zeigt Seite |
| `.gitignore` | Roman | KI | `git` | Runtime-Dateien sollen nicht versioniert werden | schließt `data/`, `logs/`, `*.exe` von Commits aus | `git status` nach Serverlauf zeigt diese Pfade nicht als untracked |

## Noch nicht existierende Ordner (Regel, kein Platzhalter)

`web/`, `configs/`, `scripts/` existieren aktuell nicht — sie wurden am 2026-10-17 als leere Platzhalter wieder entfernt (keine Nutzung seit Anlage, kein a/b/c-Auslöser). Nach dem neuen Vertrag (`CONTRACT.md`, Regel 1/2) entstehen sie erst mit ihrem ersten echten Inhalt.

| Ordner | Status |
|---|---|
| `logs/` | existiert, aber **LÜCKE**: `main.go` schreibt aktuell nicht hierhin — nur `log.Println` nach stdout. Füllt sich nur bei manueller Umleitung (`> logs/x.log`) |

## Laufzeit-Daten (nicht versioniert, nur lokal)

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `data/raplab.sqlite` (+ `-wal`/`-shm`) | Roman | **der laufende Go-Prozess selbst**, weder KI noch manuell | Go-Prozess (SQLite-Treiber) | Server-Start + jede Registrierung/Login/Löschung | persistente Nutzerdaten | Neustart-Test: Login nach Neustart funktioniert noch? |

## Entry Point

Einziger Einstiegspunkt: `cmd/lab-studio/main.go`, `func main()`. Startbefehl: `go run ./cmd/lab-studio` (Projektwurzel). Kein zweites `cmd/`-Unterverzeichnis, kein zweites CLI-Tool — bestätigt durch `grep -rl "func main" --include="*.go" .` in der vorherigen Kette.

## Zusammenfassung der LÜCKEN

1. `logs/` hat keine tatsächliche Schreibanbindung im Code; der Vertrag in `CONTRACT.md` beschreibt eine Absicht, keine Implementierung.
2. Kein automatisiertes Log-Rotations- oder Aufräumverfahren für `data/raplab.sqlite` — manuell zu beobachten.
