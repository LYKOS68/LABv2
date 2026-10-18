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
