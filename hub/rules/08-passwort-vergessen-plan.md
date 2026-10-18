# 08 — Plan: Passwort vergessen

**Auslöser:** Roman-Entscheidung (Session 2026-10-18): letzter fehlender Baustein für „gescheite Registrierung und Anmeldung", bevor LAB-CONNECT beginnt.

| Punkt | Inhalt |
|---|---|
| Ziel | Wer sein Passwort vergessen hat, kann per Link an seine E-Mail-Adresse ein neues setzen |
| Mindestkriterien | Anmeldeseite verlinkt „Passwort vergessen?"; `/forgot-password` nimmt eine E-Mail an und verschickt bei existierendem Konto einen Link `/reset-password?token=…`; dort neues Passwort + Wiederholung (gleiche Regeln wie Registrierung: min. 8 Zeichen, identisch); danach funktioniert nur noch das neue Passwort |
| Qualitätskriterien | Link 1 Stunde gültig, einmalig, nur als SHA-256-Hash gespeichert; neuer Link entwertet den alten; Antwort auf `/forgot-password` immer gleich, egal ob die E-Mail existiert (keine Konto-Ermittlung); beide POSTs CSRF-geschützt und rate-limitiert; nach dem Zurücksetzen werden alle offenen Sitzungen des Kontos beendet (wer das alte Passwort gestohlen hat, fliegt raus); der Token wird erst verbraucht, wenn das neue Passwort gültig ist (Tippfehler kostet nicht den Link); Öffnen des Links (GET) verbraucht ihn nicht (Mail-Scanner, die Links vorab öffnen) |
| Nebeneffekt | Erfolgreiches Zurücksetzen markiert das Konto als bestätigt — der Link beweist ebenfalls, dass die Person das Postfach besitzt |
| Erlaubte Mittel | stdlib, vorhandener `Mailer` — **keine neue Dependency** |
| Prüfverfahren | Unit-Tests: ungültiger/leerer/abgelaufener/benutzter Token → abgelehnt, altes Passwort gilt weiter (Always-Pass-Falle); gültiger Token → neues Passwort gilt, altes nicht, Sitzungen weg; zu kurzes Passwort → Token bleibt gültig. HTTP-Ende-zu-Ende: vergessen → Link aus Konsolen-Log → neues Passwort → Anmeldung. Gegenprobe durch Ausschalten der Token-Prüfung |
| Fehlerbedingungen | Zurücksetzen ohne gültigen Token möglich; Token mehrfach verwendbar; Antwort verrät, ob ein Konto existiert; alte Sitzungen bleiben aktiv |
| Verbote | keine Sicherheitsfragen, kein Passwort per Mail verschicken |

**Generalisierung (§6.3 Auslöser a):** Das ist die **zweite** reale Instanz von „Einmal-Link per Mail" (erste: `07-email-verifikation-plan.md`). Gemeinsam ist: Token erzeugen, als Hash mit Ablaufzeit speichern, ältere Tokens derselben Adresse löschen, genau einmal einlösen. Deshalb werden diese zwei Schritte jetzt als `createToken`/`consumeToken` in `store.go` geteilt; jede Instanz behält ihre eigene Tabelle und Gültigkeitsdauer.

**Status:** umgesetzt. `go vet` OK, `go test ./...` 18/18 PASS (nach Lesbarkeits-Umbau der Tests: gemeinsamer Test-Browser, eine Prüfung pro Zeile) (neu: leerer/falscher/abgelaufener/Hash-statt-Klartext-Token und Bestätigungs-Token als Reset-Token → abgelehnt, altes Passwort gilt weiter; gültiger Token → neues Passwort gilt, altes nicht, alte Sitzung beendet, Token nur einmal; Reset bestätigt unbestätigtes Konto; HTTP-Ende-zu-Ende inkl. identischer Antwort für unbekannte Adresse und zu kurzem Passwort ohne Linkverlust). Gegenprobe: Ablaufprüfung ausgeschaltet → 2 Tests FAIL.

BAUSTEIN: Einmal-Token (`createToken`/`consumeToken`) | AUSLÖSER: a | BELEG: `store.go` CreateVerification/VerifyEmail (Instanz 1) + Passwort-Reset (Instanz 2) brauchen identische Token-Logik | INVARIANTE: `token eingelöst ⟹ hash(token) existierte ∧ nicht abgelaufen ∧ danach gelöscht`

BAUSTEIN: Passwort-Reset (Tabelle `password_resets`) | AUSLÖSER: Roman-Entscheidung außerhalb §6.3 | BELEG: kein Weg zurück ins Konto ohne altes Passwort | INVARIANTE: `passwort geändert ohne Login ⟹ gültiger Reset-Token eingelöst`
