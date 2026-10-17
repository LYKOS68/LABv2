# 05 — Funktionale Zerlegung von main.go

**Auslöser:** `main.go` war auf 644 Zeilen gewachsen (nach vier Sicherheits-/Persistenz-Erweiterungen, s. `02-sicherheitserweiterungen-plan.md`). Diese Kette fordert, Dateien erst nach belegter Begründung anzulegen und vorzeitige Modulbildung zu vermeiden — deshalb hier die Begründung je Achse, nicht nur das Ergebnis.

**Wichtig:** Das ist reine Dateiaufteilung **innerhalb desselben Package** `main` (`cmd/lab-studio/`). Keine neuen Typen, Interfaces oder Unterpakete entstanden — deshalb greift §6.3 (zusammenarbeit-gesetze.md, "notwendig: neuer Baustein/Modul/Abstraktion") hier nicht im engeren Sinne. Trotzdem wird jede Achse mit Beleg (Datei:Zeile im alten Stand) dokumentiert, wie im Auftrag verlangt.

## Achsen und Beleg

| Datei | Achse | Beleg (vorheriger Stand) | Begründung |
|---|---|---|---|
| `password.go` | Passwort-Hashing | `main.go:98-148` | reine Funktionen, keine Abhängigkeit zu HTTP/DB — eigenständig testbar, bereits eigene Unit-Tests vorhanden |
| `ratelimit.go` | Rate-Limiting | `main.go:27-96` | vollständig self-contained Subsystem, keine Abhängigkeit zu `Store`/Templates |
| `store.go` | Persistenz | `main.go:150-325` | einzige Datei, die `*sql.DB` berührt — isoliert die SQL-Injection-Audit-Fläche auf eine Datei |
| `csrf.go` | CSRF-Schutz | `main.go:155, 278-317, 449-488` | sicherheitskritisches Subsystem, von Roman als eigener Plan-Schritt geprüft (`02-...-plan.md` Schritt 3) — vollständige Audit-Fläche in einer Datei |
| `validate.go` | Eingabevalidierung | `main.go:326-340` | reine Validierung ohne I/O |
| `templates.go` | Präsentation | `main.go:342-515` | einzige Datei mit HTML/CSS-Strings |
| `handlers.go` | HTTP-Routing | `main.go:516-644` (vorher anonyme Closures in `main()`) | Routing-Logik gehört laut Kleber-Prinzip (§5b ARBEITSWEISE.md) nicht in `main()` |
| `main.go` (neu) | Wiring | — | bleibt: Store/Limiter erzeugen, Routen registrieren, Server starten — sonst nichts |

## Verbote eingehalten

- Keine neuen Typen/Interfaces eingeführt (Closure-Factory-Pattern für Handler, kein neuer `app`-Struct).
- Kein neues Package, kein `internal/`.
- Keine Verhaltensänderung — einzige Korrektur während der Zerlegung: das `log.Printf("delete-account error: …")` wurde zunächst versehentlich weggelassen und direkt wieder ergänzt, bevor committed wurde.

## Prüfung (Schritt 4/5 der Kette)

- `go build ./cmd/lab-studio`: OK
- `gofmt -l cmd/lab-studio/`: keine Ausgabe (bereits korrekt formatiert)
- `go vet ./...`: OK
- `go test ./... -v`: alle 4 Tests PASS
- 12-Punkte-Live-Smoke-Test gegen frisch gebautes Binary (nicht gegen einen alten Prozess — das wurde über Prozess-/Zeitstempel-Vergleich verifiziert, siehe Session): Homepage, CSRF-Gate (403 ohne Token), Register, Duplikat, Login falsch/korrekt, eingeloggte Homepage, Logout, Persistenz-Neustart, Rate-Limiting (5. Versuch blockiert), Account-Löschung, Login nach Löschung schlägt fehl — alle wie vor dem Umbau.

**Ergebnis:** `main.go` 644 → 25 Zeilen. Verhalten nachweisbar unverändert.
