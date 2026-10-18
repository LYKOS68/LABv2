# Orientierung für jede KI, die in diesem Repo arbeitet

Knapp, kein Prosa-Dump. Die volle Methodik (ARBEITSWEISE.md, zusammenarbeit-gesetze.md, roman_code_review_instructions.md) liegt bewusst außerhalb dieses Repos — Roman nennt den Pfad bei Bedarf. Hier nur das Betriebs-Minimum, damit eine neue KI ohne Nacherzählen anschlussfähig ist.

## Regeln

1. Roman schreibt keinen Code — das macht ausschließlich die KI (siehe `hub/rules/04-arbeitsvertraege.md`).
2. Keine Struktur ohne Beleg. Pfade: `hub/rules/CONTRACT.md` (zwei Wurzeln, vier Regeln).
3. Jede Strukturentscheidung protokolliert:
   `BAUSTEIN: <name> | AUSLÖSER: <a/b/c oder "Roman-Entscheidung außerhalb §6.3"> | BELEG: <Zahl/Datei:Zeile> | INVARIANTE: <Prädikat>`
4. Vor jedem Bauschritt: Ziel, Mindestkriterien, Prüfverfahren schriftlich festhalten — dann erst bauen.
5. „Fertig" = ausgeführt und geprüft (Build/Test/Smoke-Test), nicht behauptet.
6. Kein Code ohne Test gegen INVALID-Eingabe, falls es ein Gate/Validator ist (Always-Pass-Falle vermeiden).
7. Kein leerer Ordner auf Vorrat. Ein Ordner entsteht erst mit seinem ersten echten Inhalt.

## Aktueller Projektstand (kurz)

Stand 2026-10-18 (Tag `stand-2026-10-18`). Go-Webanwendung (`cmd/lab-studio/`), SQLite, kein JS-Framework. Start: `go run ./cmd/lab-studio` → `127.0.0.1:8090`.

- **Startseite** im Markendesign (`10`), **Konto** mit E-Mail-Bestätigung und Passwort vergessen (`07`, `08`), Sicherheit (`02`)
- **LAB-STUDIO** `/studio`: Arbeitsfläche, Notizen, Songs, Lyrics-Editor — pro Konto, leer zum Start (`09`)
- **LAB-CONNECT** `/connect`: Oberflächen-Prototyp, Beispieldaten, speichert nichts (`11`)
- Marke: `12-brand-spec.md`, Originale in `hub/brand/`. Zweck und Ausbauplan: `13-zweck-und-vision.md`. Inventar: `14-verzeichnisinventar.md`
- 37 Tests; GitHub prüft bei jedem Push automatisch

Roman liest Code mit und achtet auf Lesbarkeit: kurze Funktionen, eine Prüfung pro Zeile, nichts doppelt kopieren. Antworten an Roman in einfacher Sprache.

## Dateien in `hub/`

- `rules/` — eingefrorene Spezifikationen/Verträge (nicht mehr verhandelbar ohne neue Freigabe)
- `chatroom.jsonl` — Nachrichten zwischen Agents/Roman, append-only, Format in `chatroom.schema.json`. Zum Hinterlassen einer Nachricht: eine Zeile anhängen, nichts umschreiben.
- `drafts/` — unfertige Ideen, **keine Bauaufträge**. Erst nach expliziter Freigabe durch Roman umsetzen.
- `agent-docs/` — diese Datei und weitere, bei Bedarf.

## Wo was liegt

- `hub/rules/` — eingefrorene Spezifikationen/Verträge
- `hub/drafts/` — unfertiges Denken (jederzeit verwerfbar)
- `cmd/lab-studio/` — der Code selbst, siehe `hub/rules/CONTRACT.md` für den internen Dateivertrag
