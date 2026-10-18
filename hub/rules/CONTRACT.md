# Verzeichnisvertrag — RAP-LAB (LABv2)

Zwei Wurzeln. Vier Regeln. Kein Ordner existiert auf Vorrat — er entsteht erst, wenn echter Inhalt da ist.

## Die zwei Wurzeln

```
LABv2/              ← erste Wurzel: Runtime/Code (was die Website baut/startet/schreibt)
└── hub/            ← zweite Wurzel: alles Nicht-Runtime (Regeln, Wissen, Skills, Agent-Koordination)
```

## Die vier Regeln

1. **Code, der die Website baut/startet** → `cmd/<tool>/`, flach, nach Funktionsachse getrennt, ein Package. Kein `internal/`, bis ein zweites Binary real davon importiert.
2. **Was der laufende Prozess selbst schreibt** (DB, Logs) → `data/` bzw. `logs/`, gitignored, nie in `cmd/`.
3. **Eingefrorene, geprüfte Entscheidungen** (Specs, Pläne, Verträge — dieser hier eingeschlossen) → `hub/rules/`, nummeriert, nie überschrieben, nur durch neues Dokument abgelöst.
4. **Unfertiges Denken + Agent-Koordination** → `hub/`, niemals von sich aus ein Bauauftrag — siehe `hub/README.md`.

Statische Assets (`web/`), handgepflegte Settings (`configs/`), Skripte, die Roman selbst triggert (`scripts/`) — alle nach Regel 1 bzw. 2, sobald der erste echte Inhalt ansteht. Bis dahin existiert der Ordner nicht. (Vier leere Platzhalter-Ordner dieser Art wurden am 2026-10-17 wieder entfernt — reine Vorrats-Entscheidung ohne Beleg, s. zusammenarbeit-gesetze.md §6.3.)

## Begründete Abweichungen von der wörtlichen Vorgabe

**Historie** bekommt keinen eigenen Ordner. Git (`.git/`) erfüllt diese Zuständigkeit bereits.

**Automatische Prüfung** liegt in `.github/workflows/test.yml` — Pfad von GitHub vorgegeben, nicht frei wählbar. Läuft bei jedem Push: `gofmt`, `go vet`, `go test`.

**Tests** liegen Go-idiomatisch als `*_test.go` neben dem Code (`cmd/lab-studio/main_test.go`), nicht in einem separaten Top-Level-`tests/`-Ordner — `go test ./...` erwartet sie dort.

## `hub/rules/` vs. der Rest von `hub/`

`hub/rules/` = eingefroren. Jede Datei beschreibt einen geprüften, abgeschlossenen Stand. Änderung nur durch neuen, benannten Schritt — nicht durch stilles Überschreiben.

`hub/drafts/`, `hub/chatroom.jsonl`, `hub/agent-docs/` = nicht eingefroren. Insbesondere `hub/drafts/`: **kein** Bauauftrag, auch wenn es wie einer aussieht. Details in `hub/README.md`.

## Innerer Pfadvertrag: `cmd/lab-studio/`

Jede Datei gehört zu genau einer Funktionsachse, alle im selben Package `main`:

| Datei | Achse | Inhalt |
|---|---|---|
| `main.go` | Wiring | `func main()` (Store/Mailer erzeugen, Server starten) + `newMux` (alle Routen an einer Stelle, auch von den Tests genutzt) — Kleber-Prinzip |
| `handlers.go` | HTTP-Routing | `currentUser` + die fünf Route-Handler |
| `store.go` | Persistenz | `User`, `Store`, alle SQL-Zugriffe — **einzige Datei, die `*sql.DB` berührt** |
| `password.go` | Passwort-Hashing | Argon2id-Parameter, `hashPassword`, `verifyPassword` |
| `csrf.go` | CSRF-Schutz | vollständige Audit-Fläche für diese Invariante in einer Datei |
| `ratelimit.go` | Rate-Limiting | `RateLimiter` und zugehörige Funktionen |
| `validate.go` | Eingabevalidierung | `validateRegistration`, `emailRe` |
| `studio.go` | LAB-STUDIO (HTTP) | Seite `/studio`, JSON-API `/studio/api/…`, Sitzungs- und CSRF-Prüfung (`studioAPI`) — s. `09-lab-studio-mvp-plan.md` |
| `web/` | Präsentation (Dateien) | per `go:embed` eingebaut: `studio.html` (Vorlage), `static/` unter `/static/` (`brand.css` = Markenfarben der ganzen Seite, `studio.css`, `studio.js`). Liegt unter `cmd/lab-studio/`, weil `go:embed` nur Unterordner des Pakets erlaubt |
| `studio_test.go` | Tests | LAB-STUDIO, Go-idiomatisch neben `studio.go` |
| `mail.go` | E-Mail-Versand | `Mailer` (SMTP aus Umgebungsvariablen, sonst Link in die Konsole) — s. `07-email-verifikation-plan.md` |
| `templates.go` | Präsentation | CSS, HTML-Templates, `pageData`, `render` |
| `main_test.go` | Tests | Go-idiomatisch neben dem Code |

Neue SQL-Query → `store.go`. Neues Formularfeld → `templates.go` + Handler. Neue Route → `handlers.go` + Registrierung in `main.go`. Begründung je Achse (mit Vorher-Beleg): `hub/rules/05-main-go-zerlegung.md`.
