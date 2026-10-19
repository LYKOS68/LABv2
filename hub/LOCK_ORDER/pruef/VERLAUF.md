# Prüfverlauf `/hub/LOCK_ORDER/pruef/`

Append-only. Nur tatsächlich ausgeführte Prüfungen werden hier abgelegt.
Ein Eintrag pro Lauf. Frühere Einträge werden nie überschrieben oder gelöscht — nur angehängt.

- **Kennung:** `PR-<YYYY-MM-DD>-<pruefid>-<stand>-L<n>` (n = laufende Nummer je Tag)
- **Statuswerte:** `erfolgreich` | `fehlgeschlagen` | `nicht ausgeführt`
- **Pflichtangaben je Eintrag:** Zeitpunkt (UTC), Kennung, Werkzeug, Entwicklungsstand, Ergebnis je Prüfung inkl. Exit-Code, Fehlerbefunde

## Geplant (nicht ausgeführt)

Keine. Geplante Prüfungen erscheinen hier zuerst und werden erst dann in „Ausgeführt" übernommen, wenn sie gelaufen sind.

## Ausgeführt

### L0 — 2026-10-18 (Uhrzeit nicht protokolliert) — `PR-2026-10-18-smoke-LOCK-2026.10.18-1-L0`

- Werkzeug: manuell (Modell, Termux) — Kommandos aus LOCK_ORDER.md §3
- Stand: LOCK-2026.10.18-1
- Ergebnis: S1–S5 alle `erfolgreich` (Exit 0), Details in LOCK_ORDER.md §3 (dort dokumentiert, bevor dieser Prüfbereich existierte)
- Fehlerbefunde: keine

### L1 — 2026-10-19T00:15Z — `PR-2026-10-19-smoke-LOCK-2026.10.19-1-L1`

- Werkzeug: `sh hub/LOCK_ORDER/tools/smoke.sh` (read-only)
- Stand: LOCK-2026.10.19-1
- Exit-Code gesamt: 0

| Prüfung | Ausgabe | Status |
|---|---|---|
| S1 | Regelwurzel vorhanden | erfolgreich |
| S2 | md-Dateien nur an erlaubten Orten | erfolgreich |
| S3 | txt-Inventar unverändert | erfolgreich |
| S4 | cmd/ unverändert | erfolgreich |
| S5 | Branch synchron mit origin, kein Push (`## termux-pull...origin/termux-pull`) | erfolgreich |

- Fehlerbefunde: keine
- Hinweis: Ausführung im Zuge der Strukturerweiterung Prüfbereich/Werkzeugbereich; vorherige Schmutzliste unverändert (`test2.txt` gelöscht uncommittet, `hub/git-befehlspal.txt`, `branch-termux-pull-manifest.md`, `hub/LOCK_ORDER/` untracked)

### L2 — 2026-10-19T00:20:43Z — `PR-2026-10-19-smoke-LOCK-2026.10.19-2-L2`

Anlass: Sprach- und Paket-Audit plus minimaler Korrekturlauf (Auftrag 1/5–4/5).

Audit (bestätigte Befunde, je mit Beleg):

- Hauptprojekt: Go, Modul `rap-lab`, `go 1.27.0` (`go.mod:1-3`); Toolchain `go1.27.0 linux/arm64` (`go version`); Testwerkzeuge `go test ./...`, `gofmt`, `go vet` (`.github/workflows/test.yml:18-25`, `*_test.go` unter `cmd/lab-studio/`).
- Werkzeuge: POSIX-`sh` (dash), `git 2.51.0`, `find`/`grep`/`sed`/`sort`/`head` — ohne Paketabhängigkeiten, kompatibel zur Laufzeit linux/Termux; Pfade relativ zur Git-Root.
- Instabilität extern: erste `go version`-Ausführung nach Toolchain-Download mit Segfault beendet, zweite Ausführung erfolgreich — reproduziert (2 Läufe), keine Korrektur (extern, kein Auftrag).
- Vermutung (nicht verifiziert): `.sh` läuft auf dem Windows-PC nur mit Git-Bash/WSL — offen, keine Änderung ohne Bedarf.
- Nachgewiesene Abweichung (korrigiert): `smoke.sh` erzeugte keine maschinenlesbaren Ergebnisse (nur Textausgabe, Nachweis: Version vom 2026-10-18, sha256 `15382274c40b8425e884b010f96bd86f04b26e8a870add57a03770e585932812`); Sicherung vor Eingriff: Kopien in `/tmp/opencode/lock-backup/`, sha256 der betroffenen Dateien: `LOCK_ORDER.md` `dfe5dcc3…`, `VERLAUF.md` `173eb9e3…`, `WERKZEUGE.md` `77cc61fe…`, `smoke.sh` `15382274…`; Korrektur: `--json <PR-Kennung>` ergänzt (eine Änderung, Stand `LOCK-2026.10.19-2`).

Lauf: Werkzeug `sh hub/LOCK_ORDER/tools/smoke.sh --json PR-…-L2` (Echtlauf, Umleitung nach `pruef/PR-2026-10-19-smoke-LOCK-2026.10.19-2-L2.json`, mit `python3 -m json.tool` validiert).

| Prüfung | Status | Exit-Code |
|---|---|---|
| S1–S5 | alle `PASS` | 0 (Gesamt) |

Fehlerbefunde: keine (neben dem externen Go-Segfault, siehe oben).

### L3 — 2026-10-19T00:20:58Z — `PR-2026-10-19-smoke-LOCK-2026.10.19-2-L3`

Anlass: Fehlerpfad-Prüfung (bewusst ungültige Eingabe, Auftrag 4/5).

- Vorbereitung: temporäre `smoke-probe.md` in der Root angelegt (außerhalb erlaubter md-Orte) — danach wieder entfernt.
- Werkzeug: `sh hub/LOCK_ORDER/tools/smoke.sh --json PR-…-L3` → Ergebnisdatei `pruef/PR-2026-10-19-smoke-LOCK-2026.10.19-2-L3.json` (echter Lauf, validiert).
- Ergebnis: S1 `PASS`, S2 `FAIL` („unerlaubte md-Dateien: smoke-probe.md"), S3–S5 `PASS`, Exit-Code 1.
- Zusätzlich geprüft: `--json` ohne Prüfkennung → Exit 2; unbekanntes Argument → Exit 2 (keine Ergebnisdatei erzeugt).
- Abschlußlauf nach Entfernen der Probe: Textmodus S1–S5 alle `PASS`, Exit 0 — Zustand wiederhergestellt.
- Statusverteilung belegt: `PASS`, `FAIL` und Eingabefehler (`FAIL EINGABE`, Exit 2) tatsächlich erzeugt; `NOT_RUN` nicht verwendet (kein Lauf ohne Prüfung).
- Fehlerbefunde: keine neuen.

### L4 — 2026-10-19T01:57:45Z — `PR-2026-10-19-smoke-LOCK-2026.10.19-3-L4`

Anlass: `hub/` als zentrale Arbeitswurzel deklariert, Zielstruktur `knowledge/`, `user/`, `shared/` angelegt (Auftrag 1/5–5/5).

- Regeländerung: LOCK_ORDER.md §2.1–§2.4 (hub-Deklaration, Einstiegspunkte, Zielstruktur, Lesebereiche), S2 um `hub/knowledge/`, `hub/user/`, `hub/shared/` erweitert, Stand `LOCK-2026.10.19-3`; `smoke.sh` S2-Liste mitgeändert (eine Iteration).
- Verzeichnisse angelegt: `hub/knowledge/` (leer), `hub/shared/` (leer); `hub/user/` existierte mit `prompt-vorlage.yaml`. Keine Inhalte verschoben — Überführung von `brand/` und `mitglieder.csv` bleibt Freigabe vorbehalten.
- Werkzeug: `sh hub/LOCK_ORDER/tools/smoke.sh --json PR-…-L4` (Echtlauf, `pruef/PR-2026-10-19-smoke-LOCK-2026.10.19-3-L4.json`, validiert).

| Prüfung | Status | Exit-Code |
|---|---|---|
| S1–S5 | alle `PASS` | 0 (Gesamt) |

Fehlerbefunde: keine. Bestand unverändert außer den zwei neuen leeren Verzeichnissen; `cmd/` unberührt; kein Commit, kein Push.

### L5 — 2026-10-19T02:07:00Z — `PR-2026-10-19-smoke-LOCK-2026.10.19-4-L5`

Anlass: Überführung freigegeben und ausgeführt — `hub/brand/` → `hub/knowledge/brand/`, `hub/mitglieder.csv` → `hub/user/mitglieder.csv` (Auftrag/„ja" des Menschen).

- Sicherung vor Eingriff: Kopie `LOCK_ORDER.md` in `/tmp/opencode/lock-backup2/`; sha256 vorher: `LOCK_ORDER.md` `3fa50fe5…`, `smoke.sh` `6391f1db…`, `branch-termux-pull-manifest.md` `1e782bae…`, `README.md` `e495ecc3…`.
- Verschiebung via `git mv` (21 brand-Dateien + `mitglieder.csv`); Inhalte angepasst: `LOCK_ORDER.md` (§1 Befunde, §2.1 Original-Pfad, §2.3 Status, §3 S2), `tools/smoke.sh` (S2-Liste), `hub/knowledge/brand/README.md` (Titel), `branch-termux-pull-manifest.md` (§5 Baum, §6 Pfadregel, §7 Zielstruktur, Zusammenfassung), `README.md` („Wo was liegt" + `hub/`-Abschnitt). Eine Iteration, Stand `LOCK-2026.10.19-4`.
- Werkzeug: `sh hub/LOCK_ORDER/tools/smoke.sh --json PR-…-L5` (Echtlauf, `pruef/PR-2026-10-19-smoke-LOCK-2026.10.19-4-L5.json`, validiert).

| Prüfung | Status | Exit-Code |
|---|---|---|
| S1–S5 | alle `PASS` | 0 (Gesamt) |

Fehlerbefunde: keine. `cmd/` unberührt; kein Commit, kein Push.

### L6 — 2026-10-19T02:24:13Z — `PR-2026-10-19-smoke-LOCK-2026.10.19-5-L6`

Anlass: Zugriffs- und Herkunftsregel festgeschrieben („es arbeitet keiner außer Roman und seinen CLIs an diesem Repository, egal wann").

- Regeländerung: `LOCK_ORDER.md` §2.5 neu (exklusiver Zugriff: nur Roman und seine CLI-Werkzeuge, keine anderen Personen/Modelle/Agenten/Bots/Dienste, kein Automatismus, Ausnahmen nur einzeln durch Roman) und §5 ergänzt; `branch-termux-pull-manifest.md` §4 um dieselbe Regel mit Verankerung auf §2.5 erweitert. Eine Iteration, Stand `LOCK-2026.10.19-5`.
- Werkzeug: `sh hub/LOCK_ORDER/tools/smoke.sh --json PR-…-L6` (Echtlauf, `pruef/PR-2026-10-19-smoke-LOCK-2026.10.19-5-L6.json`, validiert).

| Prüfung | Status | Exit-Code |
|---|---|---|
| S1–S5 | alle `PASS` | 0 (Gesamt) |

Fehlerbefunde: keine. `cmd/` unberührt; kein Commit, kein Push.
