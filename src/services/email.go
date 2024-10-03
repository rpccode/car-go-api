package services

import (
	"log"

	"gopkg.in/gomail.v2"
)

// Función para enviar correos electrónicos
func SendEmail(to, subject, body string) error {
	// Configuración del correo electrónico
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "no-reply@carsharing.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	// Configuración de Mailtrap SMTP
	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, "bf95a9a2879f61", "3a9153635ebabc")

	// Enviar correo electrónico
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("No se pudo enviar el correo a %s: %v", to, err)
		return err
	}

	log.Printf("Correo enviado a %s con éxito.", to)
	return nil
}
