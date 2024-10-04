package services

import (
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

// Función para enviar correos electrónicos
func SendEmail(to, subject, body string) error {
	// Obtener los valores de configuración desde las variables de entorno
	host := os.Getenv("MAILTRAP_HOST")
	portStr := os.Getenv("MAILTRAP_PORT")
	username := os.Getenv("MAILTRAP_USERNAME")
	password := os.Getenv("MAILTRAP_PASSWORD")
	from := os.Getenv("MAILTRAP_FROM")

	// Convertir el puerto a entero
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Error al convertir el puerto: %v", err)
		return err
	}

	// Configuración del correo electrónico
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	// Configuración de Mailtrap SMTP
	dialer := gomail.NewDialer(host, port, username, password)

	// Enviar correo electrónico
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("No se pudo enviar el correo a %s: %v", to, err)
		return err
	}

	log.Printf("Correo enviado a %s con éxito.", to)
	return nil
}
