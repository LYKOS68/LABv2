# Generator: PROMPT_METHOD-Set

**Anweisung:** Diese Datei komplett als einen Prompt an die CLI geben. Sie erzeugt 10 Methoden-Prompts für die Software-Muster, die RAP-LAB braucht. Parameter vorher prüfen und anpassen.

---

Lies zuerst `prompts/kontext.md`. Halte dich an dessen Arbeitsregeln.

## Auftrag

Erzeuge ein Set aus 10 Methoden-Prompts. Jeder überführt ein Software-Muster für RAP-LAB in eine reproduzierbare Methode, die eine CLI später auf einen konkreten Schritt anwenden kann.

## Parameter

```
MUSTER:
  M-01  Route und Handler          net/http ServeMux mit Methoden-Mustern ("GET /")
  M-02  Seite rendern              html/template mit Layout und Seitenvorlage
  M-03  Statische Dateien          CSS und Bilder per embed ausliefern
  M-04  Formular verarbeiten       POST annehmen, prüfen, Fehler zurück ins Formular
  M-05  Konfiguration              PORT und DATABASE_URL aus Umgebungsvariablen
  M-06  Datenbank-Migration        nummerierte SQL-Dateien, beim Start anwenden
  M-07  Datenbank-Zugriff          eine Funktion pro Abfrage, Kontext und Fehler weiterreichen
  M-08  Liste mit Filter           Kategorie, Art (Suche/Angebot), Seiten
  M-09  Konto und Sitzung          Passwort-Hash, Sitzungs-Cookie, geschützte Routen
  M-10  Handler-Test               net/http/httptest, ein Test pro Muss-Kriterium
DB_TREIBER: [OFFEN]
ZIELORDNER: prompts/methoden/
```

Ist ein Parameter `[OFFEN]`, erzeuge die betroffenen Methoden trotzdem, setze aber in deren EINGABEN den Parameter als Pflicht und lass ABBRUCH mit `BLOCKIERT: <parameter>` greifen.

## Vorgehen

1. Prüfe, ob `prompts/kontext.md` existiert und lesbar ist. Sonst `BLOCKIERT: kontext.md`.
2. Lege pro Muster eine Datei an: `<ZIELORDNER><ID>-<kurzname>.md`, Kurzname klein mit Bindestrich.
3. Schreibe jede Datei exakt nach dem Ausgabeschema unten.
4. Prüfe jede Datei gegen die Prüfliste. Korrigiere, bis alle Punkte erfüllt sind, höchstens 3 Durchgänge.
5. Schreibe `<ZIELORDNER>index.md`: Tabelle mit ID, Name, Datei, braucht-Parameter.
6. Hänge an `prompts/status.md` an: `| GEN-M | <Status> | <Datum> | <Anzahl Dateien> |`.

## Ausgabeschema je Methode

```markdown
# <ID> <Name>

## ZIEL
Ein Satz: was nach Anwendung im Code funktioniert.

## EINGABEN
| Name | Typ | Pflicht | Beispiel aus RAP-LAB |
|---|---|---|---|

## SCHRITTE
Nummeriert, je Schritt eine Handlung an einer benannten Datei.
Höchstens 3 Dateien insgesamt.

## INVARIANTEN
Was nach jedem Schritt wahr bleiben muss (z. B. `go build ./...` fehlerfrei, keine neue Abhängigkeit ohne Parameter).

## OUTPUT_SCHEMA
Welche Dateien, Funktionen und Routen danach existieren, mit Signatur.

## TEST
Ausführbare Befehle mit erwarteter Ausgabe. Was nur im Browser sichtbar ist: `MANUELL PRÜFEN: …`.

## ABBRUCH
- Pflichteingabe fehlt → `BLOCKIERT: <name>`
- Nach 3 Versuchen Test rot → `DURCHGEFALLEN` + echte Fehlerausgabe
- Status-Zeile in `prompts/status.md`
```

## Prüfliste

- Jede Methode nennt ein Beispiel mit echten RAP-LAB-Daten aus `kontext.md` (Texte, Routen, Felder), keine Platzhalter wie „foo“.
- Keine Methode verlangt mehr als 3 Dateien.
- Jeder TEST ist ein ausführbarer Befehl oder als `MANUELL PRÜFEN` markiert.
- Keine Methode setzt Python, Frameworks oder Bibliotheken voraus, die nicht als Parameter genannt sind.
- Keine erfundenen Werte für Name, Farben oder Schrift.

## Ergebnis an mich

Nur: Status, Liste der angelegten Dateien, offene Parameter. Keine Erklärtexte.
