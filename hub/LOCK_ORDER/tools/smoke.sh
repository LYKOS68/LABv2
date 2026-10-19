#!/bin/sh
# smoke.sh — LOCK_ORDER Smoke-Prüfungen S1–S5
# Zweck:       Prüft die in hub/LOCK_ORDER/LOCK_ORDER.md §3 definierten Invarianten S1–S5.
# Eingaben:    optional "--json <PR-Kennung>" — Ausgabe ausschließlich als JSON auf stdout;
#              ohne Argumente menschliche Textausgabe. Andere Argumente sind ungültig (Exit 2).
#              Start in einem Verzeichnis innerhalb des Projekts; Root wird via git erkannt.
# Wirkung:     ausschließlich Lesen und Konsolenausgabe; schreibt und verändert nichts.
#              Ergebnisdateien entstehen nur durch Umleitung des Ausführenden, z. B.:
#              sh hub/LOCK_ORDER/tools/smoke.sh --json PR-...-L1 > hub/LOCK_ORDER/pruef/PR-...-L1.json
# Startbefehl: sh hub/LOCK_ORDER/tools/smoke.sh [--json PR-<YYYY-MM-DD>-smoke-<stand>-L<n>]
# Exit-Codes:  0 = alle Prüfungen bestanden; 1 = mindestens eine Prüfung fehlgeschlagen;
#              2 = ungültige Eingabe; 3 = keine Git-Root gefunden.
set -u
rc=0
mode=text
kennung=""
while [ $# -gt 0 ]; do
  case "$1" in
    --json)
      mode=json
      shift
      if [ $# -gt 0 ]; then kennung=$1; shift
      else echo "FAIL EINGABE: --json benötigt eine Prüfkennung PR-..."; exit 2; fi
      ;;
    *)
      echo "FAIL EINGABE: unbekanntes Argument: $1"; exit 2 ;;
  esac
done
if [ "$mode" = json ] && [ -z "$kennung" ]; then
  echo "FAIL EINGABE: leere Prüfkennung"; exit 2
fi

if ! root=$(git rev-parse --show-toplevel 2>/dev/null); then
  echo "FAIL ROOT: kein Git-Repository gefunden"; exit 3
fi
cd "$root" || exit 3

# S1 Regelwurzel existiert
S1st=PASS; S1dt="Regelwurzel vorhanden"
if ! { [ -d hub/LOCK_ORDER ] && [ -f hub/LOCK_ORDER/LOCK_ORDER.md ]; }; then
  S1st=FAIL; S1dt="Regelwurzel fehlt"; rc=1
fi

# S2 md ausserhalb Runtime nur an erlaubten Orten
bad_md=$(find . -name '*.md' -not -path './cmd/*' -not -path './.git/*' \
  | sed 's|^\./||' \
  | grep -vE '^(README\.md|branch-termux-pull-manifest\.md|hub/LOCK_ORDER/|hub/knowledge/|hub/user/|hub/shared/)' || true)
S2st=PASS; S2dt="md-Dateien nur an erlaubten Orten"
if [ -n "$bad_md" ]; then
  S2st=FAIL
  S2dt="unerlaubte md-Dateien: $(printf '%s ' $bad_md)"
  rc=1
fi

# S3 txt-Inventar exakt die bekannten beiden Dateien
found_txt=$(find . -name '*.txt' -not -path './.git/*' | sed 's|^\./||' | sort)
expected_txt=$(printf '%s\n' \
  'cmd/lab-studio/web/static/fonts/manrope-OFL.txt' \
  'hub/git-befehlspal.txt')
S3st=PASS; S3dt="txt-Inventar unveraendert"
if [ "$found_txt" != "$expected_txt" ]; then
  S3st=FAIL
  S3dt="txt-Inventar weicht ab: $(printf '%s ' $(printf '%s\n' "$found_txt" | tr '\n' ' '))"
  rc=1
fi

# S4 Runtime unveraendert
dirty_cmd=$(git status --short cmd/)
S4st=PASS; S4dt="cmd/ unveraendert"
if [ -n "$dirty_cmd" ]; then
  S4st=FAIL
  S4dt="cmd/ hat Aenderungen: $(printf '%s ' $(printf '%s\n' "$dirty_cmd" | tr '\n' ' '))"
  rc=1
fi

# S5 kein Push, Remote-Zustand unveraendert
sb_line=$(git status -sb | head -n 1)
S5st=PASS; S5dt="Branch synchron mit origin, kein Push ($sb_line)"
case "$sb_line" in
  *\[*)
    S5st=FAIL; S5dt="Branch weicht von origin ab ($sb_line)"; rc=1 ;;
esac

ts=$(date -u +%Y-%m-%dT%H:%M:%SZ)
stand=$(grep -o 'LOCK-[0-9][0-9.]*-[0-9]*' hub/LOCK_ORDER/LOCK_ORDER.md | head -n 1)

if [ "$mode" = json ]; then
  cat <<EOF
{
  "pruefkennung": "$kennung",
  "werkzeug": "hub/LOCK_ORDER/tools/smoke.sh",
  "stand": "$stand",
  "zeitpunkt_utc": "$ts",
  "exit_code": $rc,
  "ergebnisse": [
    {"id": "S1", "status": "$S1st", "detail": "$S1dt", "betroffen": ["hub/LOCK_ORDER/"]},
    {"id": "S2", "status": "$S2st", "detail": "$S2dt", "betroffen": ["hub/", "README.md", "branch-termux-pull-manifest.md"]},
    {"id": "S3", "status": "$S3st", "detail": "$S3dt", "betroffen": ["hub/git-befehlspal.txt", "cmd/lab-studio/web/static/fonts/manrope-OFL.txt"]},
    {"id": "S4", "status": "$S4st", "detail": "$S4dt", "betroffen": ["cmd/"]},
    {"id": "S5", "status": "$S5st", "detail": "$S5dt", "betroffen": [".git/"]}
  ]
}
EOF
else
  printf '%s %s %s\n' "$S1st" S1 "$S1dt"
  printf '%s %s %s\n' "$S2st" S2 "$S2dt"
  printf '%s %s %s\n' "$S3st" S3 "$S3dt"
  printf '%s %s %s\n' "$S4st" S4 "$S4dt"
  printf '%s %s %s\n' "$S5st" S5 "$S5dt"
fi
exit $rc
