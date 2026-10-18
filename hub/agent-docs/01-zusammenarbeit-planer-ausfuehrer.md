# Zusammenarbeit Planer ↔ Ausführer

Für Arbeitsteilung zwischen einer **planenden KI ohne Repo-Zugriff** (z. B. Manus) und einem **ausführenden Werkzeug mit Dateisystem** (z. B. OpenCode, Claude Code). Eingeführt 2026-10-18 auf Romans Wunsch, um Kosten zu sparen: der Planer denkt, der Ausführer tippt und prüft.

Roman entscheidet über Ziel und Umfang. Ohne seine Freigabe kein Bauauftrag.

## 1. Was der Planer vor jedem Auftrag bekommen muss

Der Planer sieht das Repo nicht. Roman (oder der Ausführer) gibt ihm deshalb jedes Mal mit:

1. `hub/rules/CONTRACT.md` — welche Datei wofür zuständig ist
2. `hub/agent-docs/00-orientierung.md` — Regeln und aktueller Stand
3. den Plan aus `hub/rules/`, um den es geht
4. den **aktuellen** Inhalt jeder Datei, die geändert werden soll

Fehlt eine Datei, fragt der Planer danach, statt sie aus dem Gedächtnis zu rekonstruieren.

## 2. Feste technische Grenzen (gelten für jeden Auftrag)

- Go (Version laut `go.mod`), ein Paket `main` in `cmd/lab-studio/`, Dateien nach Themen getrennt (`CONTRACT.md`)
- SQLite über `modernc.org/sqlite`; SQL **nur** in `store.go`
- HTML über Go-Templates (`templates.go`) bzw. `web/*.html`; CSS/JS als Dateien in `web/static/`, per `go:embed` eingebaut
- **Kein** JavaScript-Framework, **kein** Build-Schritt, **kein** CSS-Framework, **keine** neue Go-Abhängigkeit ohne Romans Freigabe
- Schreibzugriffe: Sitzung + CSRF; jede Abfrage auf Nutzerdaten an das angemeldete Konto gebunden
- Lesbarkeit: kurze Funktionen, eine Prüfung pro Zeile, nichts doppelt kopieren; Texte für Nutzer auf Deutsch mit echten Umlauten
- Farben/Schrift nur aus `web/static/brand.css` (`hub/rules/12-brand-spec.md`)

## 3. Format: Bauauftrag (Planer → Ausführer)

```
BAUAUFTRAG <Nr> — <kurzer Titel>
Bezug: <hub/rules/NN-....md oder "Roman-Entscheidung vom JJJJ-MM-TT">

ZIEL: <ein Satz>
MINDESTKRITERIEN:
- <prüfbar formuliert>
VERBOTE:
- <was ausdrücklich nicht passieren darf>

DATEIEN:
[NEU] <pfad>
<vollständiger Inhalt>

[ÄNDERN] <pfad>
ALT (exakt so in der Datei, genau einmal):
<text>
NEU:
<text>

TESTS (in <datei>_test.go):
- <Testname>: <was er prüft>  — mindestens ein Test mit UNGÜLTIGER Eingabe
GEGENPROBE: <welche Zeile im Code man entfernt, damit der Test ROT wird>

PRÜFBEFEHL:
gofmt -l cmd/ ; go vet ./... ; go test ./...

BAUSTEIN: <name> | AUSLÖSER: <a/b/c oder Roman-Entscheidung> | BELEG: <Datei:Zeile> | INVARIANTE: <Prädikat>
(nur bei neuer Datei, neuer Tabelle oder neuer Abstraktion)
```

Regeln für den Planer:
- Änderungen an bestehenden Dateien **nur** als ALT/NEU-Paar, nie als „ersetze die ganze Datei", außer die Datei ist kurz und vollständig mitgeliefert.
- Ein Auftrag = ein Schritt, der für sich grün werden kann. Lieber drei kleine als ein großer.
- Nichts behaupten, was nicht im mitgelieferten Code steht.

## 4. Format: Rückmeldung (Ausführer → Planer)

```
RÜCKMELDUNG <Nr>
STATUS: GRÜN | ROT | ABWEICHUNG
AUSGEFÜHRT: <welche Dateien angelegt/geändert>
PRÜFAUSGABE (wörtlich, letzte Zeilen):
<Ausgabe von gofmt / go vet / go test>
GEGENPROBE: <durchgeführt? Test wurde rot? ja/nein>
ABWEICHUNGEN VOM AUFTRAG: <was anders gemacht wurde und warum, sonst "keine">
FRAGEN: <sonst "keine">
```

Regeln für den Ausführer:
- Erst umsetzen, dann **selbst** den Prüfbefehl ausführen. „Fertig" heißt GRÜN mit wörtlicher Ausgabe, nicht behauptet.
- Passt ein ALT-Text nicht zur Datei: nicht raten, STATUS ABWEICHUNG und den aktuellen Dateiausschnitt zurückmelden.
- Kein `git commit`/`git push` ohne Romans Freigabe.
- Optional: Auftrag und Rückmeldung als Zeile in `hub/chatroom.jsonl` festhalten (`chatroom.schema.json`, `type`: `task` bzw. `status`).
