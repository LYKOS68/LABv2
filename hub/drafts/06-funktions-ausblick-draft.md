# 06 — Funktions-Ausblick (Draft, nicht bindend)

**Status:** Reine Referenzliste für künftige Scoping-Entscheidungen. Nichts hiervon ist beauftragt, nichts davon ist gebaut, kein BAUSTEIN-Protokoll ausgelöst. §6.3 gilt unverändert: jede Kategorie wird erst gebaut, wenn ein konkreter Auftrag vorliegt — und erst nach der **zweiten** realen Umsetzung generalisiert, nicht vorher (§3 Kern-Definition: „wiederverwendbar" braucht zweite tatsächliche Nutzung, keine Theorie).

**Referenzmodell:** eBay Kleinanzeigen — als Beispiel für eine Plattform, die Nutzer, Inserate (Angebote) und Nachrichten zwischen Nutzern verbindet. Passt strukturell zu dem, was `LAB-CONNECT` andeutet (toter Nav-Link, bisher kein Scope).

---

## Generische Mechaniken (klein an Zahl, viele Oberflächen)

Jede konkrete Funktion unten gehört zu genau einer dieser Grundmechaniken:

| Mechanik | Kurzbeschreibung | Haben wir schon? |
|---|---|---|
| **Auth-Gate** | Session+CSRF+Rate-Limit vor jeder Mutation | ✅ gebaut (`handlers.go`, `csrf.go`, `ratelimit.go`) |
| **CRUD-Zyklus** | Erstellen/Lesen/Ändern/Löschen einer Entität | ❌ nur für `User` (teilweise — kein Update) |
| **Zustandsübergang** | Entität durchläuft definierte Zustände (State-Machine) | ❌ nicht vorhanden |
| **Liste+Filter+Sortierung** | Mehrere Entitäten durchsuchbar/sortierbar anzeigen | ❌ nicht vorhanden |
| **Mailbox/Messaging** | begrenzte Warteschlange pro Empfänger, Post+Get (s. MyOS-Vergleich) | ❌ nicht vorhanden |
| **Relations-Liste (n:m)** | User↔Entität-Beziehung (merken/folgen/blockieren) | ❌ nicht vorhanden |
| **Moderation/Reporting** | Melden + Prüfzyklus | ❌ nicht vorhanden |

---

## Konkrete Funktionen, nach Mechanik gruppiert

### Auth-Gate (Basis für alles Weitere — bereits fertig)
- Registrieren, Login, Logout, Account löschen ✅
- Passwort-Reset — fehlt, bewusst ausgeschlossen bisher
- E-Mail-Verifikation — fehlt, bewusst ausgeschlossen bisher
- Profil bearbeiten (Name/Avatar/Standort) — fehlt

### CRUD-Zyklus — Beispiel-Entität „Angebot"
- Angebot erstellen (Titel, Beschreibung, Preis, Kategorie, Bilder, Standort)
- Angebot bearbeiten
- Angebot löschen
- Eigene Angebote auflisten („Meine Anzeigen")

### Zustandsübergang — Beispiel-Entität „Angebot"
```
Entwurf → Veröffentlicht → (Pausiert) → Verkauft/Geschlossen → Gelöscht
```
- Veröffentlichen / Pausieren / Als verkauft markieren / Reaktivieren
- Automatisches Ablaufen nach X Tagen (optional, später)

### Liste+Filter+Sortierung
- Angebote durchsuchen (Volltext)
- Nach Kategorie/Preis/Standort filtern
- Nach Datum/Preis/Entfernung sortieren
- Einzelne Angebots-Detailseite
- Pagination

### Mailbox/Messaging (vgl. MyOS SendMessage/GetMessage-Vergleich, vorherige Session-Runde)
- Nachricht an Verkäufer senden
- Konversationsverlauf anzeigen
- Als gelesen/ungelesen markieren
- Nutzer blockieren
- generisches Muster: begrenzte Mailbox pro Empfänger (Kapazität/Quota gegen Spam)

### Relations-Liste (n:m)
- Angebot merken/zur Merkliste hinzufügen
- Von Merkliste entfernen
- Merkliste anzeigen
- (später denkbar: Nutzer folgen)

### Moderation/Reporting
- Angebot/Nutzer melden
- Melde-Warteschlange (falls Admin-Rolle später kommt)

---

## Ausdrücklich offen / nicht entschieden

- Ob LAB-CONNECT überhaupt ein Kleinanzeigen-artiges Marktplatz-Modell wird, ist nicht festgelegt — dieses Dokument nimmt es nur als Referenzrahmen zur Strukturierung, nicht als Produktentscheidung.
- Keine der obigen Mechaniken wird gebaut, ohne dass Roman zuerst eine konkrete erste Instanz beauftragt (z. B. „Angebot erstellen" als einzelner Schritt).
- Erst bei der **zweiten** konkreten Instanz einer Mechanik (z. B. wenn neben „Angebot erstellen" auch „Nachricht erstellen" gebaut wird) wird geprüft, ob sich wirklich eine gemeinsame Funktion lohnt (§6.3 Auslöser a) — nicht vorher.
