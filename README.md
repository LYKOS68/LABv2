# LABv2

RAP-LAB — Go-Fullstack-Homepage mit Registrierung/Login/Logout/Account-Löschen.

## Start

```
go run ./cmd/lab-studio
```

Läuft auf `http://127.0.0.1:8090`.

## Geprüfter Stand

- Homepage (Text-Platzhalter statt Bild-Assets), Registrierung, Login, Logout, Account-Löschung — funktional getestet (siehe Check-Verlauf in der Session).
- Passwörter via Argon2id gehasht (PHC-Format, OWASP-Mindestparameter m=19456 KiB, t=2, p=1), nie im Klartext gespeichert. Korrupte/fremde Hash-Strings führen zu Fehler statt Absturz (siehe `cmd/lab-studio/main_test.go`).
- Nutzerdaten in SQLite (`data/raplab.sqlite`) — **persistent**, bewiesen durch Neustart-Test (siehe `docs/02-sicherheitserweiterungen-plan.md` Schritt 1). Sessions bleiben bewusst in-memory.

## Tests

```
go test ./...
```

Deckt Argon2id-Hash/Verify-Rundlauf und Korrupt-Hash-Handling ab.

- CSRF-Schutz auf allen Mutationen (Register/Login/Logout/Delete-Account): Synchronizer-Token im versteckten Formularfeld gegen HttpOnly-Cookie, 30 Minuten gültig, zeitkonstanter Vergleich. Fehlendes/falsches Token → 403, Mutation läuft nicht. `/logout` ist jetzt ein POST-Formular statt eines GET-Links (war sonst selbst eine CSRF-Lücke).

## Offene LÜCKEN

- Siehe `docs/01-spezifikation.md` §4: exakte Hex-Farbwerte, Font-Namen, freigestellte Logo-/Signature-Assets fehlen noch (aktuell Text-Platzhalter statt Bild).
- Kein Rate-Limiting — siehe `docs/02-sicherheitserweiterungen-plan.md` Schritt 4 (geplant, noch nicht umgesetzt).

## Verzeichnisstruktur

Siehe `docs/CONTRACT.md` für den vollständigen Pfadvertrag.