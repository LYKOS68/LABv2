# 02 — Plan: Sicherheits- und Persistenz-Erweiterungen

**Auslöser:** Manus-Fragment (`manus-packet`) beschreibt ein durchdachtes Sicherheitsdesign (Argon2id, SQLite, CSRF, Rate-Limiting), liefert aber keinen lauffähigen Code (fehlender Entry Point, fehlendes `go.sum`, defektes `config.go` mit zwei Package-Deklarationen, keine ZIP auffindbar). Entscheidung: Entwurfsmuster selbst neu implementieren, nicht Code übernehmen — auf unserer bestehenden, getesteten Basis (`cmd/lab-studio`).

**Reihenfolge:** Persistenz zuerst (strukturelle Grundlage), danach Argon2id (einfache lokale Änderung, sinnvoll vor echten gespeicherten Hashes), dann CSRF, dann Rate-Limiting. Jeder Schritt einzeln gebaut, getestet, eingefroren — erst dann der nächste (§2d ARBEITSWEISE.md).

---

## Schritt 1 — SQLite-Persistenz

| Punkt | Inhalt |
|---|---|
| Ziel | User-Daten überleben einen Serverneustart; In-Memory-Store wird durch SQLite ersetzt |
| Mindestkriterien | `data/raplab.sqlite` wird beim ersten Start erzeugt; Email eindeutig (case-insensitive); Register/Login/Logout/Delete funktional identisch zu vorher (gleiche HTTP-Codes) |
| Qualitätskriterien | WAL-Modus, `busy_timeout` gesetzt, parametrisierte Queries (kein SQL-Injection-Risiko) |
| Erlaubte Mittel | `modernc.org/sqlite` (CGO-frei) als einzige neue Dependency. Kein ORM |
| Erwartetes Ergebnis | Account registrieren → Prozess beenden → neu starten → Login funktioniert weiterhin |
| Prüfverfahren | bisherige curl-Sequenz + expliziter Neustart-Test |
| Fehlerbedingungen | Daten weg nach Neustart; SQL-Injection mit `'`-Eingabe möglich; doppelte Email trotz Constraint |
| Verbote | kein Postgres/MySQL, keine Session-Persistenz (Sessions bleiben in-memory, laufen ohnehin ab) |

BAUSTEIN: SQLite-Store (`modernc.org/sqlite`) | AUSLÖSER: c | BELEG: Invariante „Nutzerdaten überleben Neustart" nur durch Persistenzschicht hinter Store-Interface erzwingbar | INVARIANTE: `user nach Neustart lesbar == user vor Neustart geschrieben`

---

## Schritt 2 — Argon2id statt bcrypt

| Punkt | Inhalt |
|---|---|
| Ziel | Passwort-Hashing auf Argon2id (OWASP-Erstempfehlung) umstellen |
| Mindestkriterien | Neue Registrierungen erzeugen Argon2id-PHC-Hash (`$argon2id$...`); Parameter mind. m=19456 KiB, t=2, p=1 (OWASP-Minimum, siehe `QUELLEN.md` im Manus-Fragment — echte Quelle, kein erfundener Wert); zeitkonstanter Vergleich bei Verifikation |
| Qualitätskriterien | Korruptes/fremdes Hash-Format führt zu Fehler, nicht zu Panic |
| Erlaubte Mittel | `golang.org/x/crypto/argon2` — **keine neue Dependency**, `x/crypto` ist bereits vorhanden (aktuell für bcrypt) |
| Erwartetes Ergebnis | Register → DB-Eintrag beginnt mit `$argon2id$` |
| Prüfverfahren | Register/Login-Zyklus wie bisher; zusätzlich: korrupten Hash-String in DB schreiben, Login versuchen → Fehler statt Absturz |
| Fehlerbedingungen | Panic bei fremdem Hash-Format; falsches Passwort trotzdem akzeptiert |
| Verbote | keine Migration bestehender bcrypt-Hashes — DB ist zu diesem Zeitpunkt noch ohne echte Nutzer, Breaking Change ist vertretbar und wird hier dokumentiert, nicht übersehen |

BAUSTEIN: Argon2id-Hashing | AUSLÖSER: c | BELEG: OWASP Password Storage Cheat Sheet (zitiert in Manus-`QUELLEN.md`, verifizierbare externe Quelle) | INVARIANTE: `gespeichertes_passwort == argon2id.Hash(klartext)`, nie Klartext, nie SHA-only

---

## Schritt 3 — CSRF-Schutz

| Punkt | Inhalt |
|---|---|
| Ziel | Mutationen (Register/Login/Logout/Delete) gegen Cross-Site-Request-Forgery absichern |
| Mindestkriterien | `GET /csrf` liefert Token + setzt HttpOnly-Cookie; alle POST-Routen prüfen Token (Hidden-Field im HTML-Formular, da server-gerendert, kein JS-Fetch) gegen Cookie-Wert; fehlender/falscher Token → 403, **Mutation wird nicht ausgeführt** |
| Qualitätskriterien | Token zeitlich begrenzt (30 Min), zeitkonstanter Vergleich (`crypto/subtle`) |
| Erlaubte Mittel | stdlib (`crypto/rand`, `crypto/subtle`, `net/http`) — keine externe CSRF-Library |
| Erwartetes Ergebnis | Formulare enthalten verstecktes CSRF-Feld |
| Prüfverfahren | **Check 1 (Always-Pass-Test, roman_code_review_instructions.md):** Request ohne Token UND mit falschem Token → beide 403. Request mit korrektem Token → geht durch wie bisher |
| Fehlerbedingungen | Gate lässt Request ohne/mit falschem Token durch (Always-Pass-Bug) |
| Verbote | kein JS-Framework |

BAUSTEIN: CSRF-Synchronizer-Token | AUSLÖSER: c | BELEG: dokumentierte LÜCKE in README.md „Kein CSRF-Schutz" | INVARIANTE: `mutation_ausgeführt ⟹ token_gültig ∧ token == cookie_wert`

---

## Schritt 4 — Rate-Limiting

| Punkt | Inhalt |
|---|---|
| Ziel | Login/Register gegen Brute-Force-/Spam-Versuche begrenzen |
| Mindestkriterien | max. 5 Versuche pro 10 Minuten pro IP (`RemoteAddr`, kein `X-Forwarded-For` — kein Proxy davor); Überschreitung → 429 + `Retry-After` |
| Qualitätskriterien | Limiter-Map wird bereinigt (keine unbegrenzt wachsende Map) |
| Erlaubte Mittel | stdlib (`sync`, `time`) — kein Redis |
| Erwartetes Ergebnis | 6. Fehlversuch von derselben IP innerhalb 10 Min → 429 |
| Prüfverfahren | 6x falsches Passwort hintereinander posten → letzter Request 429 |
| Fehlerbedingungen | Map wächst unbegrenzt (Memory-Leak); Limit umgehbar durch simplen Retry |
| Verbote | keine IP-Herleitung aus Headern |

BAUSTEIN: In-Memory-Rate-Limiter | AUSLÖSER: b | BELEG: unbegrenzte Login-Versuche = ungemessenes, aber benanntes Risiko (Brute-Force) ohne Obergrenze | INVARIANTE: `versuche_pro_ip_und_fenster ≤ 5`

---

**Status:** Alle vier Schritte umgesetzt und einzeln eingefroren (Commits `0c27703`, `1421d67`, `f3b6d94`, Rate-Limiting-Commit s. `git log`). Jeder Schritt wurde gegen sein eigenes Prüfverfahren verifiziert, bevor der nächste begann.
