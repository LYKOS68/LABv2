# 07 — Plan: E-Mail-Bestätigung bei der Registrierung

**Auslöser:** Roman-Entscheidung (Session 2026-10-18): „man registriert sich einfach" — keine Bestätigung der E-Mail-Adresse. Freigegeben: kompletter Bestätigungsablauf; solange kein Mail-Dienst eingerichtet ist, wird der Link ins Server-Konsolenfenster geschrieben, später werden nur Zugangsdaten eingetragen.

| Punkt | Inhalt |
|---|---|
| Ziel | Ein neues Konto kann sich erst anmelden, nachdem der Bestätigungslink aus der E-Mail geöffnet wurde |
| Mindestkriterien | Registrierung legt Konto als *unbestätigt* an und verschickt einen Link `/verify?token=…`; Login eines unbestätigten Kontos wird abgelehnt (auch mit korrektem Passwort); Öffnen des Links setzt das Konto auf *bestätigt*; Link ist einmalig und läuft nach 24 h ab; Link kann erneut angefordert werden |
| Qualitätskriterien | Token kryptographisch zufällig (32 Byte), in der DB nur als SHA-256-Hash gespeichert (DB-Leck verrät keine gültigen Links); Tokens überleben Neustart (SQLite, nicht in-memory); „Link erneut senden" antwortet immer gleich, egal ob die E-Mail existiert (keine Konto-Ermittlung); erneutes Senden ist CSRF-geschützt und rate-limitiert |
| Erlaubte Mittel | stdlib (`net/smtp`, `crypto/sha256`) — **keine neue Dependency** |
| Mail-Versand | Umgebungsvariablen `SMTP_HOST`, `SMTP_PORT` (Standard 587), `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`. Ist `SMTP_HOST` leer → Link wird per `log` in die Konsole geschrieben. Link-Basis: `BASE_URL` (Standard `http://127.0.0.1:8090`) |
| Bestandskonten | Konten, die vor dieser Änderung existierten, werden bei der Migration einmalig als bestätigt markiert (sonst wären sie ausgesperrt) |
| Prüfverfahren | Unit-Tests: (1) unbestätigtes Konto → Login abgelehnt; (2) ungültiger, leerer, abgelaufener und bereits benutzter Token → abgelehnt, Konto bleibt unbestätigt (Always-Pass-Falle); (3) gültiger Token → Login klappt. Dazu Live-Smoke-Test gegen frisch gebautes Binary |
| Fehlerbedingungen | Login vor Bestätigung möglich; Token mehrfach verwendbar; Klartext-Token in der DB; Bestandskonto ausgesperrt |
| Verbote | kein externer Mail-Dienst fest eingebaut, keine Zugangsdaten im Repo |

**Status:** umgesetzt. `go vet` OK, `go test ./...` 11/11 PASS (neu: unbestätigt → Login abgelehnt; leerer/falscher/abgelaufener/Hash-statt-Klartext-Token → abgelehnt; Token genau einmal gültig und nur als Hash gespeichert; erneutes Senden entwertet alten Link; Migration lässt Bestandskonten rein; Ende-zu-Ende über HTTP mit Link aus dem Konsolen-Log). Gegenprobe: Mit auskommentierter Verified-Prüfung schlagen 2 Tests fehl — kein Always-Pass. Live-Test gegen Port 8090 noch offen (Romans laufender Server belegt den Port).

BAUSTEIN: E-Mail-Bestätigung (`mail.go` + Tabelle `email_verifications`) | AUSLÖSER: Roman-Entscheidung außerhalb §6.3 | BELEG: `handlers.go` registerHandler leitet nach `CreateUser` direkt auf `/login` weiter, keine Prüfung der Adresse | INVARIANTE: `session_erzeugt ⟹ konto.verified == 1`
