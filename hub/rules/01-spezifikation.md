# 14 — Verzeichnisinventar (löst `03-verzeichnisinventar.md` ab)

Stand 2026-10-18, gegen den realen Bestand geprüft (`git ls-files`, `wc -l`, `grep -c "^func Test"`). `03` bleibt als Stand vom 2026-10-17 unverändert stehen.

Für alle Zeilen gilt, wenn nicht anders vermerkt: **Eigentümer** Roman (entscheidet über Inhalt und Umfang), **Schreiber** KI (Roman schreibt keinen Code, `04-arbeitsvertraege.md`).

## Programm — `cmd/lab-studio/` (ein Paket `main`)

| Datei | Zeilen | Zweck | Prüfverfahren |
|---|---|---|---|
| `main.go` | 46 | `main()` + `newMux` (alle Routen an einer Stelle) | `go build`, Tests nutzen `newMux` |
| `handlers.go` | 262 | Startseite, Konto: Registrieren, Bestätigen, Anmelden, Passwort vergessen/zurücksetzen, Abmelden, Löschen | `main_test.go` (Ende-zu-Ende über HTTP) |
| `studio.go` | 325 | LAB-STUDIO: Seite + JSON-Schnittstelle, Sitzungs- und CSRF-Prüfung | `studio_test.go` |
| `connect.go` | 11 | LAB-CONNECT-Seite (Oberflächen-Prototyp) | `connect_test.go` |
| `store.go` | 487 | **einzige Datei mit SQL**: Konten, Sitzungen, Einmal-Links, Songs, Studio-Felder/-Ansicht | Tests inkl. Neustart und Konto-Trennung |
| `csrf.go` | 96 | CSRF-Token | Always-Pass-Test (403 ohne/mit falschem Token) |
| `ratelimit.go` | 78 | Begrenzung der Versuche pro IP | `main_test.go` |
| `password.go` | 64 | Argon2id | Rundlauf + Korrupt-Hash-Tests |
| `mail.go` | 66 | Mail-Versand (SMTP aus Umgebungsvariablen, sonst Konsole) | Ende-zu-Ende-Tests lesen Links aus dem Log |
| `validate.go` | 26 | E-Mail-/Passwort-Regeln | Tests mit ungültigen Eingaben |
| `templates.go` | 209 | HTML/CSS von Startseite und Konto-Seiten | Tests + Screenshot-Sichtprüfung |
| `main_test.go` | 476 | 19 Tests: Konto, Links, Startseite, Navigation, Hilfs-„Test-Browser" | `go test ./...` |
| `studio_test.go` | 350 | 14 Tests LAB-STUDIO | `go test ./...` |
| `connect_test.go` | 47 | 4 Tests LAB-CONNECT + Schrift | `go test ./...` |

### `cmd/lab-studio/web/` — per `go:embed` ins Programm eingebaut

| Datei | Zweck |
|---|---|
| `studio.html`, `connect.html` | Seitengerüste LAB-STUDIO, LAB-CONNECT |
| `static/brand.css` | Markenfarben + Manrope für **alle** Seiten |
| `static/studio.css`, `static/studio.js` | LAB-STUDIO Aussehen / Verhalten |
| `static/connect.css`, `static/connect.js` | LAB-CONNECT Aussehen / Verhalten (Beispieldaten, speichert nichts) |
| `static/img/krone.jpg` | Kopf-Icon Startseite ← `hub/brand/originale/krone-logo-second.jpeg` |
| `static/img/schriftzug.png` | Schriftzug Startseite ← `rap-lab-schriftzug.png` |
| `static/img/leder-hintergrund.jpg` | blasser Hintergrund Startseite ← `unterschrift-ledercover2.jpeg` |
| `static/img/wordmark-header.jpg` | Kopf-Logo LAB-CONNECT, Ausschnitt ← `unterschrift-ledercover2.jpeg` |
| `static/fonts/manrope-latin-{400,500,600,700}-normal.woff2`, `manrope-OFL.txt` | Schrift + Lizenz (aus `@fontsource/manrope` 5.3.0, identisch zum Original-Prototyp) |

## Laufzeitdaten (nicht auf GitHub)

| Pfad | Schreiber | Zweck |
|---|---|---|
| `data/raplab.sqlite` (+ `-wal`, `-shm`) | **der laufende Server** | Konten, Links, Songs, Studio-Felder. Entsteht beim ersten Start |
| `logs/` | — | **LÜCKE**: keine Schreibanbindung, Server schreibt nur ins Konsolenfenster |

## Werkzeuge und Metadaten

| Datei | Schreiber | Zweck |
|---|---|---|
| `go.mod`, `go.sum` | `go`-Befehl, nie von Hand | Abhängigkeiten (einzige direkte: `modernc.org/sqlite`, `golang.org/x/crypto`) |
| `.github/workflows/test.yml` | KI | GitHub prüft bei jedem Push: `gofmt`, `go vet`, `go test` |
| `.vscode/tasks.json` | KI | VS-Code-Aufgabe „RAP-LAB: Server starten" |
| `.claude/settings.json` | Claude Code | Einstellungen des KI-Werkzeugs |
| `.gitignore` | KI | hält `data/`, `logs/`, `*.exe` aus Git heraus |
| `README.md` | KI | Einstieg: was die Seite kann, Start, Stand, Offenes |

## `hub/` — alles, was kein Programm ist

| Pfad | Zweck |
|---|---|
| `hub/README.md` | Einstieg in die zweite Wurzel |
| `hub/agent-docs/00-orientierung.md` | Betriebsregeln + aktueller Stand für jede KI |
| `hub/chatroom.jsonl`, `chatroom.schema.json` | Nachrichten zwischen Agents und Roman (append-only, derzeit leer) |
| `hub/drafts/06-funktions-ausblick-draft.md` | Kleinanzeigen-Mechaniken, **kein Bauauftrag** |
| `hub/brand/originale/` | 11 Marken-Originale in voller Auflösung — **nie verändern** |
| `hub/brand/referenzen/` | 7 Design-Vorlagen (Startseite, LAB-CONNECT-Screenshots, Original-Prototyp-Vorschauen) |

### `hub/rules/` — nummeriert, wird nicht überschrieben, nur abgelöst

| Nr. | Inhalt | Status |
|---|---|---|
| `01` | Startseiten-Spezifikation | LÜCKEN 1–4 abgelöst durch `12` |
| `02` | Sicherheit: SQLite, Argon2id, CSRF, Rate-Limiting | umgesetzt |
| `03` | Verzeichnisinventar 2026-10-17 | abgelöst durch `14` |
| `04` | Arbeitsverträge: wer bedient was | gültig |
| `05` | Zerlegung von `main.go` | umgesetzt |
| `07` | E-Mail-Bestätigung | umgesetzt |
| `08` | Passwort vergessen | umgesetzt |
| `09` | LAB-STUDIO MVP | umgesetzt |
| `10` | Startseite im Markendesign | umgesetzt, Roman: „noch nicht fertig" |
| `11` | LAB-CONNECT Etappe 1 (Oberfläche, an Original angeglichen) | umgesetzt als Oberfläche |
| `12` | Brand-Spec | Entwurf, offene Punkte §6 |
| `13` | Zweck und Vision | Entwurf, offene Fragen §7 |
| `14` | dieses Inventar | gültig |
| `CONTRACT.md` | Pfadvertrag (zwei Wurzeln, Dateien je Achse) | gültig |

(`06` liegt bewusst in `hub/drafts/`.)

## Einstiegspunkt

Einziger: `cmd/lab-studio/main.go`, `func main()`. Start: `go run ./cmd/lab-studio` → `http://127.0.0.1:8090`.

# 13 — Zweck und Vision: Was RAP-LAB ist und werden soll

**Status:** ENTWURF zur Freigabe durch Roman (2026-10-18). Zusammengeführt aus Romans eigenen Texten — nichts davon ist neu erfunden. Quellen:
- `R:\dw\RAP-LAB-windows-version\RAP-LAB-windows-version\__source-docs\workspace-context\idee-draft.md` (Ideenpapier, inkl. „FHIRE"-Vergleich und 1-2-1-2-Methode)
- `readme.md` desselben Prototyps (LAB-STUDIO-Entscheidungen E-1 bis E-7)
- Startseite (`01-spezifikation.md`) und LAB-CONNECT-Entwurf (`11-lab-connect-mvp-plan.md`)

---

## 1. Zweck — warum es RAP-LAB gibt

> **Musik passiert nicht nur im Studio.**

**Das Problem:** Eine Idee wird zum Song. Der Song braucht einen Beat, die Vocals einen Mix, das Release ein Cover, das Projekt ein Video. Gesucht wird über WhatsApp, Instagram, Google, Discord und irgendwelche Gruppen: 17 Tabs, 40 Chats. Und die eigenen Lines liegen verteilt in Handy-Notizen, alten Notizblöcken und Dateien auf dem PC — halb fertig, ohne Zuordnung.

**Es fehlen nicht die Leute. Es ist alles verteilt.**

**Der Zweck:** RAP-LAB holt das Verteilte an einen Ort — von der ersten Line bis zur Zusammenarbeit:

```
Line → Song → Projekt → Producer ↔ Artist → Mix → Cover → Video → Release
```

## 2. Für wen

Rapper, Artists, Songwriter, Producer, Beatmaker, Mixer, Videografen, Designer — alle, die an einem Musikprojekt arbeiten oder dafür gebraucht werden.

> Kein Label nötig, keine Millionen Follower. Es reicht, etwas zu können, das jemand anderes für sein Projekt braucht.

## 3. Die drei Bereiche

| Bereich | Zweck | Kern-Frage |
|---|---|---|
| **LAB-STUDIO** | Der Ort, an dem Lyrics geschrieben, verfeinert, entwickelt und archiviert werden. **Kein** Tonstudio: keine DAW, keine Beats, kein Mixing, kein Recording | „Wie wird aus meinem Chaos an Lines ein Song?" |
| **LAB-CONNECT** | Menschen und Leistungen der Musikbranche verbinden — Angebote und Gesuche wie Kleinanzeigen, filterbar nach Kategorie, Standort und Region | „Was kann ich selbst, und wen brauche ich für den Rest?" |
| **NEWS** | Relevante Branchen-Informationen statt Feed für alle: Plattformregeln, Zahlen, Tools, Interviews, Releases | „Was bewegt sich gerade, und was heißt das für meine Arbeit?" |

Dazu langfristig **Community** und **Kreatives Profil** (Abschnitt 5).

## 4. Was RAP-LAB besonders macht

Nicht eine einzelne Funktion, sondern die **Verbindung** aus Songwriting-Werkstatt, Kleinanzeigen für die Musikbranche, kreativem Portfolio, lokaler Vernetzung und Karriereförderung. Die bestehenden Werkzeuge decken jeweils nur ein Stück ab (Notiz-App, Instagram, Discord, Kleinanzeigen).

**Die Songwriting-Methode im Studio (Romans 1-2-1-2-Taktik):** Schreiben ist Puzzeln. 4 Bars = Line 1, Line 2, Line 1, Line 2. 16 Bars = 4 × 4 Bars. Man nimmt einen **Reimblock** (z. B. 1 Silbe: Hass, krass, was, lass), legt **Relatable Bars** daneben (gesammelte starke Zeilen, z. B. „Schwarze Augen, Seele schimmelt von dem Hass") und baut daraus Paar für Paar. Fertige gute Zeilen wandern zurück in die Relatable Bars. Genau dafür ist die verschiebbare Arbeitsfläche da: Lyrics, Reime und Zeilen nebeneinander, frei anordenbar.

## 5. Wohin es geht — Ausbaustufen

✅ = gebaut (Tag `stand-2026-10-18`) · ◐ = Oberfläche steht, Funktion fehlt · ○ = geplant

**Fundament**
- ✅ Konto: Registrieren, E-Mail-Bestätigung, Anmelden, Passwort vergessen, Konto löschen; Sicherheit (Argon2id, CSRF, Rate-Limiting)
- ✅ Startseite im Markendesign
- ○ Online gehen: Server-Anbieter, eigene Domain, Mail-Dienst

**LAB-STUDIO**
- ✅ Arbeitsfläche (verschieben, zoomen), Notizen, SONGS-Fenster, Lyrics-Editor — pro Konto, privat
- ○ Songs löschen/umbenennen; Trennung **Released / Drafts**; Cover-Bild je Song
- ○ **REIMBUCH**: Reimblöcke nach Reimart / Silbenzahl
- ○ **RELATABLE BARS**: Ablage einzelner Lines und Phrasen
- ○ **1-2-1-2-Baukasten**: 4-Bar-Blöcke aus Lines zusammenpuzzeln, Silben zählen
- ○ Fenster aneinanderlegen und schließen

**LAB-CONNECT**
- ◐ Entdecken, Anzeige, Inserieren (Handy + PC) mit Beispieldaten
- ○ Anzeigen speichern, bearbeiten, löschen; echte Filter nach Standort/Region
- ○ Bilder hochladen
- ○ Favoriten; Nachrichten zwischen Nutzern
- ○ Video-Beiträge und Feedback an Anzeigen („Zeig, was du kannst")

**Community**
- ○ Gruppen suchen, beitreten, erstellen (ähnlich Skool)
- ○ Meetups, Gigs, lokale Kontakte

**Kreatives Profil**
- ○ „Wer bin ich? Was mache ich? Was kann ich zeigen?" — Profil als Übersicht/Dashboard mit Arbeiten, Anzeigen, Kontakten

**NEWS**
- ○ Branchen-Infos, kuratiert statt Feed

## 6. Abgelöste frühere Annahmen

| Früher | Heute | Warum |
|---|---|---|
| LAB-STUDIO als lokale exe, Website getrennt (Prototyp-README „Produktmodell") | alles **eine Webanwendung**, LAB-STUDIO unter `/studio` | Prototyp-Entscheidung E-2 + Romans Entscheidung 2026-10-18 („in unsere Runtime einbauen") |
| Romans 23 Songs als Startbestand | jedes Konto startet **leer** | Romans Entscheidung 2026-10-18: die eigenen Texte waren nur Entwicklungsmaterial |
| „LAB-STUDIO zuerst, alles andere danach" (E-1) | Studio und Connect stehen nebeneinander, Ausbau abwechselnd | Stand 2026-10-18 |

## 7. Offene Fragen an Roman

1. **Anfänger oder Profis?** Die Prototyp-README sagt „richtet sich an Leute, die Musik ernsthaft machen. Ein Hobby-Werkzeug ist RAP-LAB nicht." Das Ideenpapier nennt ausdrücklich auch „Anfänger und Neulinge". Vorschlag zum Vereinen: *offen für Anfänger, gebaut für Ernsthaftigkeit* — wer anfängt, ist willkommen, aber die Werkzeuge sind für echtes Arbeiten gedacht.
2. **Geld:** Soll RAP-LAB irgendwann etwas kosten (z. B. hervorgehobene Anzeigen, Pro-Konto) oder dauerhaft kostenlos sein? Das beeinflusst, wie Konto und Anzeigen gebaut werden.
3. **Reihenfolge:** Was kommt als Nächstes — LAB-STUDIO-Werkzeuge (Reimbuch, Relatable Bars), LAB-CONNECT mit echten Anzeigen, oder zuerst online gehen?
4. **NEWS:** Wer schreibt die Inhalte — du selbst, Gastautoren, oder gesammelt aus Quellen?

# 12 — Marken-Spezifikation (Brand-Spec) RAP-LAB

**Status:** ENTWURF zur Freigabe durch Roman (2026-10-18). Abschnitte mit ✅ sind bereits durch Romans frühere Entscheidungen festgelegt, ❓ sind offen.

**Auslöser:** Roman-Entscheidung (Session 2026-10-18): „die Original-Bilder, Logos etc. und Brand-Spec festsetzen". Löst die LÜCKEN 1–4 aus `01-spezifikation.md` §4 ab.

---

## 1. Ablage der Originale

| Pfad | Inhalt | Regel |
|---|---|---|
| `hub/brand/originale/` | alle Marken-Bilder in voller Auflösung (Kopie aus `C:\Users\RPC\dev\rap-app\my-space\__source-docs\bilder`, 2026-10-18) | **nie verändern** — einzige Quelle der Wahrheit. Bearbeitete Fassungen entstehen nur als Kopie |
| `hub/brand/referenzen/` | Design-Vorlagen (Screenshots Startseite, LAB-CONNECT) | Vorlagen zum Nachbauen, keine Marken-Bilder |
| `cmd/lab-studio/web/static/img/` | verkleinerte Web-Kopien, die die Seite tatsächlich ausliefert | entstehen nur aus `originale/`; Herkunft steht in Tabelle 2 |

## 2. Logos und Bildmarken

| Datei (`originale/`) | Was | Einsatz | Web-Kopie |
|---|---|---|---|
| `krone-logo-second.jpeg` | Neon-Krone mit Skyline, Schaltkreis-Linie, Wellenform | ✅ **Icon** oben links auf der Startseite | `img/krone.jpg` (192 px, zugeschnitten) |
| `rap-lab-schriftzug.png` | Gold-Signatur mit grünem Spray, **transparent** | ✅ **Haupt-Schriftzug**: Startseite Mitte, LAB-CONNECT-Kopf | `img/schriftzug.png` (1000 px) |
| `unterschrift-ledercover2.jpeg` | Gold-Signatur, grüner Neon-Glow, auf Leder | ✅ **Hintergrund** der Startseite, nur blass (≥ 80 % abgedunkelt); ✅ **Kopf-Logo** in LAB-CONNECT als Ausschnitt (wie im Original-Prototyp) | `img/leder-hintergrund.jpg` (1280 px), `img/wordmark-header.jpg` (640 × 201, Ausschnitt) |
| `krone-logo.jpg` | schlichte Neon-Krone, schwarz, hochkant | ❓ Vorschlag: **Favicon / App-Icon** (klein gut lesbar) | — |
| `krone-logo-schrift-main.jpeg` | Krone + Signatur „RAP-LAB" | ❓ Vorschlag: Social-Media-Profilbild, Teilen-Vorschau | — |
| `krone-logo-third.jpeg` | Krone mit Skyline, Schaltkreisen, Tropfen | ❓ Vorschlag: großes Motiv (Poster, Ladebildschirm) | — |
| `krone-logo-whiteback.jpeg` | Krone auf Weiß | ❓ Vorschlag: helle Untergründe, Druck | — |
| `unterschrift-whiteback.jpeg` | Signatur auf Weiß | ❓ Vorschlag: helle Untergründe, Druck | — |
| `unterschrift-ledercover.jpeg`, `unterschrift-lederbuchcover.jpeg` | Signatur auf Leder / Notizbuch | ❓ Vorschlag: Stimmungsbilder (LAB-STUDIO = Notizbuch-Gefühl) | — |
| `kassete-cover.png` | Kassette mit Krone + Signatur | ❓ Vorschlag: Stimmungsbild, Platzhalter für Song-Cover im Studio | — |

**Regeln:**
- Die Krone MUSS auf dunklem Grund stehen (Original-Hintergrund ist Schwarz, ohne Transparenz). Auf Hell nur `krone-logo-whiteback.jpeg`.
- Der Schriftzug DARF NICHT als Text nachgebaut werden; er ist eine Handschrift und wird immer als Bild eingesetzt.
- Bilder MÜSSEN einen `alt`-Text haben („RAP-LAB").

## 3. Farben ✅

Festgelegt von Roman („genau so"), Quelle: LAB-STUDIO-Prototyp. Technisch einmal in `cmd/lab-studio/web/static/brand.css`.

| Name | Hex | Rolle |
|---|---|---|
| Schwarz | `#000000` | Grundfläche |
| Weiß | `#FFFFFF` | Fließtext, Überschriften |
| Neongrün | `#ACE182` | Aktion, Aktives, Hervorhebung (Knöpfe, gewählte Kategorie) — sparsam |
| Gold | `#C2B198` | Marke: Navigation, Taglines, Labels |
| Braun | `#543824` | Linien, Rahmen, Raster |
| Rot | `#FF6B6B` | nur Fehlermeldungen |

Abgeleitete Flächentöne in LAB-CONNECT (`#14100c`, `#1c1611`, `#2a2017`) sind dunkle Braun-Stufen und gehören zur Palette.

**Regel:** Neue Farben SOLLTEN NICHT eingeführt werden; wenn doch, nur als neue Zeile hier und in `brand.css`.

## 4. Schrift

| Punkt | Stand |
|---|---|
| Schrift | **Manrope** für die ganze Seite — ✅ (Roman, 2026-10-18; Quelle: LAB-CONNECT-Original-Prototyp). Gewichte 400/500/600/700, Teilmenge Latin (inkl. ä ö ü ß €), lokal unter `cmd/lab-studio/web/static/fonts/`, Lizenz SIL OFL 1.1 (`fonts/manrope-OFL.txt`). Ersatz: `system-ui, sans-serif` |
| Großbuchstaben-Zeilen (Nav, Tagline, Kicker) | Großbuchstaben + Buchstabenabstand 1–3 px — ✅ |
| Regel | Schriften MÜSSEN lokal ausgeliefert werden, nie von fremden Servern (Datenschutz) |

## 5. Sprache und Ton

- Du-Ansprache, direkt, auf Augenhöhe: „Denn mal ehrlich:", „Zeig, was du kannst.", „Was fehlt deinem Projekt?"
- Kurze Sätze, Kicker in Großbuchstaben über der Überschrift („DEIN SOUND. DEINE LEUTE.")
- Echte Umlaute (ä, ö, ü, ß) in allem, was Nutzer sehen
- Begriffe immer gleich (Glossar-Regel aus Romans Ideenpapier): **LAB-STUDIO**, **LAB-CONNECT**, **NEWS**, **Song**, **Lyrics**, **Anzeige**, **Inserieren**

## 6. Offene Punkte für Roman

1. Einsatz der ❓-Bilder aus Tabelle 2 so freigeben oder ändern
2. Favicon (Browser-Tab-Symbol): `krone-logo.jpg`?

# 11 — Plan: LAB-CONNECT, Etappe 1

**Status:** freigegeben als **Oberflächen-Prototyp** (Roman, 2026-10-18): „es geht nur darum die Oberfläche bzw. den Prototyp wieder einzubauen, die ganzen Klicks und Funktionen danach." Roman hat die Detailentscheidungen an die KI delegiert.

## Umfang jetzt (Oberfläche)

- Route `/connect`, ohne Anmeldung erreichbar; Nav-Link LAB-CONNECT zeigt dorthin
- Oberes Band „ETAPPE 1 · Designentwurf · keine echten Anzeigen", Schriftzug auf Leder, Profil-Symbol
- **Entdecken** mit Suchfeld, Ort-Zeile, Kategorie-Knöpfen (Für dich, Mix & Master, Beats, Songwriting, Features) — **Kategorien filtern die Beispiel-Kacheln wirklich**, Suche filtert nach Titel
- **Anzeige** (Detail) per Klick auf eine Kachel, Zurück zu Entdecken
- **Inserieren** öffnet ein Feld, das sich mit ✕ schließen lässt: biete/suche, Bilder (Platzhalter „Etappe 2"), Titel, Kategorie, Beschreibung, Preis + Einheit (`PLATZHALTER`), Ort. „Anzeige aufgeben" fügt die Anzeige **nur in der offenen Ansicht** hinzu (Hinweis: wird nicht gespeichert)
- Untere Leiste wie im Entwurf; Favoriten, Nachrichten, Profil zeigen „folgt in Etappe 2"
- Beispielanzeigen sind als BEISPIEL markiert; nichts wird an den Server geschickt

**Bewusste Abweichung von der Regel gegen tote Knöpfe:** Für den Oberflächen-Prototyp ausdrücklich von Roman gewollt; jeder noch funktionslose Knopf sagt das sichtbar („folgt in Etappe 2").

**Umsetzung:** `connect.go`, `web/connect.html`, `web/static/connect.css`, `web/static/connect.js`. `go test ./...` 36/36 PASS (neu: Nav-Link, Seite ohne Anmeldung inkl. Dateien, Inserieren-Feld mit Schließen-Knopf/Preis/`PLATZHALTER`). `node --check connect.js` OK. Sichtprüfung per Screenshot in 412 px: Entdecken, Anzeige, Inserieren und untere Leiste wie im Entwurf. **PC-Layout** (ab 900 px, Roman-Wunsch): Navigation oben rechts statt Leiste unten, Kachel-Raster über die volle Breite, Anzeige zweispaltig (Galerie links, Angaben rechts), Inserieren als Fenster in der Mitte — geprüft per Screenshot in 1440 px. Browser-Prüfung durch Roman offen.

**Angleichung an den wiedergefundenen Original-Prototyp** (Roman, 2026-10-18: „ja so loslegen"). Original: `R:\dw\RAP-LAB-Etappe-1\RAP-LAB` (React/TypeScript/Node — nur Vorlage für Aussehen und Verhalten, kein Code übernommen; Stack-Vorgaben seiner `REQUIREMENTS.md` wie MariaDB/Better Auth gelten für LABv2 nicht). Übernommen: sechs Kategorien (+ Video, Artwork) und drei weitere im Formular; Filter-Feld „Nur Remote" hinter Regler/Ortszeile; **funktionierende Merkliste** (Herz, Favoriten-Tab „Deine Merkliste", leerer Zustand „Hier ist noch Platz."); Anzeige mit „Über das Angebot", „Arbeitsproben", Künstlerprofil-Karte, Hinweis „RAP-LAB vermittelt Kontakte …", fester Nachricht-Knopf auf dem Handy (untere Leiste dort ausgeblendet); Preisart (Festpreis, Ab-Preis, Verhandelbar, Auf Anfrage, Zusammenarbeit) statt `PLATZHALTER`, Betragsfeld nur wo passend; Nachrichten/Profil mit Krone und Hinweis; Kopf-Logo als Leder-Ausschnitt `img/wordmark-header.jpg`; Schrift Manrope (seitenweit). Behalten (Romans frühere Wünsche): PC-Layout, Inserieren als schließbares Fenster, neue Anzeige erscheint in der Ansicht. Farben: festgelegte Markenfarben statt der leicht abweichenden Prototyp-Töne. Prüfung: `go test ./...` 37/37 PASS, Screenshots Handy (5 Zustände) und PC (Anzeige).

Danach (eigener Plan): Speichern in der Datenbank, echte Filter, Favoriten, Nachrichten, Bilder.

---

## Ursprünglicher Entwurf (Fragen inzwischen durch Delegation entschieden)

**Auslöser:** Roman-Entscheidung (Session 2026-10-18). Der ursprüngliche Prototyp ist verloren; Vorlage sind drei Handy-Screenshots eines Designentwurfs (`C:\Users\RPC\dev\rap-app\my-space\__source-docs\bilder\kleinanzeigen-bsp1.jpeg` bis `bsp3.jpeg`). Strukturrahmen: `hub/drafts/06-funktions-ausblick-draft.md` (Kleinanzeigen-Metapher).

## Was die Vorlage zeigt

| Bildschirm | Inhalt |
|---|---|
| **Entdecken** | Kopfzeile „DEIN SOUND. DEINE LEUTE." / „Was fehlt deinem Projekt?"; Suchfeld „Beats, Mixing, Features …" mit Filter-Symbol; Ort „Deutschland & Remote"; runde Kategorie-Knöpfe (Für dich, Mix & Master, Beats, Songwriting, …); Kacheln „Entdecke Möglichkeiten" mit Bild, Herz, Ort, Titel |
| **Anzeige** | Zurück zu Entdecken; Bildergalerie; Ort; Herz (merken); Titel; Preis („50 € pro Song"); Kategorie; grüner Knopf „Nachricht schreiben" |
| **Inserieren** | „DEIN KÖNNEN ZÄHLT." / „Zeig, was du kannst. Oder sag, wen du für dein Projekt suchst."; Auswahl **Ich biete / Ich suche**; Bilder hinzufügen (im Entwurf selbst als „Etappe 2" markiert); Titel; Kategorie; Beschreibung; … |
| überall | untere Leiste: Entdecken, Favoriten, Inserieren, Nachrichten, Profil; dunkelbraune Karten, Gold-Schrift, Neongrün für Aktives |

## Vorschlag Etappe 1 (Mechaniken aus Draft 06)

| Mechanik | in Etappe 1 | Funktion |
|---|---|---|
| CRUD | ✅ | Anzeige erstellen (biete/suche, Titel, Kategorie, Beschreibung, Preis + Einheit oder „Verhandlungssache", Ort oder Remote); eigene Anzeigen unter Profil → „Meine Anzeigen" ansehen |
| Liste + Filter | ✅ | Entdecken: alle Anzeigen, neueste zuerst; Suche im Titel/Text; Filter nach Kategorie, biete/suche, Remote |
| Detailseite | ✅ | Anzeige ansehen |
| Relations-Liste (Favoriten) | ❓ | Herz + Favoriten-Tab |
| Mailbox (Nachrichten) | ❓ | „Nachricht schreiben" + Nachrichten-Tab |
| Bilder-Upload | ❌ Etappe 2 | bis dahin Platzhalter mit Kategorie-Symbol, wie im Entwurf |
| Bearbeiten/Löschen/Pausieren | ❌ später | analog zu LAB-STUDIO (Roman: Löschen/Umbenennen später) |

**Regel gegen tote Knöpfe:** Was nicht in Etappe 1 ist, wird nicht als Knopf angezeigt (kein Herz ohne Favoriten, kein „Nachricht schreiben" ohne Nachrichten).

## Offene Fragen an Roman

1. **Favoriten** und **Nachrichten** schon in Etappe 1 — oder nur Inserieren/Entdecken/Anzeige?
2. **Ohne Anmeldung stöbern?** Wie bei Kleinanzeigen: Entdecken und Anzeige für alle, Inserieren nur angemeldet. Oder alles nur angemeldet?
3. **Kategorien:** vollständige Liste? Sichtbar sind Mix & Master, Beats, Songwriting. Vorschlag zum Ergänzen: Features, Recording/Studio, Video, Cover/Artwork, Promotion, Sonstiges.
4. **Preis-Einheiten:** „pro Song" ist sichtbar. Vorschlag: pro Song, pro Stunde, pro Projekt, Festpreis, Verhandlungssache, kostenlos.

## Kriterien (nach Freigabe zu präzisieren)

| Punkt | Inhalt |
|---|---|
| Prüfverfahren | Go-Tests: Anzeige erstellen nur angemeldet + CSRF; Pflichtfelder fehlen → abgelehnt, nichts gespeichert (Always-Pass-Falle); Filter liefern nur passende Anzeigen; „Meine Anzeigen" zeigt nur eigene; fremder Nutzer kann keine Anzeige im Namen eines anderen anlegen. Danach Sichtprüfung per Screenshot (Handy- und PC-Breite) und durch Roman |
| Verbote | kein JS-Framework, keine neue Go-Dependency, keine Bild-Uploads in Etappe 1 |

# 10 — Plan: Startseite nach Romans Vorlage

**Auslöser:** Roman-Entscheidung (Session 2026-10-18). Vorlage: Handy-Screenshot einer verlorenen Version (`C:\Users\RPC\dev\rap-app\my-space\__source-docs\bilder\homepage-bsp.jpeg`). Schließt die LÜCKE „Logo-/Signature-Assets fehlen" aus `01-spezifikation.md` §4.

**Status:** umgesetzt. `go test ./...` 33/33 PASS (neu: `TestHomeShowsBrandImages`). Sichtprüfung per Edge-Screenshot in 1280 px und in einem 412-px-Rahmen (Handybreite): Aufbau wie Vorlage, keine seitliche Verschiebung. Browser-Prüfung durch Roman offen.

| Punkt | Inhalt |
|---|---|
| Ziel | Startseite sieht aus wie die Vorlage |
| Mindestkriterien | oben links Kronen-Icon (`krone-logo-second.jpeg`), rechts LAB-STUDIO / LAB-CONNECT / Anmelden in Gold; Mitte Schriftzug-Bild (`rap-lab-schriftzug.png`) statt Text-Überschrift; Tagline Gold in Großbuchstaben; vier Punkte Weiß fett; Fließtext-Block über die ganze Breite mit Leder-Schriftzug (`unterschrift-ledercover2.jpeg`) blass im Hintergrund; „RAP-LAB." Gold am Ende |
| Qualitätskriterien | Bilder als verkleinerte Web-Kopien unter `cmd/lab-studio/web/static/img/` (Krone 192 px ≈ 12 KB, Schriftzug 1000 px, Leder 1280 px) — Originale bleiben unverändert; Bilder haben `alt`-Text; Seite funktioniert auch auf dem Handy |
| Prüfverfahren | Test: Startseite bindet die drei Bilder ein, alle drei werden ausgeliefert. Danach Sichtprüfung durch Roman im Browser (Handy + PC) |
| Verbote | keine neue Dependency, keine Schrift-Dateien von außen |

BAUSTEIN: Marken-Bilder (`web/static/img/`) | AUSLÖSER: Roman-Entscheidung außerhalb §6.3 | BELEG: `01-spezifikation.md` §4 LÜCKE Logo/Signature | INVARIANTE: `jedes eingebundene Bild wird mit Status 200 ausgeliefert`


# 09 — Plan: LAB-STUDIO als MVP in die Runtime übernehmen

**Status:** umgesetzt, Browser-Test durch Roman offen. `go vet` OK, `go test ./...` 32/32 PASS (davon 14 Studio: Zugang ohne Anmeldung → Login/401; Schreiben ohne/falsches CSRF-Token → 403, nichts gespeichert; neues Konto leer; Song speichern/laden; Feld-Änderung lässt ID/Typ unverändert; Konto B kann Songs/Felder von A weder lesen noch ändern noch darauf zeigen; ungültige Eingaben → 400, nichts gespeichert; Zoom begrenzt; Daten überleben Neustart; Konto löschen löscht Studio-Daten; Seite lädt ihre Dateien). Gegenprobe: Konto-Bindung (`owner = ?`) aus den 4 ID-Abfragen entfernt → beide Konto-Trennungs-Tests FAIL. `node --check studio.js` OK.

**Auslöser:** Roman-Entscheidung (Session 2026-10-18): Der LAB-STUDIO-Prototyp (`R:\dw\RAP-LAB-windows-version\RAP-LAB-windows-version\main.go`, 685 Zeilen, eigener Server auf Port 3000) wird refaktoriert und als MVP in `cmd/lab-studio` eingebaut. Romans eigene Lyrics waren nur Entwicklungsmaterial — **jeder Nutzer startet mit leerem Studio**, auch Roman.

## Was der Prototyp kann (Ist-Stand, gelesen)

- Endlose Arbeitsfläche: verschieben mit Maus/Finger, zoomen mit Mausrad/Zwei-Finger/Knöpfen (25 %–400 %), Punkt-Raster im Hintergrund
- Linke Leiste: `SONGS`, `+ FELD` (Notiz), `ZOOM +`, `ZOOM −`, `ANSICHT` (zurücksetzen), Zoom-Anzeige
- Felder auf der Fläche, an der Kopfzeile verschiebbar:
  - **Notiz**: Textfeld, speichert automatisch 0,5 s nach dem Tippen
  - **Songs**: Kacheln (Bildplatzhalter + Name) → Klick öffnet den Text im Editor mit `Zurück`, `Speichern` und Strg+S
- Ansicht und Felder überleben Neuladen und Neustart (`data/workspace.json`), Songs liegen als `.txt`-Dateien
- Schreibzugriffe laufen im Browser nacheinander (keine überholten Speichervorgänge); Server ändert eine Kopie und übernimmt nur bei Erfolg
- Farben: Schwarz `#000000`, Weiß `#FFFFFF`, Neongrün `#ACE182`, Gold `#C2B198`, Braun `#543824`

## Was sich beim Einbau ändern muss (und warum)

| Prototyp | MVP in der Runtime | Grund |
|---|---|---|
| ein gemeinsamer Song-Ordner, kein Konto | Songs, Felder und Ansicht **pro Konto** in SQLite | Website mit fremden Nutzern |
| Songs nur aus Seed-Dateien, kein Anlegen | **neu:** `+ SONG` legt einen leeren Song mit Titel an | leeres Studio ohne Seed wäre sonst unbenutzbar |
| Seed aus `__source-docs/lykos-lyrics` | kein Seed | Romans Entscheidung |
| eigener Server, Port 3000, nur Loopback | Route `/studio` im bestehenden Server (8090), nur angemeldet | eine Runtime |
| Schutz nur über `Origin`-Prüfung | Sitzung + CSRF-Token als Header `X-CSRF-Token` bei jedem Schreibzugriff | bestehendes Gesetz `mutation_ausgeführt ⟹ token_gültig` |
| alles in einer `main.go` (Prototyp-Regel E-3) | Aufteilung nach `CONTRACT.md` (s. Dateien) | LABv2-Vertrag gilt, nicht der des Prototyps |
| JavaScript als dichte Einzeiler im Go-String | eigene Datei, lesbar: kurze Funktionen, eine Anweisung pro Zeile | Romans Lesbarkeitsregel |
| Prüfskripte `00_worm` (Python, AST auf `main.go`, Port 3000) | normale Go-Tests in `*_test.go` | Skripte prüfen eine Datei/einen Port, die es so nicht mehr gibt |

## Dateien (neu/geändert)

| Datei | Achse | Inhalt |
|---|---|---|
| `store.go` | Persistenz | Tabellen `songs`, `studio_fields`, `studio_view` + Zugriffe — bleibt einzige Datei mit SQL |
| `studio.go` (neu) | HTTP-Routing Studio | JSON-API unter `/studio/api/…`, Seite `/studio` |
| `web/studio.html`, `web/studio.css`, `web/studio.js` (neu, unter `cmd/lab-studio/`) | Präsentation Studio | per `go:embed` ins Programm eingebaut — erster echter Inhalt für `web/` (CONTRACT Regel 1); `go:embed` erlaubt nur Unterordner des Pakets, daher nicht Top-Level |
| `main.go` | Wiring | Studio-Routen in `newMux` |
| `templates.go` | Präsentation | Nav-Link `LAB-STUDIO` → `/studio` |

## Kriterien

| Punkt | Inhalt |
|---|---|
| Ziel | Angemeldeter Nutzer öffnet `/studio`, legt Songs und Notizen an, schreibt Texte, verschiebt die Fläche — alles ist nach Neuladen und Server-Neustart wieder da |
| Mindestkriterien | alle Prototyp-Funktionen oben + `+ SONG`; neues Konto sieht leere Fläche ohne Songs; nicht angemeldet → `/studio` leitet auf `/login`, API antwortet 401 |
| Qualitätskriterien | jede SQL-Abfrage ist an das Konto gebunden; Text max. 1 MB; Zoom serverseitig auf 0,25–4 begrenzt; unbekannter Feldtyp → 400; Ändern = ganz oder gar nicht (Transaktion) |
| Prüfverfahren | Go-Tests: (1) Konto B kann Song/Feld von Konto A weder lesen noch ändern (404) — **Always-Pass-Falle**: Test muss rot werden, wenn die Konto-Bindung aus der Abfrage entfernt wird; (2) ohne Sitzung 401, ohne/falsches CSRF-Token 403, nichts gespeichert; (3) ungültiger Typ/Größe/Body → 400, nichts gespeichert; (4) Song speichern → neuer Store auf derselben DB-Datei liest denselben Text; (5) neues Konto: leere Liste. Danach Live-Test im Browser durch Roman |
| Fehlerbedingungen | fremde Songs sichtbar; Schreiben ohne Token möglich; Änderung nach Neustart weg; Romans Lyrics tauchen in einem Konto auf |
| Verbote | kein JS-Framework, kein Build-Schritt, keine neue Go-Dependency; nichts aus der Prototyp-Liste „Nicht im ersten Abschnitt" (Reimbuch, Released/Drafts, Relatable-Bars, 1-2-1-2, Multi-User-Canvas) |

## Entscheidungen (Roman, 2026-10-18)

1. **Farben:** Die Prototyp-Farben gelten für die **ganze** Seite. Sie liegen einmal in `web/static/brand.css` und werden von Startseite/Anmeldung und Studio gemeinsam genutzt (zweite reale Verwendung → gemeinsame Datei, §6.3 a).
2. **Löschen/Umbenennen** von Songs und Notizen: **nicht** im MVP, kommt später.
3. **Speichern:** wie im Prototyp — Songtext per Knopf/Strg+S, Notizen automatisch.
4. `+ SONG` sitzt als kleines Eingabefeld „Neuer Song" im SONGS-Fenster über den Kacheln; die linke Leiste bleibt exakt wie im Prototyp.

BAUSTEIN: LAB-STUDIO (Tabellen `songs`/`studio_fields`/`studio_view`, `studio.go`, `web/`) | AUSLÖSER: Roman-Entscheidung außerhalb §6.3 | BELEG: Prototyp `main.go` (685 Z.) läuft getrennt, ohne Konten | INVARIANTE: `song/feld gelesen oder geändert ⟹ song.owner == angemeldeter_nutzer`

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


# 04 — Arbeitsverträge: wer bedient was

Macht jede ausführbare Handlung ohne eigene Programmierkenntnisse bedienbar. Roman schreibt keinen Code — das übernimmt in dieser Zusammenarbeit durchgehend die KI (bestätigt in `roman_code_review_instructions.md`: seine reale Lücke ist „Code selbst schreiben", nicht Konzepte verstehen).

## Wer bearbeitet was

| Bereich | Bearbeitet von |
|---|---|
| Alle `.go`-Dateien (`cmd/lab-studio/*.go`) | **ausschließlich KI** |
| Alle Markdown-Dokumente (`hub/rules/*.md`, `README.md`) | **ausschließlich KI** |
| `.vscode/tasks.json`, `.gitignore` | **ausschließlich KI** |
| `go.mod`/`go.sum` | **ausschließlich die `go`-CLI** (über `go get`/`go mod tidy`, von der KI ausgeführt) — nie von Hand |
| Visuelle Freigabe der Homepage im Browser | **ausschließlich Roman** — nur er kann beurteilen, ob es „richtig aussieht" |
| Entscheidung über Struktur/Umfang/Fertigstellung | **ausschließlich Roman** (§3.1 zusammenarbeit-gesetze.md) |

## Was die KI automatisch ausführt (keine Freigabe nötig)

| Prüfung | Befehl | Wann |
|---|---|---|
| Kompilier-Check | `go build ./cmd/lab-studio` | nach jeder Code-Änderung |
| Statische Analyse | `go vet ./...` | nach jeder Code-Änderung |
| Unit-Tests | `go test ./...` | nach jeder Code-Änderung |
| Funktionaler Smoke-Test (curl-Sequenz) | — | nach jedem abgeschlossenen Schritt |

Das sind reine Lesevorgänge bzw. reversible lokale Testläufe gegen Wegwerf-Datenbanken — kein Freigabebedarf.

## Was ausdrückliche Freigabe braucht (nie automatisch)

| Aktion | Warum |
|---|---|
| `git commit` | verändert dauerhafte Historie |
| `git push` | verändert geteilten/entfernten Zustand |
| Jede Strukturentscheidung (neue Datei außerhalb bestehender Konventionen, neues Modul) | braucht BAUSTEIN-Protokollzeile + Romans Entscheidung bei fehlendem a/b/c-Auslöser |
| Löschen von Dateien/Ordnern, die die KI nicht selbst in dieser Session erzeugt hat | z. B. `manus-packet` — nur auf explizite Bestätigung |
| Erweiterung des Funktionsumfangs über die aktuelle Schritt-Definition hinaus | §3.4: aus einer Aufgabe entsteht nie ungefragt eine zweite |

## Was Roman selbst im Terminal/Browser starten muss

Die KI kann den Server zu Testzwecken selbst kurz starten und wieder stoppen (passiert laufend in dieser Session). Für Romans eigene Nutzung/Demo/Abnahme muss er selbst aktiv werden — dafür gibt es keinen Ersatz:

| Aktion | Befehl / Ort | Voraussetzungen | Erwartetes Ergebnis | Fehlerverhalten | Rücksetzung |
|---|---|---|---|---|---|
| Server lokal starten | `go run ./cmd/lab-studio` in der Projektwurzel, **oder** VS Code: Terminal → Run Task → „RAP-LAB: Server starten" | Go 1.27+ installiert, Port 8090 frei | Konsole zeigt `RAP-LAB laeuft auf http://127.0.0.1:8090` | Fehlermeldung + Programmabbruch, z. B. bei fehlendem Schreibrecht auf `data/` oder belegtem Port | `Strg+C` im Terminal, keine Nebenwirkungen |
| Seite ansehen/prüfen | Browser → `http://127.0.0.1:8090` | Server läuft (siehe oben) | Homepage lädt | Verbindung abgelehnt → Server läuft nicht | — |
| Tests selbst nachvollziehen | `go test ./... -v` in der Projektwurzel | — | `PASS` je Test | `FAIL` mit Zeilenangabe | — |
| Commit erstellen/pushen | `git commit` / `git push` | — | — | — | — |

## Trennung automatisch vs. manuell — Zusammenfassung

```
Automatisch (KI, keine Freigabe):  build → vet → test → Smoke-Test
Manuell (Roman, keine KI-Ausführung): Server starten, Browser pruefen, git push, Loeschungen ausserhalb der Session, Strukturfreigaben
```

# 03 — Verzeichnisinventar: Zweck, Verwendung, Verantwortung


Vollständiges, wahrheitsgetreues Inventar der RAP-LAB-Verzeichnisstruktur, Stand 2026-10-17. Jede Zeile ist gegen den realen Dateibestand geprüft (`find`, `wc -l`, `git log`), nicht aus dem Gedächtnis geschrieben.

**Spaltenbedeutung:**
- **Eigentümer** — wer über Inhalt/Scope entscheidet (§3.1 zusammenarbeit-gesetze.md: immer Roman)
- **Schreiber** — wer die Datei tatsächlich erzeugt/ändert
- **Leser** — wer/was die Datei zur Entwicklungs- oder Laufzeit liest
- **Auslöser** — was einen Schreibvorgang auslöst
- **Zweck** — wofür die Datei da ist
- **Prüfverfahren** — wie man feststellt, dass der Inhalt korrekt ist

## Code

**Stand nach Funktionsachsen-Zerlegung** (s. `05-main-go-zerlegung.md`): `main.go` ist auf reines Wiring reduziert (25 Zeilen), die Logik liegt in sieben nach Achse getrennten Dateien, alle im selben Package `main`.

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `cmd/lab-studio/main.go` | Roman | KI | `go build`/`run`, Roman (Review) | neue Route/neuer Start-Parameter | Wiring: Store/Limiter erzeugen, Routen registrieren, Server starten | `go build`, manueller Start |
| `cmd/lab-studio/handlers.go` | Roman | KI | `go build`, Roman (Review) | neue Route | HTTP-Handler für alle fünf Routen + `currentUser` | Smoke-Test (curl-Sequenz) |
| `cmd/lab-studio/store.go` | Roman | KI | `go build`/`test`, Roman (Review) | neue Datenanforderung | SQLite-Persistenz, einzige Datei mit SQL | Neustart-Test, `go test` |
| `cmd/lab-studio/password.go` | Roman | KI | `go test`, Roman (Review) | Sicherheitsanforderung | Argon2id-Hashing | `go test ./... -v` (Rundlauf + Korrupt-Hash) |
| `cmd/lab-studio/csrf.go` | Roman | KI | `go build`, Roman (Review) | CSRF-Lücke (Plan Schritt 3) | CSRF-Token-Lebenszyklus + Validierung | Always-Pass-Test (403 ohne/mit falschem Token) |
| `cmd/lab-studio/ratelimit.go` | Roman | KI | `go test`, Roman (Review) | Brute-Force-Lücke (Plan Schritt 4) | In-Memory-Rate-Limiter | `go test ./... -v` + 6-Versuche-Live-Test |
| `cmd/lab-studio/validate.go` | Roman | KI | `go build`, Roman (Review) | Eingabevalidierung | E-Mail-/Passwort-Regeln | manueller INVALID-Test |
| `cmd/lab-studio/templates.go` | Roman | KI | `go build`, Roman (Review) | neues Formularfeld/neue Sektion | HTML/CSS, Rendering | Browser-/curl-Abgleich gegen `01-spezifikation.md` |
| `cmd/lab-studio/main_test.go` | Roman | KI | `go test` | neue Funktion mit Prüfbedarf | automatisierte Tests: Argon2id-Rundlauf, Korrupt-Hash, Rate-Limiter-Grenzen | `go test ./... -v` → PASS |

## Modul-Metadaten

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `go.mod` | Roman | **`go`-CLI selbst** (`go get`/`go mod init`), nie von Hand editiert | `go build`/`run`/`test` | neue Dependency | Modulname, Go-Version, direkte Abhängigkeiten | `go build ./...` schlägt fehl, wenn inkonsistent |
| `go.sum` | Roman | `go`-CLI selbst | `go build`/`run`/`test` | wie `go.mod` | Checksummen aller (auch transitiven) Abhängigkeiten | dito |

## Governance / Dokumentation

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `README.md` | Roman | KI | Roman, GitHub-Weboberfläche | jeder abgeschlossene Schritt | Einstiegspunkt: Startbefehl, geprüfter Stand, offene LÜCKEN | Startbefehl tatsächlich ausführen — stimmt er mit der realen Route überein? |
| `hub/rules/01-spezifikation.md` | Roman | KI | Roman, KI (Frontend-Referenz) | Referenzbild-Analyse (erste Kette, Schritt 1/5) | Homepage-Spezifikation, LÜCKEN-Liste zu Assets/Farben/Fonts | Abgleich gegen `content-home`-Template in `main.go` |
| `hub/rules/CONTRACT.md` | Roman | KI | Roman, KI (vor jeder neuen Datei) | Strukturentscheidung (zweite Kette) | verbindlicher Pfadvertrag für alle Verzeichnisse | jede neue Datei gegen diese Tabelle prüfen |
| `hub/rules/02-sicherheitserweiterungen-plan.md` | Roman | KI | Roman, KI | Sicherheitsaudit nach Manus-Fragment-Review | 4-Schritte-Plan (Persistenz/Argon2id/CSRF/Rate-Limit) + Nachweis je Schritt | genannte Commit-Hashes in `git log` auffindbar? |
| `hub/rules/03-verzeichnisinventar.md` (diese Datei) | Roman | KI | Roman, KI | diese Kette, Schritt 1/5 | vollständiges Datei-für-Datei-Inventar | Abgleich gegen `find .` |

## Startsteuerung / Werkzeugkonfiguration

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `.vscode/tasks.json` | Roman | KI | VS Code | Startkomfort ohne Terminal-Tippen | kapselt `go run ./cmd/lab-studio` als Task | Terminal → Run Task → Server startet, Browser zeigt Seite |
| `.gitignore` | Roman | KI | `git` | Runtime-Dateien sollen nicht versioniert werden | schließt `data/`, `logs/`, `*.exe` von Commits aus | `git status` nach Serverlauf zeigt diese Pfade nicht als untracked |

## Noch nicht existierende Ordner (Regel, kein Platzhalter)

`web/`, `configs/`, `scripts/` existieren aktuell nicht — sie wurden am 2026-10-17 als leere Platzhalter wieder entfernt (keine Nutzung seit Anlage, kein a/b/c-Auslöser). Nach dem neuen Vertrag (`CONTRACT.md`, Regel 1/2) entstehen sie erst mit ihrem ersten echten Inhalt.

| Ordner | Status |
|---|---|
| `logs/` | existiert, aber **LÜCKE**: `main.go` schreibt aktuell nicht hierhin — nur `log.Println` nach stdout. Füllt sich nur bei manueller Umleitung (`> logs/x.log`) |

## Laufzeit-Daten (nicht versioniert, nur lokal)

| Datei | Eigentümer | Schreiber | Leser | Auslöser | Zweck | Prüfverfahren |
|---|---|---|---|---|---|---|
| `data/raplab.sqlite` (+ `-wal`/`-shm`) | Roman | **der laufende Go-Prozess selbst**, weder KI noch manuell | Go-Prozess (SQLite-Treiber) | Server-Start + jede Registrierung/Login/Löschung | persistente Nutzerdaten | Neustart-Test: Login nach Neustart funktioniert noch? |

## Entry Point

Einziger Einstiegspunkt: `cmd/lab-studio/main.go`, `func main()`. Startbefehl: `go run ./cmd/lab-studio` (Projektwurzel). Kein zweites `cmd/`-Unterverzeichnis, kein zweites CLI-Tool — bestätigt durch `grep -rl "func main" --include="*.go" .` in der vorherigen Kette.

## Zusammenfassung der LÜCKEN

1. `logs/` hat keine tatsächliche Schreibanbindung im Code; der Vertrag in `CONTRACT.md` beschreibt eine Absicht, keine Implementierung.
2. Kein automatisiertes Log-Rotations- oder Aufräumverfahren für `data/raplab.sqlite` — manuell zu beobachten.


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

# 01 — Referenzanalyse und Spezifikation: RAP-LAB Homepage

**Ursprüngliche Vorgabe:** „Erstelle eine Golang-Fullstack-Website mit exakt einer Homepage nach der Bildreferenz.“
**Backreferenz:** keine — Ausgangspunkt.
**Referenzquellen:**
- `homepage-bsp.jpeg` — Mobile-Screenshot eines laufenden Prototyps auf `127.0.0.1:8090` (maßgebliche Referenz für Layout/Text)
- `krone-logo.jpg` — freigestelltes grünes Neon-Kronen-Icon, schwarzer Hintergrund
- `unterschrift-ledercover2.jpeg` — hochauflösende Quelle des „RAP-LAB“-Signature-Schriftzugs (Gold + grüner Neon-Glow), auf Leder-Textur

---

## 1. Ziel, Kontext, Zielgruppe

- **Ziel der Seite:** Landingpage für RAP-LAB — eine Plattform, die Musikschaffende (Rap/Musikproduktion) mit Kollaborateuren (Producer, Mixer, Videograf) und Brancheninformationen zusammenbringt, statt dass der Prozess aus verstreuten Einzelteilen (Notizen-App, Instagram, Discord) zusammengesucht wird.
- **Kontext:** Erste von geplant fünf Schritten — hier ausschließlich die Homepage. Zwei Nav-Punkte (`LAB-STUDIO`, `LAB-CONNECT`) deuten künftige Unterbereiche an, sind aber laut Auftrag **nicht** Teil der zu bauenden Seite.
- **Zielgruppe:** Rapper/Musiker, die ihren Produktionsprozess professionalisieren wollen.
- **Grenze für Schritt 1–2:** exakt eine Seite, eine Route, statisch (kein Formular-State nötig).

---

## 2. Layout — Sektionen von oben nach unten

1. **Header** (schwarz, über dünner Trennlinie von Hero abgesetzt)
   - links: Kronen-Logo (klein, Icon-Größe)
   - rechts: Navigation `LAB-STUDIO` · `LAB-CONNECT`, goldene Caps-Schrift
2. **Hero**
   - Signature-Schriftzug „RAP-LAB“ (Bild-Asset, Gold + grüner Neon-Glow), zentriert
   - Tagline, zweizeilig, Gold, Caps, zentriert:
     > EIN ORT FÜR DAS, WAS ZWISCHEN IDEE UND MUSIK PASSIERT.
3. **Punkte-Liste** (weiß, fett, Caps, zentriert, je Zeile ein eigener Absatz):
   ```
   MUSIK ENTWICKELN.
   MENSCHEN FINDEN.
   ZUSAMMENARBEITEN.
   WEITERKOMMEN.
   ```
4. **Fließtext-Sektion** — leicht abgesetzter Hintergrundton (dunkelgrau statt reines Schwarz), linksbündig, weiß, Normaltext. Exakter Text (wortgleich aus Referenz, nicht umformulieren):
   ```
   Denn mal ehrlich:

   Du kannst den nächsten Song auch wieder in deinen Handy-Notizen schreiben.
   Den Producer kannst du wieder über fünf Ecken suchen.
   Den Mixer wieder in irgendeiner Instagram-Story finden.
   Den Videografen wieder irgendwo in einem Discord-Server suchen.
   Und deine Brancheninformationen kannst du dir wieder aus zehn verschiedenen Quellen zusammensuchen.
   Oder du hörst irgendwann auf, deinen Musikprozess aus Einzelteilen zusammenzubauen.
   ```
5. **Abschluss:** `RAP-LAB.` zentriert, Gold, Caps.

Reihenfolge und Textinhalt sind 1:1 aus `homepage-bsp.jpeg` übernommen, nichts ergänzt.

---

## 3. Farben, Typografie (benannt, soweit aus dem Bild ableitbar)

| Rolle | Wert | Status |
|---|---|---|
| Hintergrund Haupt | Schwarz | gesichert, exakter Hex-Ton = LÜCKE |
| Hintergrund Fließtext-Sektion | dunkelgrau, etwas heller als Haupt-Schwarz | gesichert, exakter Hex-Ton = LÜCKE |
| Akzentfarbe 1 (Logo-Schrift, Nav, Taglines, Abschluss) | Gold/warmes Beige | gesichert, exakter Hex-Ton = LÜCKE |
| Akzentfarbe 2 (Krone, Glow unter Signature) | Neon-Grün | gesichert, exakter Hex-Ton = LÜCKE |
| Fließtext/Listenfarbe | Weiß / sehr helles Grau | gesichert, exakter Hex-Ton = LÜCKE |
| Schriftart Caps-Texte (Nav, Tagline, Liste, Abschluss) | fette, serifenlose Grotesk | **LÜCKE** — Name nicht zweifelsfrei identifizierbar |
| Schriftart Fließtext | serifenlose Normalschrift | **LÜCKE** — Name nicht zweifelsfrei identifizierbar |
| Schriftzug „RAP-LAB“ (Hero) | individuelle Signature/Handschrift | kein Web-Font — **muss als Bild-Asset eingebunden werden** |

---

## 4. LÜCKEN — explizit offene Punkte

1. Exakte Hex-Farbwerte (nur aus Screenshot geschätzt, keine Farbpipette gezogen).
2. Exakte Font-Namen für Caps- und Fließtext — Ersatz: System-Font-Stack mit ähnlichem Charakter, bis Roman eine Quelle nennt.
3. „RAP-LAB“-Signature als freigestelltes Web-Asset (transparentes PNG/SVG) existiert noch nicht — `unterschrift-ledercover2.jpeg` hat Leder-Hintergrund, ist nicht direkt einsetzbar. Muss freigestellt oder zugeschnitten werden.
4. Kronen-Logo (`krone-logo.jpg`) hat durchgehend schwarzen Hintergrund, kein Alpha-Kanal — unproblematisch nur wenn Seiten-Hintergrund exakt gleiches Schwarz ist.
5. Zielrouten von `LAB-STUDIO` / `LAB-CONNECT` unbekannt — werden in diesem Umfang als nicht-funktionale Platzhalter (`#` oder `href="#"`) behandelt, bis Roman Ziele nennt.
6. Nur Mobile-Referenz vorhanden, keine Desktop-Referenz — Desktop-Breakpoint-Verhalten wird aus Mobile-Layout abgeleitet (Zentrierung, max-width-Container), nicht aus einer zweiten Bildvorlage.
7. Unklar, ob `RAP-LAB.` das tatsächliche Seitenende ist oder der Screenshot vor weiterem Inhalt endet (Scroll-Grenze unbekannt) — wird als Seitenende angenommen, bis Roman widerspricht.
8. Meta-Daten (Seitentitel, Favicon, Description) nicht im Referenzbild enthalten.

---

## 5. Mindestkriterien (für Schritt 2)

- Exakt eine HTML-Seite, eine Route.
- Alle 5 Layout-Sektionen aus §2 vorhanden, in dieser Reihenfolge.
- Alle Texte wortgleich aus §2 übernommen.
- Hintergrund durchgehend schwarz (bis auf abgesetzte Fließtext-Sektion).
- Responsive: bricht weder auf Mobile- noch auf Desktop-Breite.
- Keine Platzhalterbilder/Lorem Ipsum — Logo/Signature als echte Bild-Assets referenziert (Pfad ggf. als LÜCKE markiert, bis Dateien in LABv2 kopiert sind).

## 6. Verbote

- Keine zweite Seite oder Route in Schritt 2.
- Keine weiteren Navigationspunkte außer `LAB-STUDIO`, `LAB-CONNECT`.
- Keine erfundenen Hex-Werte ohne Kennzeichnung als Schätzung.
- Keine CSS-Frameworks (Bootstrap/Tailwind) oder JS-Frameworks ohne benannten Auslöser (a/b/c nach § 6.3 zusammenarbeit-gesetze.md) — Default ist reines HTML/CSS.
- Keine externen Font-/Asset-CDNs ohne Rückfrage.
- Kein Auth-/Backend-Code in Schritt 1–2.

---

**Ergebnis dieses Schritts:** überprüfbare Spezifikation, s.o. Nächster Schritt (2/5) erst nach Freigabe durch Roman.
