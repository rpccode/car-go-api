package utils

import (
	"go-auth-api/src/config"
	"go-auth-api/src/models"
)

// Función para enviar recordatorio de devolución del vehículo
func SendReturnReminder(userID int, userEmail string) error {
	message := "Recuerda devolver el vehículo a tiempo para evitar cargos adicionales."
	notification := models.Notification{
		UserID:  userID,
		Message: message,
	}

	// Guardar la notificación en la base de datos
	if err := notification.Send(config.DB); err != nil {
		return err
	}

	// Enviar recordatorio por correo electrónico
	emailBody := "<p>" + message + "</p>"
	return SendEmail(userEmail, "Recordatorio de Devolución", emailBody)
}
