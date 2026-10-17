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
- Nutzerdaten in SQLite (`data/raplab.sqlite`) — **persistent**, bewiesen durch Neustart-Test (siehe `docs/02-sicherheitserweiterungen-plan.md` Schritt 1). Sessions bleiben bewusst in-memory.

## Offene LÜCKEN

- Siehe `docs/01-spezifikation.md` §4: exakte Hex-Farbwerte, Font-Namen, freigestellte Logo-/Signature-Assets fehlen noch (aktuell Text-Platzhalter statt Bild).
- Keine Tests (`go test ./...` → no test files).
- Kein CSRF-Schutz auf den Formularen — siehe `docs/02-sicherheitserweiterungen-plan.md` Schritt 3 (geplant, noch nicht umgesetzt).
- Kein Rate-Limiting — siehe `docs/02-sicherheitserweiterungen-plan.md` Schritt 4 (geplant, noch nicht umgesetzt).
- Passwort-Hashing noch `bcrypt`, Umstellung auf Argon2id in Schritt 2 des Sicherheitsplans geplant.

## Verzeichnisstruktur

Siehe `docs/CONTRACT.md` für den vollständigen Pfadvertrag.