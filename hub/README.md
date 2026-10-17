# hub/

Arbeitsraum für Ideen, Drafts und Agent-Koordination — getrennt von `docs/` (eingefrorene Spezifikationen/Verträge) und von `cmd/` (fertiger Code).

| Pfad | Zweck |
|---|---|
| `chatroom.jsonl` | Append-only Nachrichtenlog zwischen Agents (z. B. Claude CLI, „agy") und Roman. Format: `chatroom.schema.json`. |
| `chatroom.schema.json` | JSON-Schema für jede Zeile in `chatroom.jsonl`. |
| `drafts/` | Unfertige Ideen, Brainstorming. **Kein Bauauftrag**, bis Roman es ausdrücklich freigibt. |
| `agent-docs/` | Knappe Betriebsregeln für jede KI, die hier arbeitet — siehe `agent-docs/00-orientierung.md` als Einstieg. |

**Nicht hierher gehört:** fertiger Code (→ `cmd/`), eingefrorene Specs (→ `docs/`), Laufzeitdaten (→ `data/`).
