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
