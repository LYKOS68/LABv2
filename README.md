# RAP-LAB

**Musik passiert nicht nur im Studio.** RAP-LAB holt an einen Ort, was sonst über Handy-Notizen, Instagram, Discord und 40 Chats verteilt ist — von der ersten Line bis zur Zusammenarbeit. Zweck und Ausbauplan: `hub/rules/01-spezifikation.md`, Abschnitt 13.

Eine Go-Webanwendung, ein Programm, eine SQLite-Datenbank. Kein JavaScript-Framework, kein Build-Schritt.

## Starten

```
go run ./cmd/lab-studio
```

Dann im Browser: **http://127.0.0.1:8090**. Alternativ in VS Code: Terminal → Run Task → „RAP-LAB: Server starten". Beenden mit `Strg+C`.

> Nach jeder Code-Änderung den Server neu starten — `go run` übersetzt den Code nur beim Start.

## Was die Seite kann (Stand: Tag `stand-2026-10-18`)

| Bereich | Adresse | Stand |
|---|---|---|
| **Startseite** | `/` | im Markendesign: Neon-Krone, Gold-Schriftzug, Leder-Hintergrund — Handy und PC. Roman: „noch nicht fertig" |
| **Konto** | `/register`, `/login`, `/forgot-password` | Registrieren mit E-Mail-Bestätigung, Anmelden, Passwort vergessen, Abmelden, Konto löschen |
| **LAB-STUDIO** | `/studio` (nur angemeldet) | verschiebbare, zoombare Arbeitsfläche; Notizen (speichern automatisch); SONGS-Fenster mit „Neuer Song"; Lyrics-Editor (Speichern / Strg+S). Jedes Konto startet leer und sieht nur seine eigenen Daten |
| **LAB-CONNECT** | `/connect` | **Oberflächen-Prototyp** nach dem Original-Entwurf: Entdecken mit Suche, Kategorien und Filter, Merkliste, Anzeige-Seite, Inserieren-Fenster. Arbeitet mit Beispielanzeigen, **speichert noch nichts** |

### Sicherheit

- Passwörter mit Argon2id (OWASP-Mindestparameter), nie im Klartext
- CSRF-Schutz auf allen Schreibzugriffen (Formulare und Studio-Schnittstelle)
- Rate-Limiting: max. 5 Versuche pro 10 Minuten und IP bei Registrieren, Anmelden, Links anfordern
- E-Mail- und Passwort-Links: zufällig, einmalig, befristet (24 h / 1 h), in der Datenbank nur als Hash
- LAB-STUDIO: jede Datenbankabfrage an das angemeldete Konto gebunden — per Test geprüft, dass ein zweites Konto fremde Songs weder sieht noch ändert

### E-Mails

Ohne Mail-Dienst stehen Bestätigungs- und Passwort-Links im **Server-Konsolenfenster** (Zeile beginnt mit `[mail]`). Für echten Versand einmalig vor dem Start setzen: `SMTP_HOST`, `SMTP_PORT` (Standard 587), `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`, `BASE_URL`.

## Marke

Farben, Schrift, Logos und ihre Regeln: `hub/rules/01-spezifikation.md`, Abschnitt 12.

- Farben: Schwarz, Weiß, Neongrün `#ACE182`, Gold `#C2B198`, Braun `#543824` — einmal in `cmd/lab-studio/web/static/brand.css`
- Schrift: Manrope (lokal, SIL Open Font License)
- Originalbilder unverändert in `hub/brand/originale/`, Design-Vorlagen in `hub/brand/referenzen/`

## Prüfen

```
go test ./...
```

37 Tests (Konto, E-Mail- und Passwort-Links, Studio inkl. Konto-Trennung, Connect, Startseite, Schrift). Bei sicherheitskritischen Prüfungen wurde jeweils gegengeprobt, dass der Test rot wird, wenn die Prüfung im Code fehlt. GitHub führt bei jedem Push `gofmt`, `go vet` und alle Tests automatisch aus (`.github/workflows/test.yml`).

## Wo was liegt

```
cmd/lab-studio/        das Programm (ein Go-Paket, nach Themen in Dateien getrennt)
  web/                 Seiten, CSS, JavaScript, Bilder, Schrift — ins Programm eingebaut
data/                  Datenbank, entsteht beim ersten Start, nie auf GitHub
hub/                   alles, was kein Programm ist
  rules/               CONTRACT.md + 01-spezifikation.md (alle Pläne in einem Dokument)
  brand/               Marken-Originale und Design-Vorlagen
  drafts/              Ideen ohne Bauauftrag
  agent-docs/          Einstieg für jede KI, die hier arbeitet
  user/                nur für Roman: Prompts zum Kopieren (YAML-Gruppen, append-only)
```

Verbindlicher Pfadvertrag: `hub/rules/CONTRACT.md`. Alle Pläne, Inventar, Marke und Vision in einem Dokument: `hub/rules/01-spezifikation.md` (Abschnitte nach alter Nummer).

## Offen

- **LAB-CONNECT** speichert noch nichts: Anzeigen in der Datenbank, echte Filter nach Ort, Bilder, Nachrichten, Profile (`hub/rules/01-spezifikation.md`, Abschnitt 11)
- **LAB-STUDIO**: Songs löschen/umbenennen, Released/Drafts, Reimbuch, Relatable Bars, 1-2-1-2-Baukasten (`hub/rules/01-spezifikation.md`, Abschnitt 13 §5)
- **Online gehen**: Server-Anbieter, Domain, Mail-Dienst
- Offene Fragen an Roman: `hub/rules/01-spezifikation.md`, Abschnitt 12 §6 und 13 §7
- `logs/` hat keine Schreibanbindung (Server schreibt nur ins Konsolenfenster)
