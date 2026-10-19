# LOCK_ORDER — verbindliche Architekturregeln

Subwurzel `/hub/LOCK_ORDER/` (deklariert 2026-10-18, Auftrag 1/5).
Geltung: Projekt-Root `/root/github/LABv2`, Branch `termux-pull`.
Regeln liegen als Artefakte in dieser Datei, nicht im Kontextgedächtnis.

**Entwicklungsstand: `LOCK-2026.10.19-5`** — bestätigt durch die Smoke-Prüfungen (Ausführungsnachweise in `pruef/VERLAUF.md`, maschinenlesbar als `pruef/PR-*.json` aus echten Läufen).

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
| `hub/knowledge/brand/README.md` | md | konform (menschliche Doku) |
| `hub/LOCK_ORDER/LOCK_ORDER.md` | md | konform (diese Datei) |
| `hub/git-befehlspal.txt` | txt | konform (Inhalt für Roman); enthält Markdown-Syntax, toleriert |
| `cmd/lab-studio/web/static/fonts/manrope-OFL.txt` | txt | Ausnahme: Schrift-Lizenztext, Runtime-Bestand, bleibt unverändert |
| `hub/chatroom.jsonl`, `hub/chatroom.schema.json` | json | konform (KI-Kommunikation/-Schema) |
| `hub/user/prompt-vorlage.yaml` | yaml | konform (Prompt-Vorlage) |
| `hub/user/mitglieder.csv` | csv | außerhalb der Drei-Format-Ordnung; persönlicher Draft unter `user/`, kein Änderungsbedarf |
| `hub/LOCK_ORDER/pruef/VERLAUF.md`, `hub/LOCK_ORDER/tools/WERKZEUGE.md` | md | konform (Prüfhistorie, Werkzeug-Registry) |
| `hub/LOCK_ORDER/tools/smoke.sh` | sh | konform (Werkzeugbereich, read-only) |

Konvertieren oder Verschieben jeder Datei nur mit nachgewiesenem Bedarf und Freigabe.

## 2. Kontext- und Token-Grenzen

Trennung: Runtime-Code = `cmd/lab-studio/`. Projektwissen = `hub/` und Root-Dokumente.

### 2.1 `hub/` als zentrale Arbeitswurzel

`hub/` ist der exklusive Arbeits- und Wissensbereich: Romans Arbeitsmittel, CLI-Kontext,
Projektkontext und Knowledge — strikt getrennt von der Runtime unter `cmd/lab-studio/`.

- **Einstiegspunkte (jederzeit verfügbar):** die Dateien direkt in `hub/` — aktuell
  `chatroom.jsonl`, `chatroom.schema.json`, `git-befehlspal.txt` — sowie die in §2.2 genannten Einstiegsdateien.
- **Verboten:** das ungefragte Einlesen tieferer `hub/`-Verzeichnisse und ihrer Dateiinhalte.
  Zusätzliche Inhalte nur bei konkretem Aufgabenbezug oder ausdrücklicher Aufforderung.
- `hub/knowledge/brand/originale/` niemals laden und niemals verändern (Originaldateien).

### 2.2 Einstiegspunkte für Modelle (in dieser Reihenfolge, nur was der Auftrag braucht)

1. Konkreter Auftrag des Menschen.
2. `hub/LOCK_ORDER/LOCK_ORDER.md` — bei jeder Struktur-, Format- oder Regel-Frage.
3. `branch-termux-pull-manifest.md` — bei Fragen zu Branch-Zweck, Herkunft, Arbeitsregeln.
4. `README.md` — nur bei Projekt-/Pfadfragen.

### 2.3 Zielstruktur `hub/` (deklariert; Überführung erst nach Prüfung und Freigabe)

| Bereich | Zweck | Status |
|---|---|---|
| `hub/knowledge/` | Markenmaterial, Referenzen, Projektkontext, wiederkehrende Skills/Prompts; Unterordner nie automatisch laden | deklariert; enthält `brand/` |
| `hub/user/` | Persönliche Arbeitsmittel Romans | enthält `prompt-vorlage.yaml` und `mitglieder.csv` (Draft) |
| `hub/shared/` | Übergabedokumente, Prüf-/Testergebnisse, bestätigte Entwicklungsstände, Sicherungspunkte; maschinenlesbare Ergebnisse nur durch tatsächlich ausgeführte Prüfungen | deklariert, leer |

Trennungsprinzip: dauerhaftes Projektwissen (`knowledge`) ≠ persönliche Mittel (`user`) ≠
operative Übergabe- und Prüfergebnisse (`shared`). Keine doppelten Ablagen: eine Datei hat
genau einen Ort; bestehendes `hub/LOCK_ORDER/pruef/` bleibt vorerst, Überführung nach `shared/`
erst nach geprüfter Freigabe.

### 2.4 Lesebereiche (begrenzt, auftragsbezogen)

- `cmd/lab-studio/`: nur bei konkretem Codeauftrag; gezielt die nötigen Dateien, kein Vollladen des Verzeichnisses.
- `hub/`: nur die im Auftrag genannten Dateien oder die in §2.1 freigegebenen Einstiegspunkte.
- Root-Ebene: keine rekursive Gesamtsicht, keine Muster wie `**/*` über Root oder `hub/`.

Verbotene Praktiken:

- Unkontrolliertes Laden von Root oder `hub/` (einschließlich Unterordner) in den Kontext.
- Verbindliche Regeln ausschließlich im Kontext behalten, ohne sie als Datei zu sichern (diese Datei ist die Sicherung).
- Runtime-Dateien unter `cmd/` aus Regeln heraus ändern — braucht separaten Auftrag.
- Verschieben bestehender Inhalte in `knowledge/`, `user/` oder `shared/` ohne geprüfte Freigabe.

### 2.5 Zugriffs- und Herkunftsregel (verbindlich)

An diesem Repository arbeitet **ausschließlich Roman und seine eigenen CLI-Werkzeuge**
(Termux/OpenCode samt eingesetztem Modell sowie vom Menschen direkt ausgelöste lokale CLI-Läufe) —
**zu keinem Zeitpunkt** jemand oder etwas anderes.

- Keine anderen Personen, Modelle, Agenten, Bots oder Dienste lesen, schreiben, committen oder pushen.
- Kein automatischer Zugriff, keine Integrationen, keine Hintergrund- oder Fremdprozesse.
- Jede Ausnahme gilt nur, wenn Roman sie ausdrücklich und einzeln anweist.
- Herkunft jedes Stands ist damit eindeutig: Roman bzw. dessen CLI. Diese Herkunft ist nicht umdeutbar.

## 3. Smoke-Prüfungen (Beleg)

Jede Prüfung läuft aus der Projekt-Root. Ausgeführt schlägt behauptet: Ergebnis hier dokumentieren.

| # | Prüfung | Kommando |
|---|---|---|
| S1 | Regelwurzel existiert | `test -d hub/LOCK_ORDER && test -f hub/LOCK_ORDER/LOCK_ORDER.md` |
| S2 | Nur erlaubte md-Dateien außerhalb Runtime | `find . -name '*.md' -not -path './cmd/*' -not -path './.git/*'` → nur: `README.md`, `branch-termux-pull-manifest.md` sowie alles unter `hub/LOCK_ORDER/`, `hub/knowledge/`, `hub/user/`, `hub/shared/` |
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
- `hub/` als zentrale Arbeitswurzel deklariert (§2.1); Zielstruktur `knowledge/`, `user/`, `shared/` festgeschrieben (§2.3).
- Überführung ausgeführt: `hub/brand/` → `hub/knowledge/brand/`, `hub/mitglieder.csv` → `hub/user/mitglieder.csv` (Sicherung vorab, `git mv`, alle Pfadverweise angepasst).
- Formatordnung unverändert (§1); S2 an die neue Struktur angepasst (§3).
- Smoke-Prüfungen S1–S5 als Werkzeug nach der Regeländerung neu ausgeführt; Nachweis in `pruef/VERLAUF.md` und `pruef/PR-*.json`.
- Zugriffs- und Herkunftsregel festgeschrieben (§2.5): ausschließlich Roman und seine CLIs arbeiten an diesem Repository, zu keinem Zeitpunkt jemand anderes.
- Keine Commits, kein Push, keine Verschiebungen von Bestandsinhalten, keine Runtime-Änderungen.
- Wartet auf den nächsten ausdrücklichen Auftrag.
