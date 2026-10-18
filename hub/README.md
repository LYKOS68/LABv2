# hub/

Zweite Wurzel des Projekts — alles Nicht-Runtime: Regeln, Wissen, Skills, Agent-Koordination. Getrennt von `cmd/` (fertiger Code, erste Wurzel). Siehe `hub/rules/CONTRACT.md` für den vollständigen Pfadvertrag.

| Pfad | Zweck |
|---|---|
| `rules/` | Pläne, Spezifikationen, Verträge — nummeriert, **verbindlich**, werden nicht überschrieben, nur durch ein neues Dokument abgelöst. Übersicht: `rules/14-verzeichnisinventar.md` |
| `chatroom.jsonl` | Append-only Nachrichtenlog zwischen Agents (z. B. Claude CLI, „agy") und Roman. Format: `chatroom.schema.json`. |
| `chatroom.schema.json` | JSON-Schema für jede Zeile in `chatroom.jsonl`. |
| `drafts/` | Unfertige Ideen, Brainstorming. **Kein Bauauftrag**, bis Roman es ausdrücklich freigibt. |
| `brand/` | Marken-Originale (Logos, Schriftzüge) und Design-Vorlagen — nie verändern.  |
| `agent-docs/` | Knappe Betriebsregeln für jede KI, die hier arbeitet — siehe `agent-docs/00-orientierung.md` als Einstieg. |
| `user/` | Nur für Roman: Texte zum Kopieren ins Chatfenster. One-Shot-Prompts als YAML-Gruppen (`prompts-<gruppe>.yaml`), **append-only**. Schema und Regeln: `user/prompt-vorlage.yaml`. |

**Nicht hierher gehört:** fertiger Code (→ `cmd/`), Laufzeitdaten (→ `data/`).
