# hub/

Zweite Wurzel des Projekts — alles Nicht-Runtime: Regeln, Wissen, Skills, Agent-Koordination. Getrennt von `cmd/` (fertiger Code, erste Wurzel). Siehe `hub/rules/CONTRACT.md` für den vollständigen Pfadvertrag.

| Pfad | Zweck |
|---|---|
| `rules/` | eingefrorene Spezifikationen/Verträge — **verbindlich**, nicht ohne neuen Schritt verhandelbar |
| `chatroom.jsonl` | Append-only Nachrichtenlog zwischen Agents (z. B. Claude CLI, „agy") und Roman. Format: `chatroom.schema.json`. |
| `chatroom.schema.json` | JSON-Schema für jede Zeile in `chatroom.jsonl`. |
| `drafts/` | Unfertige Ideen, Brainstorming. **Kein Bauauftrag**, bis Roman es ausdrücklich freigibt. |
| `agent-docs/` | Knappe Betriebsregeln für jede KI, die hier arbeitet — siehe `agent-docs/00-orientierung.md` als Einstieg. |

**Nicht hierher gehört:** fertiger Code (→ `cmd/`), Laufzeitdaten (→ `data/`).
