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
