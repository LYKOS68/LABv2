# Orientierung für jede KI, die in diesem Repo arbeitet

Knapp, kein Prosa-Dump. Die volle Methodik (ARBEITSWEISE.md, zusammenarbeit-gesetze.md, roman_code_review_instructions.md) liegt bewusst außerhalb dieses Repos — Roman nennt den Pfad bei Bedarf. Hier nur das Betriebs-Minimum, damit eine neue KI ohne Nacherzählen anschlussfähig ist.

## Regeln

1. Roman schreibt keinen Code — das macht ausschließlich die KI (siehe `docs/04-arbeitsvertraege.md`).
2. Keine Struktur ohne Beleg. Pfade: `docs/CONTRACT.md`.
3. Jede Strukturentscheidung protokolliert:
   `BAUSTEIN: <name> | AUSLÖSER: <a/b/c oder "Roman-Entscheidung außerhalb §6.3"> | BELEG: <Zahl/Datei:Zeile> | INVARIANTE: <Prädikat>`
4. Vor jedem Bauschritt: Ziel, Mindestkriterien, Prüfverfahren schriftlich festhalten — dann erst bauen.
5. „Fertig" = ausgeführt und geprüft (Build/Test/Smoke-Test), nicht behauptet.
6. Kein Code ohne Test gegen INVALID-Eingabe, falls es ein Gate/Validator ist (Always-Pass-Falle vermeiden).

## Aktueller Projektstand (kurz)

Go-Fullstack-Webserver (`cmd/lab-studio/`), eine Homepage, Registrierung/Login/Logout/Account-Löschung. SQLite-Persistenz, Argon2id-Hashing, CSRF-Schutz, Rate-Limiting — alle vier geprüft und eingefroren (`docs/02-sicherheitserweiterungen-plan.md`). Start: `go run ./cmd/lab-studio`, läuft auf `127.0.0.1:8090`.

## Dateien in `hub/`

- `chatroom.jsonl` — Nachrichten zwischen Agents/Roman, append-only, Format in `chatroom.schema.json`. Zum Hinterlassen einer Nachricht: eine Zeile anhängen, nichts umschreiben.
- `drafts/` — unfertige Ideen, **keine Bauaufträge**. Erst nach expliziter Freigabe durch Roman umsetzen.
- `agent-docs/` — diese Datei und weitere, bei Bedarf.

## Wo was liegt

- `docs/` — eingefrorene Spezifikationen/Verträge (nicht mehr verhandelbar ohne neue Freigabe)
- `hub/drafts/` — unfertiges Denken (jederzeit verwerfbar)
- `cmd/lab-studio/` — der Code selbst, siehe `docs/CONTRACT.md` für den internen Dateivertrag
