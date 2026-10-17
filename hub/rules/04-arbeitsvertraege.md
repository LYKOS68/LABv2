# 04 — Arbeitsverträge: wer bedient was

Macht jede ausführbare Handlung ohne eigene Programmierkenntnisse bedienbar. Roman schreibt keinen Code — das übernimmt in dieser Zusammenarbeit durchgehend die KI (bestätigt in `roman_code_review_instructions.md`: seine reale Lücke ist „Code selbst schreiben", nicht Konzepte verstehen).

## Wer bearbeitet was

| Bereich | Bearbeitet von |
|---|---|
| Alle `.go`-Dateien (`cmd/lab-studio/*.go`) | **ausschließlich KI** |
| Alle Markdown-Dokumente (`hub/rules/*.md`, `README.md`) | **ausschließlich KI** |
| `.vscode/tasks.json`, `.gitignore` | **ausschließlich KI** |
| `go.mod`/`go.sum` | **ausschließlich die `go`-CLI** (über `go get`/`go mod tidy`, von der KI ausgeführt) — nie von Hand |
| Visuelle Freigabe der Homepage im Browser | **ausschließlich Roman** — nur er kann beurteilen, ob es „richtig aussieht" |
| Entscheidung über Struktur/Umfang/Fertigstellung | **ausschließlich Roman** (§3.1 zusammenarbeit-gesetze.md) |

## Was die KI automatisch ausführt (keine Freigabe nötig)

| Prüfung | Befehl | Wann |
|---|---|---|
| Kompilier-Check | `go build ./cmd/lab-studio` | nach jeder Code-Änderung |
| Statische Analyse | `go vet ./...` | nach jeder Code-Änderung |
| Unit-Tests | `go test ./...` | nach jeder Code-Änderung |
| Funktionaler Smoke-Test (curl-Sequenz) | — | nach jedem abgeschlossenen Schritt |

Das sind reine Lesevorgänge bzw. reversible lokale Testläufe gegen Wegwerf-Datenbanken — kein Freigabebedarf.

## Was ausdrückliche Freigabe braucht (nie automatisch)

| Aktion | Warum |
|---|---|
| `git commit` | verändert dauerhafte Historie |
| `git push` | verändert geteilten/entfernten Zustand |
| Jede Strukturentscheidung (neue Datei außerhalb bestehender Konventionen, neues Modul) | braucht BAUSTEIN-Protokollzeile + Romans Entscheidung bei fehlendem a/b/c-Auslöser |
| Löschen von Dateien/Ordnern, die die KI nicht selbst in dieser Session erzeugt hat | z. B. `manus-packet` — nur auf explizite Bestätigung |
| Erweiterung des Funktionsumfangs über die aktuelle Schritt-Definition hinaus | §3.4: aus einer Aufgabe entsteht nie ungefragt eine zweite |

## Was Roman selbst im Terminal/Browser starten muss

Die KI kann den Server zu Testzwecken selbst kurz starten und wieder stoppen (passiert laufend in dieser Session). Für Romans eigene Nutzung/Demo/Abnahme muss er selbst aktiv werden — dafür gibt es keinen Ersatz:

| Aktion | Befehl / Ort | Voraussetzungen | Erwartetes Ergebnis | Fehlerverhalten | Rücksetzung |
|---|---|---|---|---|---|
| Server lokal starten | `go run ./cmd/lab-studio` in der Projektwurzel, **oder** VS Code: Terminal → Run Task → „RAP-LAB: Server starten" | Go 1.27+ installiert, Port 8090 frei | Konsole zeigt `RAP-LAB laeuft auf http://127.0.0.1:8090` | Fehlermeldung + Programmabbruch, z. B. bei fehlendem Schreibrecht auf `data/` oder belegtem Port | `Strg+C` im Terminal, keine Nebenwirkungen |
| Seite ansehen/prüfen | Browser → `http://127.0.0.1:8090` | Server läuft (siehe oben) | Homepage lädt | Verbindung abgelehnt → Server läuft nicht | — |
| Tests selbst nachvollziehen | `go test ./... -v` in der Projektwurzel | — | `PASS` je Test | `FAIL` mit Zeilenangabe | — |
| Commit erstellen/pushen | `git commit` / `git push` | — | — | — | — |

## Trennung automatisch vs. manuell — Zusammenfassung

```
Automatisch (KI, keine Freigabe):  build → vet → test → Smoke-Test
Manuell (Roman, keine KI-Ausführung): Server starten, Browser pruefen, git push, Loeschungen ausserhalb der Session, Strukturfreigaben
```
