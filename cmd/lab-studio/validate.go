package main

import (
	"regexp"
	"strings"
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func validateRegistration(email, password, confirm string) string {
	email = strings.TrimSpace(email)
	if email == "" || !emailRe.MatchString(email) {
		return "Bitte eine gueltige E-Mail-Adresse angeben."
	}
	if len(password) < 8 {
		return "Passwort muss mindestens 8 Zeichen haben."
	}
	if password != confirm {
		return "Passwoerter stimmen nicht ueberein."
	}
	return ""
}
