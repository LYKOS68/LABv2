# AUFTRAG — RAP-LAB (LABv2)

*Gesichert am Start der Refactoring-Phase. Dieses Dokument ist der verbindliche Rahmen für die kommenden Sessions.*

## 1. Projekt-Vision

Das Endprodukt ist eine **Marktplatz-App für Musikschaffende**: Community-Plattform mit Künstlerprofilen, Nachrichten, Feedback und einem Kleinanzeigenbereich für Angebote und Gesuche (Beats, Mixing, Mastering, Ghostwriting, Musikvideos, Covergestaltung, gemeinsame Projekte). Leitfrage: „Was kann ich selbst, und wen brauche ich für den Rest?"

Stack: Go als zentraler Orchestrator (Router, Businesslogik, SQLite), Frontend über HTML-Vorlagen, als PWA lauffähig. Repo: `LYKOS68/LABv2`, Lizenz Polyform Noncommercial 1.0.0, Owner Roman Hild (LYKOS68) als Product Owner & Architect.

## 2. Aktueller Stand (Inventur)

**Wo liegt der Code:** alles flach in `cmd/raplab/` (Package `main`), 5 Go-Dateien, 2415 Zeilen, dazu 42 grüne Tests.

| Datei | Zeilen | Status |
|---|---:|---|
| `auth.go` | 797 | **fertig**: Register, E-Mail-Verifikation, Login/Logout, Passwort-Reset, Account-Löschung, Sessions, Tokens, bcrypt, CSRF, Rate-Limiter, Mailer |
| `main_test.go` | 613 | 28 Tests, grün |
| `studio.go` | 501 | Studio-Page + API + CSRF, erste App-Fläche |
| `studio_test.go` | 350 | 14 Tests, grün |
| `main.go` | 146 | Basis-Setup: ServeMux, Store/SQLite, Static, Render, Home/Connect |
| `web/web.go` | 8 | toter Stub |

**Fertig:** Etappe 1 (DB, Router, Config) und Etappe 2 (Auth komplett inkl. Tests).
**Fehlt komplett:** Profile, Marktplatz-Inserate (CRUD), Chat/Nachrichten, Feedback/Ranglisten, Deployment.
**Strukturproblem:** alles in einem Package, `auth.go` als 800-Zeilen-Monolith — Aufteilung in `internal/`-Pakete steht an.

## 3. Die ZTI-II-Architektur

`internal/zti/` (4 Dateien, 466 Zeilen, 12 Tests grün) ist das **stille Fundament**: generische Chain mit Bits/Status/Station-Trace, konkrete Operationen IPC-Lookup und Object-Register, gebaut als Beweis des ZTIIKetten-Musters (Z → T → I → II → COMMIT / DROP).

Es wird für die Marktplatz-App **erstmal nicht aktiv genutzt** und bleibt unverändert als Werkzeugkasten im Repo. Kein Code der App darf von `internal/zti/` abhängig werden, bevor ein konkreter Bedarf aus den Etappen 1–7 nachweislich ist.

## 4. Zukunfts-Vision: Die Ad-Tech-Engine (Ricks Material)

**Konzept:** Eine eigene **Event-API und Consent-Engine** (Vorbild: TikTok Events-API), die später auf dem Marktplatz läuft — Nutzer-Events sauber erfassen, mit Einwilligung, dedupliziert und an Werbenetzwerke weiterreichen.

**Blaupause Ricks MyOS-Code für die spätere Go-Implementierung:**

- `myqueue.c` → **Event-Queues**: geblockte/coaleszbare Events, Prioritäten, Cap/Drop-Verhalten als Vorlage für die Event-Pipeline.
- `myobject.c` → **Advanced Matching & Dedupe**: Free-Stack-Registry, Kollisions-Veto, Handles als Vorlage für Event-/Objekt-Identifier und Deduplizierung.
- **ZTI-II als Consent-Gatekeeper**: jeder Event-Flow läuft als Kette — Consent (I) und Substrat (Z) sind Gatter; ohne alle vier Bits kein COMMIT, also kein Event-Export.

**Wichtig: Das ist ein Zukunfts-Schritt und beginnt erst, wenn die Basis (Etappe 1–7) läuft.** Keine Arbeit an der Ad-Tech-Engine vor Abschluss der Marktplatz-Kernfunktionen.

## 5. Die nächsten Schritte (Morgen)

1. `docs/PRODUCTION_COT.md` schreiben — den kompletten Produktionsprozess der App als Etappen 1–7 (Basis-Setup, Auth, Profile, Inserate-CRUD, Chat, Feedback/Ranglisten, Deployment).
2. `docs/REFACTOR_MATRIX.md` schreiben — pro Etappe: was ist schon da, was ist Schrott, was muss nur verschoben werden, was fehlt komplett.
3. Den `cmd/raplab`-Monolithen (besonders `auth.go`) in saubere `internal/`-Pakete aufteilen — Refactoring statt Neubau, Tests laufen grün halten.

**Vorgehensregel:** Erst CoT und Refactoring-Matrix, dann Code anfassen. Ricks Code nur dann heranziehen, wenn eine Etappe ein konkretes Problem hat (z. B. Chat → `myqueue.c`), nicht früher.
