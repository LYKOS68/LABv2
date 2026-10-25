# Generator: PROMPT_CASCADE für eine Phase

**Anweisung:** Diese Datei komplett als einen Prompt an die CLI geben, vorher `PHASE` setzen. Sie erzeugt eine Kette fertiger Schritt-Prompts, die nacheinander an die CLI gehen.

---

Lies zuerst `prompts/kontext.md` und `prompts/methoden/index.md`. Halte dich an die Arbeitsregeln.

## Parameter

```
PHASE: H
ZIEL: Homepage mit Leitsatz, zwei Bereichen „Sag, was du wirklich brauchst.“ und
      „Zeig, was du kannst.“ mit den Beispielen aus kontext.md und Platzhalter „Gruppen – kommt bald“.
ZIELDATEI: prompts/kaskaden/<PHASE>.md
```

Fehlt `prompts/methoden/index.md` → `BLOCKIERT: methoden`. Ist eine vorige Phase in `prompts/status.md` nicht komplett BESTANDEN → `BLOCKIERT: phase <vorige>`.

## Vorgehen

Denke die Zerlegung zuerst Schritt für Schritt durch, bevor du schreibst:

1. **Zerlegen:** Welche Teilergebnisse braucht das ZIEL? Jedes muss mit einem Befehl oder `MANUELL PRÜFEN` nachweisbar sein.
2. **Chunks:** Schneide so, dass jeder Schritt eine Methode aus dem Index anwendet und höchstens 3 Dateien anfasst. Höchstens 10 Schritte.
3. **Abhängigkeiten:** Für jeden Schritt die Vorgänger-IDs. Keine Kreise.
4. **Übergaben typisieren:** Für jeden Schritt, was danach existiert: Datei, Funktion mit Signatur, Route.
5. **Teilprüfung:** Je Schritt die Muss-Kriterien und Testbefehle.
6. **Gesamtprüfung:** Ein letzter Schritt `<PHASE>-ABNAHME`, der das ZIEL als Ganzes prüft.

Schreibe nur das Ergebnis in die ZIELDATEI, nicht die Überlegung.

## Ausgabeschema der ZIELDATEI

````markdown
# Kaskade <PHASE>: <ZIEL in einem Satz>

| ID | Methode | hängt ab von | Übergabe |
|---|---|---|---|

## <PHASE>-01 <Titel>

```text
Lies prompts/kontext.md und prompts/methoden/<Methode>.md.
GATE: Prüfe in prompts/status.md, dass <Vorgänger-IDs> BESTANDEN sind. Sonst BLOCKIERT: <ID>.
ZIEL: <ein Satz>
DATEIEN: <höchstens 3 Pfade>
EINGABEN: <Werte für die EINGABEN der Methode, aus kontext.md>
ÜBERGABE: <was danach existiert, mit Signaturen>
MUSS:
- <Kriterium> → TEST: <Befehl> → ERWARTET: <Ausgabe>
ABBRUCH: nach Arbeitsregeln aus kontext.md, Status-Zeile für <ID> anhängen.
```
````

Jeder Schritt-Block ist ein vollständiger, eigenständiger Prompt zum Einfügen in die CLI.

## Prüfliste

- Jeder Schritt hat ein GATE auf seine Vorgänger. Kein Schritt ohne gültigen Vorgänger.
- Jede Übergabe nennt konkrete Dateien und Signaturen, keine Beschreibungen wie „die Logik“.
- Jede Methode im Schritt existiert im Index.
- Texte der Seite kommen wörtlich aus `kontext.md`. Offenes bleibt `[OFFEN: …]`.
- Die Kette deckt das ZIEL vollständig ab. Was nicht abgedeckt ist, steht am Ende unter `## Lücken`.

Höchstens 3 Durchgänge, bis die Prüfliste erfüllt ist, sonst `DURCHGEFALLEN`.
Hänge an `prompts/status.md` an: `| GEN-<PHASE> | <Status> | <Datum> | <Anzahl Schritte> |`.

## Ergebnis an mich

Nur: Status, Pfad der Zieldatei, Tabelle der Schritte, Lücken. Keine Erklärtexte.
