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
