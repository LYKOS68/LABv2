# Branch-Manifest: termux-pull

Verbindliche Grundlage der mobilen Entwicklungsarbeit.
Erstellt in der lokalen Projekt-Root. Stand: 2026-10-18.
Status: aktiv (Schritt 5/5 abgeschlossen — nur dieses Manifest wurde erstellt).

---

## 1. Auftrag und Identität

| Punkt | Festlegung | Beleg |
|---|---|---|
| Auftragsname | Repository `RAP-LABv2` | Auftrag 1/5 |
| Geprüfter Remote | `https://github.com/LYKOS68/LABv2.git` | `git remote -v` |
| Projektbezeichnung | RAP-LAB (Modul `rap-lab`) | `go.mod:1`, `README.md:3` |
| Branch | `termux-pull`, HEAD `4e87fa1`, synchron mit `origin/termux-pull` | `git branch -vv` |
| Lokaler Arbeitsbereich | Projekt-Root `/root/github/LABv2` | aktuelles Verzeichnis |
| Laufzeitkontext | Linux-Konsole, OpenCode CLI, Modell Big Pickle | Umgebungssitzung |

**LÜCKE:** Die exakte Entsprechung „RAP-LABv2" = GitHub-Repo `LYKOS68/LABv2` ist eine Auftragsangabe, kein verifizierter GitHub-Name. Namensdifferenz notiert, nicht aufgelöst.

## 2. Zweck der Branch

Die Branch `termux-pull` ist der mobile Entwicklungsarm des Windows-Projekts LabRap:
von unterwegs über das Handy (Termux, OpenCode CLI, Modell Big Pickle) werden
Verbesserungen hierhereingezogen und weiterentwickelt; das Windows-Projekt
bleibt die PC-Seite.

Beleg der Branch-Bestimmung im Projekt selbst: `README.md:1` —
„HIER LADE ICH VON UNTERWEGS ÜBERS HANDY PULL + VERBESSERUNG REIN UND KANN SIE SPÄTER VOM PC WIEDER LADEN".

## 3. Arbeitsregeln

1. ** zunächst ausschließlich lokal.** Alle Arbeit dieser Branch findet im lokalen Arbeitsbereich `/root/github/LABv2` statt.
2. **Remote-Änderungen und Pushes sind bis zur ausdrücklichen Freigabe verboten.** Kein `git push`, keine Änderung an `origin`, keine Umbenennungen, keine Remote-Deletes.
3. Bestehende lokale, uncommittete Zustände bleiben unberührt (siehe §4), solange kein Auftrag sie adressiert.
4. **Späterer PC-Workflow (nach Freigabe):** Änderungen dieser Branch am PC prüfen (`git log`, `git diff` gegen den PC-Stand) → einzeln übernehmen (Merge/Cherry-Pick) oder verwerfen → erst danach PC-seitig committen/pushen. Übernahme und Verwerfung sind getrennte, je zu beurteilende Entscheidungen.

## 4. Herkunft und Bearbeitungszustand

**Herkunftsregel:** Jeder Änderungsstand dieser Branch stammt aus der mobilen
Bearbeitung (Handy, Termux, OpenCode CLI, Big Pickle). Diese Herkunft ist
Kennung jedes Stands, nicht bloßer Kommentar.

Drei Zustände werden streng getrennt:

| Zustand | Wer | Was gilt |
|---|---|---|
| Mobile Bearbeitung | Modell auf dem Handy | Herkunft dieses Stands; Quelle aller lokal erzeugten Änderungen |
| PC-Prüfung | Modell/Mensch am Windows-PC | nur Lesen und Beurteilen, noch keine Übernahme |
| Spätere Übernahme | PC, nach Freigabe | erst dann wird mobiler Stand PC-seitiger Stand |

Regeln:

- **Keine unbelegten Annahmen** über Bearbeitungsort, Pfade oder Dateistrukturen: jede Aussage braucht einen Beleg (Git, Dateiinhalt, Aussage Romans) — andernfalls Kennzeichnung als **LÜCKE**.
- Ein PC-Modell darf mobile Änderungen **nicht** als PC-seitige Implementierung lesen oder behandeln. Herkunft „mobil" bleibt bis zur ausdrücklichen PC-Übernahme bestehen.
- Laufzeitkontext dieser Sitzung ist Linux (`/root/github/LABv2`); die physische Geräte-Herkunft „Termux auf dem Handy" ist Auftragsangabe und hier nicht unabhängig verifiziert → **LÜCKE**.

## 5. Tatsächlicher Root-Zustand (Ist-Aufnahme, 2026-10-18)

Aufnahmebasis für die spätere Root-Neuordnung. Nichts wurde verschoben oder gelöscht.

```
/root/github/LABv2/
  .claude/settings.json      Werkzeugkonfiguration (Windows-PC-Pfade, siehe Lücken)
  .git/                      Repository
  .github/workflows/test.yml CI: go test bei jedem Push (README.md:54)
  .gitignore                 ignoriert /data/*
  .vscode/tasks.json         VS-Code-Aufgabe „Server starten" (PC-Herkunft)
  README.md                  Projekt- und Pfaddoku, Branch-Zweck in Zeile 1
  cmd/lab-studio/            das Programm (Go, Web-Assets unter cmd/lab-studio/web/)
  go.mod / go.sum            Modul `rap-lab`, Go 1.27.0
  hub/                       Non-Code-Bereich: brand/, user/, chatroom.jsonl,
                             chatroom.schema.json, mitglieder.csv,
                             git-befehlspal.txt (untracked, 2026-10-18 23:40)
```

Git-Arbeitszustand (unverändert belassen):

- gelöscht, nicht committet: `test2.txt`
- untracked: `hub/git-befehlspal.txt`

Bestehende Abhängigkeiten:

- `data/` entsteht zur Laufzeit, ist git-ignoriert, nie auf GitHub (`README.md:60`, `.gitignore`).
- `cmd/lab-studio` läuft mit `go run ./cmd/lab-studio`, Prüfung `go test ./...` (`README.md:49`, `README.md:54`).
- `hub/` ist laut README „beschlossen, aber noch nicht umgesetzt" im Neubau; alte Regel-/Plandokus wurden 2026-10-18 entfernt, ältere Fassungen nur in der Git-Historie (`README.md:67-69`).
- Nicht vorhanden: `LOCK_ORDER.md`, `__hub-docs/` → die Plan-/Ledger-Infrastruktur des rap-lab-Über-Projekts existiert hier **nicht**.

### Lücken dieser Aufnahme

- `.claude/settings.json` verweist auf Windows-Pfade (`C:/Users/RPC/dev/rap-app`, `R:/dw`): PC-Herkunft, hier nicht verifizierbar → **LÜCKE**.
- Echte Verzeichnisstruktur der PC-Seite (`C:/Users/RPC/dev/rap-app`) ist unbekannt → **LÜCKE**.
- Inhalt/Stand von `hub/git-befehlspal.txt` und Ursache der `test2.txt`-Löschung nicht beauftragt, nicht bewertet → **LÜCKE**.

## 6. Verbindliche Pfadregeln

1. Root dieser Branch-Arbeit ist ausschließlich `/root/github/LABv2`.
2. **Programm-Bereich:** `cmd/` — nur Codeänderungen mit eigenem Auftrag.
3. **Knowledge-Bereich:** `hub/` — alles, was kein Programm ist (`README.md:61-65`). `hub/brand/` und `hub/user/` unverändert belassen.
4. **Laufzeit-Bereich:** `data/` — git-ignoriert, nie anfassen.
5. **Werkzeug-Bereich:** `.claude/`, `.vscode/`, `.github/` — Konfiguration, Herkunft dokumentiert, keine Änderung ohne Auftrag.
6. Root-Ebene: nur bestehende Dateien plus ausdrücklich freigegebene Manifeste. Keine neuen Verzeichnisse ohne Beleg und Freigabe.
7. Kein Verschieben, kein Löschen, kein Umbenennen ohne **belegten Bedarf und ausdrückliche Freigabe** (Auftrag 3/5).

## 7. Zielstruktur (Plan — nicht umgesetzt, nicht beschlossen)

Getrennte Bereiche für CLI-Arbeit und Knowledge, angestrebt:

| Bereich | Ort | Status |
|---|---|---|
| Programm | `cmd/` | Bestand |
| Knowledge / Non-Code | `hub/` | Bestand, Neubau beschlossen, offen (`README.md:67-69`) |
| CLI-Arbeit (Skripte, Werkzeuge) | eigener Bereich in der Root, Name offen | **Vorschlag, kein Beleg → LÜCKE** |
| Laufzeitdaten | `data/` | Bestand, ignoriert |
| Manifeste / Pläne | Root-Ebene, nur mit Freigabe | Bestand (diese Datei) |

Keine Umbauten: diese Tabelle ist Zielbild für den nächsten Bauauftrag, keine Handlungsanweisung für diese Session.

## 8. Backreferenzen und Abschluss

- Schritt 1/5 → Auftrag, Zweck, Regeln (§1–§3)
- Schritt 2/5 → Herkunft und Zustände (§4)
- Schritt 3/5 → Ist-Aufnahme und Pfadregeln (§5–§7)
- Schritt 4/5 → Rückkopplung: Repository, Branch, Arbeitsumgebung, Regeln, Pfad- und Herkunftsaussagen gegen den realen Zustand abgeglichen; Widersprüche korrigiert (Namensdifferenz §1, nicht existierende Plan-Infrastruktur §5, PC-Pfade als Lücken §5)
- Schritt 5/5 → nur dieses Manifest erstellt; lokale Zustände bewahrt; gewartet auf den nächsten ausdrücklichen Bauauftrag

**Zusammenfassung:** Manifest 1 | lokale Zustände unverändert | Remote/Push unberührt | Lücken 5 (offen) | Umbauten 0
