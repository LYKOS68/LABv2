# LABv2

RAP-LAB — Go-Fullstack-Homepage mit Registrierung/Login/Logout/Account-Löschen.

## Start

```
go run ./cmd/lab-studio
```

Läuft auf `http://127.0.0.1:8090`.

## Geprüfter Stand

- Homepage (Text-Platzhalter statt Bild-Assets), Registrierung, Login, Logout, Account-Löschung — funktional getestet (siehe Check-Verlauf in der Session).
- Passwörter via `bcrypt` gehasht, nie im Klartext gespeichert.
- Nutzer-/Session-Daten liegen in-memory — **nicht persistent**, gehen bei Neustart verloren.

## Offene LÜCKEN

- Siehe `01-spezifikation.md` §4: exakte Hex-Farbwerte, Font-Namen, freigestellte Logo-/Signature-Assets fehlen noch (aktuell Text-Platzhalter statt Bild).
- Keine Tests (`go test ./...` → no test files).
- Keine `.vscode/tasks.json`.
- Keine Persistenz (In-Memory-Store).
- Kein CSRF-Schutz auf den Formularen (nicht gefordert, für Produktivbetrieb nachzuholen).