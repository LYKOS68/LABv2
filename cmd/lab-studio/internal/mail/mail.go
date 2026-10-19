package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os"
)

// Mailer verschickt Bestaetigungslinks per SMTP. Ohne SMTP_HOST wird der Link nur in die
// Konsole geschrieben, damit der Ablauf lokal ohne Mail-Dienst testbar bleibt.
type Mailer struct {
	host, port, user, pass, from string
	baseURL                      string
}

func NewFromEnv() *Mailer {
	m := &Mailer{
		host:    os.Getenv("SMTP_HOST"),
		port:    os.Getenv("SMTP_PORT"),
		user:    os.Getenv("SMTP_USER"),
		pass:    os.Getenv("SMTP_PASS"),
		from:    os.Getenv("SMTP_FROM"),
		baseURL: os.Getenv("BASE_URL"),
	}
	if m.port == "" {
		m.port = "587"
	}
	if m.from == "" {
		m.from = m.user
	}
	if m.baseURL == "" {
		m.baseURL = "http://127.0.0.1:8090"
	}
	return m
}

func (m *Mailer) SendVerificationLink(email, token string) error {
	return m.sendLink(email, "E-Mail-Adresse bestaetigen", "/verify",
		"bitte bestaetige deine E-Mail-Adresse fuer RAP-LAB:",
		"Der Link ist 24 Stunden gueltig. Wenn du dich nicht registriert hast, ignoriere diese Mail.", token)
}

func (m *Mailer) SendPasswordResetLink(email, token string) error {
	return m.sendLink(email, "Passwort zuruecksetzen", "/reset-password",
		"ueber diesen Link kannst du ein neues Passwort fuer RAP-LAB setzen:",
		"Der Link ist 1 Stunde gueltig. Wenn du das nicht angefordert hast, ignoriere diese Mail - dein Passwort bleibt unveraendert.", token)
}

func (m *Mailer) sendLink(email, subject, path, intro, outro, token string) error {
	link := m.baseURL + path + "?token=" + url.QueryEscape(token)
	if m.host == "" {
		log.Printf("[mail] kein SMTP_HOST gesetzt - %s fuer %s: %s", subject, email, link)
		return nil
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: RAP-LAB: %s\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n\r\n"+
		"Hallo,\r\n\r\n%s\r\n\r\n%s\r\n\r\n%s\r\n",
		m.from, email, subject, intro, link, outro)
	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.pass, m.host)
	}
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{email}, []byte(msg))
}
