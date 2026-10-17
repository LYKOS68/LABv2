# Verzeichnisvertrag — RAP-LAB (LABv2)

Verbindliche Pfad- und Schreibregeln für dieses Repository. Jede neue Datei wird nach dieser Tabelle zugeordnet, nicht nach Gefühl.

| Bereich | Pfad | Zweck | Git |
|---|---|---|---|
| Go-Code | `cmd/<tool>/` | Einstiegspunkte, ein Unterordner pro ausführbarem Programm | versioniert |
| Frontend | `web/` | ausgelagerte Bild-/HTML-/CSS-Assets | versioniert |
| Konfiguration | `configs/` | handgepflegte Settings, Env-Vorlagen | versioniert |
| Runtime-Daten + maschinenlesbare Artefakte | `data/` | DB-Dateien, JSON-Exports, alles zur Laufzeit erzeugt oder importiert | **gitignored** |
| Skripte | `scripts/` | Wartungs-/Automatisierungsskripte | versioniert |
| Logs | `logs/` | Laufzeit-Logs | **gitignored** |
| Prüfprotokolle | `verification/` | Check-Ergebnisse, Smoke-Test-Records | versioniert |
| Verträge/Governance | `docs/` | Spezifikationen, dieser Pfadvertrag | versioniert |
| Startsteuerung | `.vscode/tasks.json` | VS-Code-Task für Serverstart | versioniert |

## Begründete Abweichungen von der wörtlichen Vorgabe

**Historie** bekommt keinen eigenen Ordner. Git (`.git/`) erfüllt diese Zuständigkeit bereits — ein zusätzlicher `history/`-Ordner wäre eine doppelte Zuständigkeit (ausdrücklich zu vermeiden laut Auftrag).

**Tests** liegen Go-idiomatisch als `*_test.go` neben dem jeweiligen Code (z. B. künftig `cmd/lab-studio/main_test.go`), nicht in einem separaten Top-Level-`tests/`-Ordner. `go test ./...` erwartet Tests dort — ein eigener Ordner würde diese Konvention brechen.

## Protokoll dieser Strukturentscheidung

```
BAUSTEIN: web/, configs/, data/, scripts/, logs/, verification/, docs/ | AUSLÖSER: keiner von a/b/c (§6.3 zusammenarbeit-gesetze.md) | BELEG: explizite Entscheidung von Roman, "vollständigen 12-Bereiche-Vertrag jetzt anlegen" | INVARIANTE: jede neue Datei liegt gemäß obiger Tabelle, keine Kategorie wird doppelt zugeordnet
```

Diese Ordner enthalten aktuell keinen Inhalt außer `docs/` und `cmd/`. Das ist eine bewusste Vorrats-Entscheidung, keine regelkonforme Notwendigkeit nach §6.3 — festgehalten hier, damit das bei einer späteren Prüfung nachvollziehbar bleibt.

## Innerer Pfadvertrag: `cmd/lab-studio/`

Innerhalb des Go-Codes gilt ein zweiter, feinerer Vertrag — jede Datei gehört zu genau einer Funktionsachse, alle im selben Package `main` (keine Unterpakete, kein `internal/`):

| Datei | Achse | Inhalt |
|---|---|---|
| `main.go` | Wiring | nur `func main()` — Store/Limiter erzeugen, Routen registrieren, Server starten. Kein eigener Logik-Code (Kleber-Prinzip, §5b ARBEITSWEISE.md) |
| `handlers.go` | HTTP-Routing | `currentUser` + die fünf Route-Handler (`homeHandler`, `registerHandler`, `loginHandler`, `logoutHandler`, `deleteAccountHandler`) |
| `store.go` | Persistenz | `User`, `Store`, alle SQL-Zugriffe. **Einzige Datei, die `*sql.DB` berührt** |
| `password.go` | Passwort-Hashing | Argon2id-Parameter, `hashPassword`, `verifyPassword`. Keine HTTP-/DB-Abhängigkeit |
| `csrf.go` | CSRF-Schutz | `csrfEntry`, CSRF-Store-Methoden, `csrfTokenFor`, `checkCSRF`. Vollständige Audit-Fläche für diese Invariante in einer Datei |
| `ratelimit.go` | Rate-Limiting | `RateLimiter` und zugehörige Funktionen. Keine Abhängigkeit zu Store/Templates |
| `validate.go` | Eingabevalidierung | `validateRegistration`, `emailRe`. Reine Funktionen, kein I/O |
| `templates.go` | Präsentation | CSS, HTML-Templates, `pageData`, `render` |
| `main_test.go` | Tests | Go-idiomatisch neben dem Code, s. o. |

Neue SQL-Query → gehört in `store.go`. Neues Formularfeld → `templates.go` + zugehöriger Handler. Neue Route → `handlers.go` + Registrierung in `main.go`. Diese Aufteilung ist reine Dateiorganisation innerhalb eines Package — kein neues Modul, kein neues Interface, daher kein BAUSTEIN-Protokoll nach §6.3 nötig (der Vertrag selbst dokumentiert die Achsen-Begründung, s. `docs/05-main-go-zerlegung.md`).
