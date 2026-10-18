package main

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

func newMailerFromEnv() *Mailer {
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
	link := m.baseURL + "/verify?token=" + url.QueryEscape(token)
	if m.host == "" {
		log.Printf("[mail] kein SMTP_HOST gesetzt - Bestaetigungslink fuer %s: %s", email, link)
		return nil
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: RAP-LAB: E-Mail-Adresse bestaetigen\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n\r\n"+
		"Hallo,\r\n\r\nbitte bestaetige deine E-Mail-Adresse fuer RAP-LAB:\r\n\r\n%s\r\n\r\n"+
		"Der Link ist 24 Stunden gueltig. Wenn du dich nicht registriert hast, ignoriere diese Mail.\r\n",
		m.from, email, link)
	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.pass, m.host)
	}
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{email}, []byte(msg))
}
