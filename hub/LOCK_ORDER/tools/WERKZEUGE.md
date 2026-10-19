# Werkzeuge `/hub/LOCK_ORDER/tools/`

Registry aller Werkzeuge. Ein Skript wird **nur** bei nachgewiesenem Steuerungsbedarf angelegt;
jeder Eintrag braucht Zweck, Eingaben, Wirkung und Startbefehl.

Verbindliche Beschränkungen:

- Werkzeuge sind **read-only** gegenüber dem Projekt: keine Änderungen, keine Commits, keine Pushes, kein automatisches Schreiben von Dateien.
- Dateizugriff nur auf ausdrücklich erlaubte Pfade: Lesen der gesamten Root (außer `.git/`-Interna), Schreiben **nirgends**.
- Jede Ausführung wird anschließend als Eintrag in `../pruef/VERLAUF.md` protokolliert (vom Ausführenden, nicht vom Skript).
- Start immer manuell durch den Menschen oder das Modell nach Auftrag — keine Automatik, keine Cron-, Hook- oder Event-Auslösung.

## Registry

| Skript | Zweck | Eingaben | Wirkung | Startbefehl |
|---|---|---|---|---|
| `smoke.sh` | Führt die in `LOCK_ORDER.md` §3 definierten Smoke-Prüfungen S1–S5 aus und meldet je Prüfung PASS/FAIL mit Exit-Code | optional `--json <PR-Kennung>` (z. B. `PR-2026-10-19-smoke-LOCK-2026.10.19-2-L1`); sonst keine; arbeitet aus der Projekt-Root (via `git rev-parse`) | ausschließlich Konsolenausgabe; liest Dateien, verändert nichts; Exit 0 = alle grün, 1 = mind. 1 Prüfung fehlgeschlagen, 2 = ungültige Eingabe, 3 = keine Git-Root | `sh hub/LOCK_ORDER/tools/smoke.sh [--json <PR-Kennung>]` |

### JSON-Ergebnisformat (`--json`)

Ausgabe nur auf stdout; wird vom Ausführenden umgeleitet, z. B.
`sh hub/LOCK_ORDER/tools/smoke.sh --json PR-...-L1 > hub/LOCK_ORDER/pruef/PR-...-L1.json`.
Es entsteht nur bei tatsächlichem Lauf — es werden keine Ergebnisse vorab erzeugt.

Felder: `pruefkennung`, `werkzeug`, `stand` (aus `LOCK_ORDER.md` gelesen), `zeitpunkt_utc`,
`exit_code`, `ergebnisse[]` mit `id` (S1–S5), `status` (`PASS` | `FAIL`; `NOT_RUN` bleibt für
nicht ausgeführte Prüfungen reserviert und darf nie ohne Lauf erscheinen), `detail`, `betroffen[]`.
