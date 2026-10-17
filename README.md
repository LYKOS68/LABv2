# LABv2

RAP-LAB — Go-Fullstack-Homepage mit Registrierung/Login/Logout/Account-Löschen.

## Start

```
go run ./cmd/lab-studio
```

Läuft auf `http://127.0.0.1:8090`. Alternativ VS Code: Terminal → Run Task → „RAP-LAB: Server starten".

## Geprüfter Stand

- Homepage (Text-Platzhalter statt Bild-Assets), Registrierung, Login, Logout, Account-Löschung — funktional getestet.
- Passwörter via Argon2id gehasht (PHC-Format, OWASP-Mindestparameter m=19456 KiB, t=2, p=1), nie im Klartext gespeichert. Korrupte/fremde Hash-Strings führen zu Fehler statt Absturz.
- Nutzerdaten in SQLite (`data/raplab.sqlite`) — **persistent**, bewiesen durch Neustart-Test. Sessions bleiben bewusst in-memory.
- CSRF-Schutz auf allen Mutationen (Register/Login/Logout/Delete-Account): Synchronizer-Token gegen HttpOnly-Cookie, 30 Minuten gültig, zeitkonstanter Vergleich. Fehlendes/falsches Token → 403.
- Rate-Limiting auf Register/Login: max. 5 Versuche pro 10 Minuten pro IP, danach `429` + `Retry-After`. In-Memory, begrenzte Map (max. 10.000 Schlüssel).

Damit sind alle vier Schritte aus `docs/02-sicherheitserweiterungen-plan.md` umgesetzt.

`cmd/lab-studio/main.go` ist auf reines Wiring reduziert (25 Zeilen); die Logik liegt in sieben nach Funktionsachse getrennten Dateien im selben Package — siehe `docs/05-main-go-zerlegung.md` für die Begründung je Datei und `docs/CONTRACT.md` für den vollständigen Pfadvertrag.

## Tests

```
go test ./...
```

Deckt ab: Argon2id-Hash/Verify-Rundlauf, Korrupt-Hash-Handling (kein Panic), Rate-Limiter-Grenzen.

## Offene LÜCKEN

- Siehe `docs/01-spezifikation.md` §4: exakte Hex-Farbwerte, Font-Namen, freigestellte Logo-/Signature-Assets fehlen noch (aktuell Text-Platzhalter statt Bild).
- `web/`, `configs/`, `scripts/`, `verification/` sind angelegt, aber komplett leer (bewusste Vorratsentscheidung, s. `docs/CONTRACT.md`).
- `logs/` hat keine tatsächliche Schreibanbindung im Code (nur `log.Println` nach stdout).

## Verzeichnisstruktur und Zuständigkeiten

- `docs/CONTRACT.md` — verbindlicher Pfadvertrag (Top-Level + interner Vertrag für `cmd/lab-studio/`)
- `docs/03-verzeichnisinventar.md` — vollständiges Datei-für-Datei-Inventar (Eigentümer/Schreiber/Leser/Prüfverfahren)
- `docs/04-arbeitsvertraege.md` — was die KI automatisch ausführt vs. was Roman selbst im Terminal/Browser starten muss
