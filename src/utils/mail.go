package utils

import (
	"crypto/tls"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

// Función para enviar un correo electrónico
func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	// Configurar correo electrónico
	m.SetHeader("From", "no-reply@carsharing.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configurar servidor SMTP
	d := gomail.NewDialer(
		"sandbox.smtp.mailtrap.io",
		25, // Puerto de Mailtrap
		"9145ccd8e3e3c7",
		"4c5d0e6a1e71ad",
	)
	d.SSL = false
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Enviar correo
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("no se pudo enviar el correo: %v", err)
	}

	return nil
}
