# BRANCH : termux-pull : HIER LADE ICH VON UNTERWEGS ÜBERS HANDY PULL + VERBESSERUNG REIN UND KANN SIE SPÄTER VOM PC WIEDER LADEN:

#RAP-LAB

**Musik passiert nicht nur im Studio.**

## Warum es RAP-LAB gibt

Eine Idee wird zum Song. Der Song braucht einen Beat, die Vocals einen Mix, das Release ein Cover, das Projekt ein Video. Gesucht wird dann über WhatsApp, Instagram, Google, Discord und irgendwelche Gruppen: 17 Tabs, 40 Chats. Und die eigenen Lines liegen verteilt in Handy-Notizen, alten Notizblöcken und Dateien auf dem PC — gute Sätze, aber nicht rund, nicht zusammengebaut.

**Es fehlen nicht die Leute. Es ist alles verteilt.**

RAP-LAB holt das an einen Ort, von der ersten Line bis zur Zusammenarbeit:

```
Line → Song → Projekt → Producer ↔ Artist → Mix → Cover → Video → Release
```

## Für wen

Rapper, Artists, Songwriter, Producer, Beatmaker, Mixer, Videografen, Designer — alle, die an einem Musikprojekt arbeiten oder dafür gebraucht werden. Kein Label nötig, keine Millionen Follower: Es reicht, etwas zu können, das jemand anderes für sein Projekt braucht.

## Die Bereiche

- **LAB-STUDIO** — der Ort, an dem Lyrics geschrieben, verfeinert und archiviert werden. Eine Arbeitsfläche, auf der man Songs, Lines und Notizen nebeneinander legt und frei verschiebt, weil Schreiben wie Puzzeln ist. Kein Tonstudio: keine Beats, kein Mixing, kein Recording.
- **LAB-CONNECT** — Angebote und Gesuche der Musikbranche, wie Kleinanzeigen. Die Frage ist nicht „Kann ich alles alleine?", sondern „Was kann ich selbst, und wen brauche ich für den Rest?"
- **NEWS** (später) — relevante Infos für die eigene Arbeit statt Feed für alle.

## Wo es gerade steht (2026-10-19)

| Bereich | Stand |
|---|---|
| Startseite | im Markendesign, auf Handy und PC |
| Konto | Registrieren mit E-Mail-Bestätigung, Anmelden, Passwort vergessen, Konto löschen |
| LAB-STUDIO | nutzbar: Arbeitsfläche zum Verschieben und Zoomen, Notizen, Songs anlegen, Lyrics schreiben und speichern. Jedes Konto sieht nur seine eigenen Texte |
| LAB-CONNECT | nur die Oberfläche, mit Beispielanzeigen — speichert noch nichts |
| NEWS | noch nicht begonnen |

## Wohin es geht

- **LAB-STUDIO**: Songs getrennt nach *Released* und *Drafts*; ein **Reimbuch** zum Nachschlagen von Reimen; **Relatable Bars** als Ablage der stärksten Lines zum Wiederverwenden
- **LAB-CONNECT**: echte Anzeigen erstellen und nach Standort und Region filtern; Videos und Feedback an Anzeigen, damit man zeigen kann, was man kann
- **Community**: Gruppen, Meetups, Gigs und lokale Kontakte
- **Kreatives Profil**: „Wer bin ich? Was mache ich? Was kann ich zeigen?"

## Ausprobieren

```
go run ./cmd/lab-studio
```

Dann im Browser **http://127.0.0.1:8090** öffnen (in VS Code auch über die Aufgabe „RAP-LAB: Server starten"). Solange kein Mail-Dienst eingerichtet ist, erscheinen die Bestätigungs- und Passwort-Links im Konsolenfenster des Servers.

Prüfen mit `go test ./...` — GitHub macht das bei jedem Push automatisch.

## Wo was liegt

```
cmd/lab-studio/   das Programm
data/             die Datenbank, entsteht beim ersten Start, nie auf GitHub
hub/              alles, was kein Programm ist
  brand/          Logos, Schriftzüge und Design-Vorlagen im Original — nie verändern
  user/           nur für Roman: Vorlage für Prompts zum Kopieren
  chatroom.jsonl  Nachrichten zwischen den Modellen, die hier mitarbeiten
```

## Nächster Neubau: `hub/`

`hub/` wird neu aufgebaut. Die alten Regel- und Plandokumente wurden am 2026-10-18 entfernt, weil sie veraltet und teils von der KI dazugedichtet waren; ältere Fassungen liegen in der Git-Historie. Der Neubau ist **beschlossen, aber noch nicht umgesetzt**. Bis dahin gilt als Beleg nur, was im Code steht, was Git zeigt und was Roman selbst sagt.
