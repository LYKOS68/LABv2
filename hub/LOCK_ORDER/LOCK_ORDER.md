# LOCK_ORDER — verbindliche Architekturregeln

Subwurzel `/hub/LOCK_ORDER/` (deklariert 2026-10-18, Auftrag 1/5).
Geltung: Projekt-Root `/root/github/LABv2`, Branch `termux-pull`.
Regeln liegen als Artefakte in dieser Datei, nicht im Kontextgedächtnis.

**Entwicklungsstand: `LOCK-2026.10.19-2`** — bestätigt durch die Smoke-Prüfungen (Ausführungsnachweise in `pruef/VERLAUF.md`, maschinenlesbar als `pruef/PR-*.json` aus echten Läufen).

---

## 1. Formatordnung (verbindlich)

Gilt für alle Nicht-Runtime-Dateien: Projekt-Root-Ebene und `hub/`.
Runtime unter `cmd/lab-studio/` (Go, HTML, CSS, JS, Fonts, Bilder) ist ausgenommen.

| Format | Zweck | Zielgruppe |
|---|---|---|
| `.md` | Projektregeln, globale Architektur, menschliche Dokumentation | alle |
| `.txt` | ausschließlich Inhalte für Roman | Roman |
| `.json` / `.jsonl` / `.yaml` | KI-Daten, Konfiguration, maschinenlesbare Vgaben | Maschinen/KI |
| `.sh` | ausschließlich in `hub/LOCK_ORDER/tools/` — read-only Werkzeuge, manuell auslösbar | Mensch/Modell |

Prüfung vorhandener Dateien (Bestand, nichts verschoben oder konvertiert — kein Bedarf nachgewiesen):

| Datei | Format | Befund |
|---|---|---|
| `README.md`, `branch-termux-pull-manifest.md` | md | konform (Projekt-/Branch-Doku) |
| `hub/brand/README.md` | md | konform (menschliche Doku) |
| `hub/LOCK_ORDER/LOCK_ORDER.md` | md | konform (diese Datei) |
| `hub/git-befehlspal.txt` | txt | konform (Inhalt für Roman); enthält Markdown-Syntax, toleriert |
| `cmd/lab-studio/web/static/fonts/manrope-OFL.txt` | txt | Ausnahme: Schrift-Lizenztext, Runtime-Bestand, bleibt unverändert |
| `hub/chatroom.jsonl`, `hub/chatroom.schema.json` | json | konform (KI-Kommunikation/-Schema) |
| `hub/user/prompt-vorlage.yaml` | yaml | konform (Prompt-Vorlage) |
| `hub/mitglieder.csv` | csv | außerhalb der Drei-Format-Ordnung; Bestand bleibt, kein Änderungsbedarf |
| `hub/LOCK_ORDER/pruef/VERLAUF.md`, `hub/LOCK_ORDER/tools/WERKZEUGE.md` | md | konform (Prüfhistorie, Werkzeug-Registry) |
| `hub/LOCK_ORDER/tools/smoke.sh` | sh | konform (Werkzeugbereich, read-only) |

Konvertieren oder Verschieben jeder Datei nur mit nachgewiesenem Bedarf und Freigabe.

## 2. Kontext- und Token-Grenzen

Trennung: Runtime-Code = `cmd/lab-studio/`. Projektwissen = `hub/` und Root-Dokumente.

Einstiegspunkte für Modelle (in dieser Reihenfolge, nur was der Auftrag braucht):

1. Konkreter Auftrag des Menschen.
2. `hub/LOCK_ORDER/LOCK_ORDER.md` — bei jeder Struktur-, Format- oder Regel-Frage.
3. `branch-termux-pull-manifest.md` — bei Fragen zu Branch-Zweck, Herkunft, Arbeitsregeln.
4. `README.md` — nur bei Projekt-/Pfadfragen.

Lesebereiche (begrenzt, auftragsbezogen):

- `cmd/lab-studio/`: nur bei konkretem Codeauftrag; gezielt die nötigen Dateien, kein Vollladen des Verzeichnisses.
- `hub/`: nur die im Auftrag genannten Dateien. `hub/brand/originale/` niemals laden und niemals verändern (Originaldateien).
- Root-Ebene: keine rekursive Gesamtsicht, keine Muster wie `**/*` über Root oder `hub/`.

Verbotene Praktiken:

- Unkontrolliertes Laden von Root oder `hub/` in den Kontext.
- Verbindliche Regeln ausschließlich im Kontext behalten, ohne sie als Datei zu sichern (diese Datei ist die Sicherung).
- Runtime-Dateien unter `cmd/` aus Regeln heraus ändern — braucht separaten Auftrag.

## 3. Smoke-Prüfungen (Beleg)

Jede Prüfung läuft aus der Projekt-Root. Ausgeführt schlägt behauptet: Ergebnis hier dokumentieren.

| # | Prüfung | Kommando |
|---|---|---|
| S1 | Regelwurzel existiert | `test -d hub/LOCK_ORDER && test -f hub/LOCK_ORDER/LOCK_ORDER.md` |
| S2 | Nur erlaubte md-Dateien außerhalb Runtime | `find . -name '*.md' -not -path './cmd/*' -not -path './.git/*'` → nur: `README.md`, `branch-termux-pull-manifest.md`, `hub/brand/README.md` sowie alles unter `hub/LOCK_ORDER/` |
| S3 | txt-Dateien inventarisiert | `find . -name '*.txt' -not -path './.git/*'` → nur: `hub/git-befehlspal.txt`, `cmd/.../manrope-OFL.txt` |
| S4 | Runtime unberührt | `git status --short cmd/` → leer |
| S5 | Kein Push, Remote unverändert | `git status -sb` zeigt keine ahead/behind-Marker gegenüber `origin/termux-pull` |

Ausführung: einzeln oder gebündelt via Werkzeug `tools/smoke.sh` (siehe `tools/WERKZEUGE.md`).
Textmodus für Menschen; `--json <PR-Kennung>` für maschinenlesbare Ergebnisse — Ausgabe nur auf stdout,
Umleitung und Ablage in `pruef/PR-*.json` erfolgt durch den Ausführenden, nur bei tatsächlichem Lauf.
Dauerhafte Ergebnisnachweise gehören nach jeder Ausführung zusätzlich in `pruef/VERLAUF.md` —
hier bleibt nur der jeweils dokumentierte Regelstand.

## 4. Änderungsregeln für diese Subwurzel

- Struktur: `LOCK_ORDER.md` (Regeln), `pruef/` (nur ausgeführte Prüfungen, append-only), `tools/` (nur bei Bedarf angelegte read-only Skripte).
- Jede Regeländerung erzeugt einen neuen, eindeutig benannten Entwicklungsstand (`LOCK-JJJJ.MM.TT-n`) und wird in S1–S5 erneut nachgewiesen; der Nachweis wird in `pruef/VERLAUF.md` abgelegt.
- Skripte entstehen ausschließlich bei konkretem Steuerungsbedarf, mit Eintrag in `tools/WERKZEUGE.md`, ohne automatische Änderungen/Commits/Pushes.
- Keine Module, Router oder Konfigurationssysteme auf Vorrat.
- Bestehende Runtime und bestehende Inhalte bleiben unverändert bis zu einem separaten Auftrag.

## 5. Bestätigter Stand

- Subwurzel `/hub/LOCK_ORDER/` mit `LOCK_ORDER.md`, `pruef/` (Verlauf) und `tools/` (Registry + ein Smoke-Werkzeug).
- Formatordnung erweitert (`.sh` nur im Werkzeugbereich) und gegen den realen Baum geprüft; Befunde in §1.
- Kontextgrenzen unverändert in §2; Schreibrechte: Werkzeuge read-only, Prüfergebnisse nur in `pruef/`.
- Smoke-Prüfungen S1–S5 als Werkzeug ausgeführt (Text- und JSON-Modus); Ausführungsnachweise in `pruef/VERLAUF.md` und `pruef/PR-*.json` — Ergebnisdateien nur aus echten Läufen, inkl. geprüftem Fehlerpfad (Exit 1/2).
- Audit (Sprach-/Paketkompatibilität) durchgeführt, Befunde und Sicherungs-Hashes in `pruef/VERLAUF.md` (L2).
- Keine Commits, kein Push, keine Verschiebungen, keine Runtime-Änderungen.
- Wartet auf den nächsten ausdrücklichen Auftrag.
