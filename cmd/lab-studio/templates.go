package main

import (
	"html/template"
	"log"
	"net/http"
)

// Farben kommen aus web/static/brand.css (gemeinsam mit LAB-STUDIO).
const baseCSS = `
body{background:var(--bg);color:var(--fg);font-family:var(--font);margin:0;padding:0;overflow-x:hidden}
header{display:flex;justify-content:space-between;align-items:center;padding:12px 24px;border-bottom:1px solid #1a1a1a}
header nav{display:flex;align-items:center;flex-wrap:wrap;justify-content:flex-end}
header a{color:var(--gold);text-decoration:none;letter-spacing:2px;margin-left:24px}
.brand{display:block;margin:0;line-height:0}
.brand img{height:48px;width:auto}
main{max-width:720px;margin:0 auto;padding:32px 24px}
h1{color:var(--gold);text-align:center}
.schriftzug{margin:16px 0 24px;line-height:0}
.schriftzug img{display:block;width:100%;height:auto}
.tagline{color:var(--gold);text-align:center;text-transform:uppercase;letter-spacing:1px;font-size:20px;line-height:1.5;margin:0 0 40px}
.points{margin-bottom:32px}
.points p{color:var(--fg);text-align:center;font-weight:bold;text-transform:uppercase;letter-spacing:2px;font-size:19px;margin:6px 0}
.fliesstext{margin:0 calc(50% - 50vw);padding:32px max(24px,calc(50vw - 336px));white-space:pre-line;font-size:18px;line-height:1.6;background:linear-gradient(rgba(0,0,0,.82),rgba(0,0,0,.82)),url(/static/img/leder-hintergrund.jpg) center/cover}
.ende{color:var(--gold);text-align:center;font-weight:bold;letter-spacing:3px;font-size:20px;margin-top:32px}
form{display:flex;flex-direction:column;gap:12px;max-width:360px;margin:32px auto}
input{padding:10px;background:var(--bg);border:1px solid var(--brown);border-radius:8px;color:var(--fg)}
button{padding:10px;background:transparent;border:1px solid var(--neon);border-radius:6px;color:var(--neon);font-weight:bold;cursor:pointer}
button:hover{background:var(--neon);color:var(--bg)}
.error{color:var(--error);text-align:center}
.ok{color:var(--neon);text-align:center}
.account-bar{text-align:center;margin-top:24px}
.switch{text-align:center;color:var(--gold);font-size:13px}
.switch a{color:var(--neon)}
.nav-btn{background:none;border:none;color:var(--gold);text-decoration:none;font-weight:normal;letter-spacing:2px;margin-left:24px;padding:0;cursor:pointer;font-size:inherit;font-family:inherit}
.nav-btn:hover{background:none;color:var(--neon)}
@media (max-width:600px){header{padding:10px 16px}header a,.nav-btn{margin-left:14px;letter-spacing:1px;font-size:14px}.brand img{height:40px}.tagline{font-size:17px}.points p{font-size:16px}.fliesstext{font-size:16px}}
`

var tpl = template.Must(template.New("base").Parse(`
{{define "layout"}}
<!doctype html>
<html lang="de">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>RAP-LAB</title><link rel="stylesheet" href="/static/brand.css"><style>` + baseCSS + `</style></head>
<body>
<header>
  <a class="brand" href="/"><img src="/static/img/krone.jpg" alt="RAP-LAB Startseite"></a>
  <nav>
    <a href="/studio">LAB-STUDIO</a>
    <a href="/connect">LAB-CONNECT</a>
    {{if .LoggedIn}}
      <form method="post" action="/logout" style="display:inline">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <button type="submit" class="nav-btn">Logout</button>
      </form>
    {{else}}
      <a href="/login">Anmelden</a>
    {{end}}
  </nav>
</header>
<main>{{template "content" .}}</main>
</body>
</html>
{{end}}

{{define "content-home"}}
<h1 class="schriftzug"><img src="/static/img/schriftzug.png" alt="RAP-LAB"></h1>
<p class="tagline">Ein Ort für das, was zwischen Idee und Musik passiert.</p>
<div class="points">
<p>Musik entwickeln.</p>
<p>Menschen finden.</p>
<p>Zusammenarbeiten.</p>
<p>Weiterkommen.</p>
</div>
<div class="fliesstext">Denn mal ehrlich:

Du kannst den nächsten Song auch wieder in deinen Handy-Notizen schreiben.
Den Producer kannst du wieder über fünf Ecken suchen.
Den Mixer wieder in irgendeiner Instagram-Story finden.
Den Videografen wieder irgendwo in einem Discord-Server suchen.
Und deine Brancheninformationen kannst du dir wieder aus zehn verschiedenen Quellen zusammensuchen.
Oder du hörst irgendwann auf, deinen Musikprozess aus Einzelteilen zusammenzubauen.</div>
<p class="ende">RAP-LAB.</p>
{{if .LoggedIn}}
<div class="account-bar">
  Angemeldet als {{.Email}} &middot;
  <form method="post" action="/delete-account" style="display:inline">
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
    <button type="submit" onclick="return confirm('Account wirklich endgueltig loeschen?')">Account loeschen</button>
  </form>
</div>
{{end}}
{{end}}

{{define "content-register"}}
<h1>Registrieren</h1>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/register">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="email" name="email" placeholder="E-Mail" value="{{.Email}}" required>
  <input type="password" name="password" placeholder="Passwort (min. 8 Zeichen)" required>
  <input type="password" name="confirm" placeholder="Passwort wiederholen" required>
  <button type="submit">Konto erstellen</button>
</form>
<p class="switch">Schon ein Konto? <a href="/login">Anmelden</a></p>
{{end}}

{{define "content-login"}}
<h1>Anmelden</h1>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/login">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="email" name="email" placeholder="E-Mail" value="{{.Email}}" required>
  <input type="password" name="password" placeholder="Passwort" required>
  <button type="submit">Anmelden</button>
</form>
{{if .ShowResend}}
<form method="post" action="/resend-verification">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="hidden" name="email" value="{{.Email}}">
  <button type="submit">Bestaetigungslink erneut senden</button>
</form>
{{end}}
<p class="switch"><a href="/forgot-password">Passwort vergessen?</a></p>
<p class="switch">Noch kein Konto? <a href="/register">Registrieren</a></p>
{{end}}

{{define "content-forgot"}}
<h1>Passwort vergessen</h1>
<p class="switch">Gib deine E-Mail-Adresse ein. Wir schicken dir einen Link, mit dem du ein neues Passwort setzen kannst.</p>
<form method="post" action="/forgot-password">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="email" name="email" placeholder="E-Mail" required>
  <button type="submit">Link anfordern</button>
</form>
<p class="switch"><a href="/login">Zurueck zur Anmeldung</a></p>
{{end}}

{{define "content-forgot-sent"}}
<h1>Passwort vergessen</h1>
<p class="ok">Falls ein Konto mit {{.Email}} existiert, haben wir dir einen Link geschickt.</p>
<p class="switch">Der Link ist 1 Stunde gueltig.</p>
<p class="switch"><a href="/login">Zur Anmeldung</a></p>
{{end}}

{{define "content-reset"}}
<h1>Neues Passwort</h1>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/reset-password">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="hidden" name="token" value="{{.Token}}">
  <input type="password" name="password" placeholder="Neues Passwort (min. 8 Zeichen)" required>
  <input type="password" name="confirm" placeholder="Passwort wiederholen" required>
  <button type="submit">Passwort speichern</button>
</form>
{{end}}

{{define "content-reset-result"}}
<h1>Neues Passwort</h1>
{{if .Error}}
<p class="error">{{.Error}}</p>
<p class="switch"><a href="/forgot-password">Neuen Link anfordern</a></p>
{{else}}
<p class="ok">Dein Passwort wurde geaendert.</p>
<p class="switch"><a href="/login">Jetzt anmelden</a></p>
{{end}}
{{end}}

{{define "content-verify-sent"}}
<h1>Fast geschafft</h1>
<p class="ok">Wir haben einen Bestaetigungslink an {{.Email}} geschickt.</p>
<p class="switch">Klicke auf den Link in der E-Mail, um dein Konto freizuschalten. Der Link ist 24 Stunden gueltig.</p>
<p class="switch"><a href="/login">Zur Anmeldung</a></p>
{{end}}

{{define "content-verify-result"}}
<h1>E-Mail-Bestaetigung</h1>
{{if .Error}}
<p class="error">{{.Error}}</p>
<p class="switch">Melde dich an, um einen neuen Link anzufordern. <a href="/login">Anmelden</a></p>
{{else}}
<p class="ok">Deine E-Mail-Adresse ist bestaetigt.</p>
<p class="switch"><a href="/login">Jetzt anmelden</a></p>
{{end}}
{{end}}
`))

type pageData struct {
	LoggedIn   bool
	Email      string
	Error      string
	CSRFToken  string
	ShowResend bool
	Token      string
}

func render(w http.ResponseWriter, r *http.Request, store *Store, contentName string, data pageData) {
	data.CSRFToken = csrfTokenFor(w, r, store)
	clone, err := tpl.Clone()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	clone = template.Must(clone.Parse(`{{define "content"}}{{template "` + contentName + `" .}}{{end}}`))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := clone.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("render error: %v", err)
	}
}
