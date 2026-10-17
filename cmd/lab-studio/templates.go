package main

import (
	"html/template"
	"log"
	"net/http"
)

const baseCSS = `
body{background:#0a0a0a;color:#f2f2f2;font-family:Arial,Helvetica,sans-serif;margin:0;padding:0}
header{display:flex;justify-content:space-between;align-items:center;padding:16px 24px;border-bottom:1px solid #222}
header a{color:#d4af37;text-decoration:none;font-weight:bold;letter-spacing:1px;margin-left:24px}
.brand{color:#d4af37;font-weight:bold;letter-spacing:1px}
main{max-width:640px;margin:0 auto;padding:32px 24px}
h1{color:#d4af37;text-align:center}
.tagline{color:#d4af37;text-align:center;font-weight:bold;text-transform:uppercase}
.points p{color:#fff;text-align:center;font-weight:bold;text-transform:uppercase;margin:4px 0}
.fliesstext{background:#141414;padding:24px;white-space:pre-line}
.ende{color:#d4af37;text-align:center;font-weight:bold;margin-top:24px}
form{display:flex;flex-direction:column;gap:12px;max-width:360px;margin:32px auto}
input{padding:10px;background:#111;border:1px solid #333;color:#fff}
button{padding:10px;background:#d4af37;border:none;color:#000;font-weight:bold;cursor:pointer}
.error{color:#ff6b6b;text-align:center}
.ok{color:#8fe388;text-align:center}
.account-bar{text-align:center;margin-top:24px}
.nav-btn{background:none;border:none;color:#d4af37;text-decoration:none;font-weight:bold;letter-spacing:1px;margin-left:24px;padding:0;cursor:pointer;font-size:inherit;font-family:inherit}
`

var tpl = template.Must(template.New("base").Parse(`
{{define "layout"}}
<!doctype html>
<html lang="de">
<head><meta charset="utf-8"><title>RAP-LAB</title><style>` + baseCSS + `</style></head>
<body>
<header>
  <span class="brand">RAP-LAB</span>
  <nav>
    <a href="#">LAB-STUDIO</a>
    <a href="#">LAB-CONNECT</a>
    {{if .LoggedIn}}
      <form method="post" action="/logout" style="display:inline">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <button type="submit" class="nav-btn">Logout</button>
      </form>
    {{else}}
      <a href="/register">Registrieren</a>
      <a href="/login">Login</a>
    {{end}}
  </nav>
</header>
<main>{{template "content" .}}</main>
</body>
</html>
{{end}}

{{define "content-home"}}
<h1>RAP-LAB</h1>
<p class="tagline">Ein Ort fuer das, was zwischen Idee und Musik passiert.</p>
<div class="points">
<p>Musik entwickeln.</p>
<p>Menschen finden.</p>
<p>Zusammenarbeiten.</p>
<p>Weiterkommen.</p>
</div>
<div class="fliesstext">Denn mal ehrlich:

Du kannst den naechsten Song auch wieder in deinen Handy-Notizen schreiben.
Den Producer kannst du wieder ueber fuenf Ecken suchen.
Den Mixer wieder in irgendeiner Instagram-Story finden.
Den Videografen wieder irgendwo in einem Discord-Server suchen.
Und deine Brancheninformationen kannst du dir wieder aus zehn verschiedenen Quellen zusammensuchen.
Oder du hoerst irgendwann auf, deinen Musikprozess aus Einzelteilen zusammenzubauen.</div>
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
{{end}}

{{define "content-login"}}
<h1>Login</h1>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/login">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="email" name="email" placeholder="E-Mail" value="{{.Email}}" required>
  <input type="password" name="password" placeholder="Passwort" required>
  <button type="submit">Anmelden</button>
</form>
{{end}}
`))

type pageData struct {
	LoggedIn  bool
	Email     string
	Error     string
	CSRFToken string
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
