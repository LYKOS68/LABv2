# 10 — Plan: Startseite nach Romans Vorlage

**Auslöser:** Roman-Entscheidung (Session 2026-10-18). Vorlage: Handy-Screenshot einer verlorenen Version (`C:\Users\RPC\dev\rap-app\my-space\__source-docs\bilder\homepage-bsp.jpeg`). Schließt die LÜCKE „Logo-/Signature-Assets fehlen" aus `01-spezifikation.md` §4.

**Status:** umgesetzt. `go test ./...` 33/33 PASS (neu: `TestHomeShowsBrandImages`). Sichtprüfung per Edge-Screenshot in 1280 px und in einem 412-px-Rahmen (Handybreite): Aufbau wie Vorlage, keine seitliche Verschiebung. Browser-Prüfung durch Roman offen.

| Punkt | Inhalt |
|---|---|
| Ziel | Startseite sieht aus wie die Vorlage |
| Mindestkriterien | oben links Kronen-Icon (`krone-logo-second.jpeg`), rechts LAB-STUDIO / LAB-CONNECT / Anmelden in Gold; Mitte Schriftzug-Bild (`rap-lab-schriftzug.png`) statt Text-Überschrift; Tagline Gold in Großbuchstaben; vier Punkte Weiß fett; Fließtext-Block über die ganze Breite mit Leder-Schriftzug (`unterschrift-ledercover2.jpeg`) blass im Hintergrund; „RAP-LAB." Gold am Ende |
| Qualitätskriterien | Bilder als verkleinerte Web-Kopien unter `cmd/lab-studio/web/static/img/` (Krone 192 px ≈ 12 KB, Schriftzug 1000 px, Leder 1280 px) — Originale bleiben unverändert; Bilder haben `alt`-Text; Seite funktioniert auch auf dem Handy |
| Prüfverfahren | Test: Startseite bindet die drei Bilder ein, alle drei werden ausgeliefert. Danach Sichtprüfung durch Roman im Browser (Handy + PC) |
| Verbote | keine neue Dependency, keine Schrift-Dateien von außen |

BAUSTEIN: Marken-Bilder (`web/static/img/`) | AUSLÖSER: Roman-Entscheidung außerhalb §6.3 | BELEG: `01-spezifikation.md` §4 LÜCKE Logo/Signature | INVARIANTE: `jedes eingebundene Bild wird mit Status 200 ausgeliefert`
