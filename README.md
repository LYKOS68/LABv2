# RAP-LAB

**Musik passiert nicht nur im Studio.** RAP-LAB holt an einen Ort, was sonst über Handy-Notizen, Instagram, Discord und 40 Chats verteilt ist — von der ersten Line bis zur Zusammenarbeit.

Eine Go-Webanwendung, ein Programm, eine SQLite-Datenbank. Kein JavaScript-Framework, kein Build-Schritt.

## Starten

```
go run ./cmd/lab-studio
```

Dann im Browser: **http://127.0.0.1:8090**. Alternativ in VS Code: Terminal → Run Task → „RAP-LAB: Server starten". Beenden mit `Strg+C`.

> Nach jeder Code-Änderung den Server neu starten — `go run` übersetzt den Code nur beim Start.

## Was die Seite kann (Stand: Commit `b1a3c75`, 2026-10-18)

| Bereich | Adresse | Stand |
|---|---|---|
| **Startseite** | `/` | im Markendesign: Krone, Schriftzug, Leder-Hintergrund — Handy und PC. Roman: „noch nicht fertig" |
| **Konto** | `/register`, `/login`, `/forgot-password` | Registrieren mit E-Mail-Bestätigung, Anmelden, Passwort vergessen, Abmelden, Konto löschen |
| **LAB-STUDIO** | `/studio` (nur angemeldet) | verschiebbare, zoombare Arbeitsfläche; Notizen (speichern automatisch); SONGS-Fenster mit „Neuer Song"; Lyrics-Editor (Speichern / Strg+S). Jedes Konto startet leer und sieht nur seine eigenen Daten |
| **LAB-CONNECT** | `/connect` | **Oberflächen-Prototyp**: Entdecken mit Suche, Kategorien und Filter, Merkliste, Anzeige-Seite, Inserieren-Fenster. Arbeitet mit Beispielanzeigen, **speichert nichts** |

### Sicherheit

- Passwörter mit Argon2id, nie im Klartext
- CSRF-Schutz auf allen Schreibzugriffen (Formulare und Studio-Schnittstelle)
- Rate-Limiting: max. 5 Versuche pro 10 Minuten und IP bei Registrieren, Anmelden, Bestätigungs-Link erneut anfordern, Passwort vergessen und Passwort neu setzen
- E-Mail- und Passwort-Links: zufällig, einmalig, befristet (24 h / 1 h), in der Datenbank nur als Hash
- LAB-STUDIO: jede Datenbankabfrage an das angemeldete Konto gebunden — per Test geprüft, dass ein zweites Konto fremde Songs und Felder weder sieht noch ändert

### E-Mails

Ohne Mail-Dienst stehen Bestätigungs- und Passwort-Links im **Server-Konsolenfenster** (Zeile beginnt mit `[mail]`). Für echten Versand einmalig vor dem Start setzen: `SMTP_HOST`, `SMTP_PORT` (Standard 587), `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`, `BASE_URL`.

## Marke

- Farben: Schwarz, Weiß, Neongrün `#ACE182`, Gold `#C2B198`, Braun `#543824` — einmal in `cmd/lab-studio/web/static/brand.css`
- Schrift: Manrope (lokal, SIL Open Font License)
- Originalbilder unverändert in `hub/brand/originale/`, Design-Vorlagen in `hub/brand/referenzen/`

## Prüfen

```
go test ./...
```

37 Tests (Konto, E-Mail- und Passwort-Links, Rate-Limiting, Passwort-Hashing, Studio inkl. Konto-Trennung, Connect, Startseite, Schrift). GitHub führt bei jedem Push `gofmt`, `go vet` und alle Tests automatisch aus (`.github/workflows/test.yml`).

## Wo was liegt

```
cmd/lab-studio/        das Programm (ein Go-Paket, nach Themen in Dateien getrennt)
  web/                 Seiten, CSS, JavaScript, Bilder, Schrift — ins Programm eingebaut
data/                  Datenbank, entsteht beim ersten Start, nie auf GitHub
hub/                   alles, was kein Programm ist — wird neu aufgebaut (siehe unten)
  brand/               Marken-Originale und Design-Vorlagen — nie verändern
  user/                nur für Roman: Vorlage für Prompts zum Kopieren
  chatroom.jsonl       Nachrichten zwischen den Modellen, nur anhängen (Format: chatroom.schema.json)
```

## Nächster Neubau: `hub/`

Am 2026-10-18 wurden die alten Regel- und Plandokumente aus `hub/` entfernt (Commit `b1a3c75`), weil sie veraltet und teils von der KI dazugedichtet waren. Ältere Fassungen liegen in der Git-Historie. Der Neubau von `hub/` ist **deklariert, aber noch nicht umgesetzt**. Bis dahin gelten als Beleg nur der Code, die Git-Historie und Romans eigene Aussagen.

## Offen

- **LAB-CONNECT** speichert noch nichts (keine Server-Anbindung)
- **LAB-STUDIO**: Songs löschen und umbenennen gibt es noch nicht
- **Online gehen**: Server-Anbieter, Domain, Mail-Dienst
- Kommentare in `connect.go`, `brand.css` und `connect.js` verweisen noch auf die entfernten Plandokumente
