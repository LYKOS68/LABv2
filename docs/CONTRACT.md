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
