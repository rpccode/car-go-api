package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func init() {
	// Cargar las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Función para enviar un correo electrónico
func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	// Obtener los valores desde las variables de entorno
	from := os.Getenv("MAILTRAP_FROM")
	host := os.Getenv("MAILTRAP_HOST")
	port := os.Getenv("MAILTRAP_PORT")
	username := os.Getenv("MAILTRAP_USERNAME")
	password := os.Getenv("MAILTRAP_PASSWORD")

	// Configurar correo electrónico
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configurar servidor SMTP
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}
	d := gomail.NewDialer(host, portNum, username, password)
	d.SSL = false
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Enviar correo
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("no se pudo enviar el correo: %v", err)
	}

	return nil
}
