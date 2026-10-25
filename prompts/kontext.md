# Projektkontext LABv2

Jeder Prompt in `prompts/` setzt voraus, dass du diese Datei zuerst gelesen hast.

## Projekt

LABv2 ist eine Kleinanzeigen-Plattform für Musikschaffende ohne Label.
Der globale Bereich sind Kleinanzeigen: Gesuche und Angebote.
Ein Gruppenbereich wie skool.com kommt später, bis dahin nur als Platzhalter-Seite.

## Entwicklungsplan

| Phase | Präfix | Ziel |
|---|---|---|
| 1 | H | Homepage, das Herz der Seite |
| 2 | K | Kleinanzeigen: Konto, Anzeige erstellen, Liste, Detail, Kontakt |
| 3 | G | Gruppenbereich als Platzhalter-Seite |

Eine Phase beginnt erst, wenn alle Schritte der vorigen in `prompts/status.md` als BESTANDEN stehen.

## Technik

- Go 1.27, Modul `rap-lab`, nur Standardbibliothek, solange eine Aufgabe nichts anderes verlangt
- Windows 11, PowerShell 7
- Datenbank: PostgreSQL 18 (von Roman bestätigt 2026-10-25), lokal als Dienst, Verbindung über `DATABASE_URL` aus `.env`
- `psql` ist im VS-Code-Terminal über `.vscode/settings.json` im PATH
- Port: `PORT` aus `.env`, Standard 3000
- Node 24 nur für Werkzeuge
- Kein Python, nirgends

## Marke

- Ton: duzen, kurz, direkt
- Leitsatz: „Du musst kein Label sein. Du brauchst auch keine Millionen Follower. Du musst nur etwas können, das jemand anderes für sein Projekt braucht.“
- Gesuche: „Sag, was du wirklich brauchst.“
  - „Suche Producer für Deutschrap. Düsterer Sound. Langfristige Zusammenarbeit.“
  - „Suche Sängerin für Hook auf einem fertigen Song.“
  - „Suche Mixer für zwei fertige Tracks.“
  - „Suche Videograf für Musikvideo – Konzept steht.“
- Angebote: „Zeig, was du kannst.“
  - „Mix & Master – 50 € pro Song.“
  - „Beats für Rapper – individuelle Produktionen.“
  - „Ghostwriting – Hooks und Parts.“
  - „Musikvideo – Konzept, Dreh und Schnitt.“
  - „Coverdesign für Artists.“
- Name, Farben, Schrift, Logo: offen. Nicht erfinden, sichtbaren Platzhalter `[OFFEN: …]` setzen.

## Code-Regeln

- Kurze Funktionen mit sprechenden Namen
- Eine Prüfung pro `if`
- Keine kopierten Blöcke, Gemeinsames in eine Funktion
- Nur Dateien anfassen, die der Schritt nennt

## Arbeitsregeln für die CLI

**Chunking**
- Ein Schritt = ein Ziel, höchstens 3 Dateien, endet mit `go build ./...` ohne Fehler.
- Ein Schritt setzt nur auf Dateien auf, die in `prompts/status.md` als BESTANDEN geführt sind.

**Modalität**
- Erlaubt: Dateien lesen und schreiben, Befehle ausführen (`go`, `git`, `curl`), deren echte Ausgabe lesen.
- Nicht erlaubt: Aussehen im Browser, Bilder oder Audio beurteilen, ohne ein Werkzeug dafür zu haben. Dann steht im Ergebnis `MANUELL PRÜFEN: <was>`.
- Nie eine Befehlsausgabe ausdenken. Was nicht ausgeführt wurde, gilt als nicht geprüft.

**Abbruch**
- Fehlt ein Pflichtparameter: sofort `BLOCKIERT: <parameter>`, nichts bauen.
- Höchstens 3 Versuche pro Schritt. Danach `DURCHGEFALLEN` mit der letzten echten Fehlerausgabe.
- `BESTANDEN` nur mit Nachweis: Befehl und tatsächliche Ausgabe.
- Am Ende jedes Schritts eine Zeile an `prompts/status.md` anhängen:
  `| <ID> | <BESTANDEN/DURCHGEFALLEN/BLOCKIERT> | <Datum> | <Nachweis kurz> |`

## Rollen

- Roman entscheidet und gibt frei.
- Die CLI führt genau den Schritt aus, der ihr gegeben wird.
