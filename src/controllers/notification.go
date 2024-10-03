package controllers

import (
	"go-auth-api/src/config"
	"go-auth-api/src/models"
	"go-auth-api/src/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Función para enviar notificaciones
func sendNotification(c *gin.Context, notification models.Notification, subject string, body string) {
	// Registrar notificación en la base de datos
	if err := notification.Send(config.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la notificación: " + err.Error()})
		return
	}

	// Canal para comunicar el resultado del envío del correo
	emailResult := make(chan error)

	// Enviar correo en segundo plano usando goroutine
	go func() {
		recipientEmail := "recipient@example.com" // Obtén el correo del usuario desde el JSON o la base de datos
		err := utils.SendEmail(recipientEmail, subject, body)
		emailResult <- err // Enviar el resultado (nil si es exitoso, error si falla)
	}()

	// Esperar el resultado del envío de correo por un tiempo limitado
	select {
	case err := <-emailResult:
		if err != nil {
			// Si hay un error en el envío del correo, lo informamos
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar la notificación por correo: " + err.Error()})
		} else {
			// Si el correo fue exitoso
			c.JSON(http.StatusOK, gin.H{"message": "Notificación enviada exitosamente y correo enviado"})
		}
	case <-time.After(5 * time.Second): // Si el envío tarda más de 5 segundos, se considera un timeout
		log.Println("Advertencia: el envío del correo está tardando demasiado")
		c.JSON(http.StatusOK, gin.H{"message": "Notificación enviada exitosamente, pero el correo está tardando más de lo esperado"})
	}
}

// Enviar notificación
func SendNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	log.Println("Notificación recibida:", notification)

	// Crear el cuerpo del correo
	body := "<p>Tienes una nueva notificación:<br> " + notification.Message + "</p>"
	subject := "Nueva notificación"

	// Llamar a la función de envío de notificaciones
	sendNotification(c, notification, subject, body)
}

// Enviar recordatorio de notificación
func SendNotificationReminder(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	log.Println("Recordatorio de notificación recibido:", notification)

	// Crear el cuerpo del correo
	body := "<p>Recuerda devolver el vehículo a tiempo para evitar cargos adicionales.</p>"
	subject := "Recordatorio de Devolución"

	sendNotification(c, notification, subject, body)
}

// Enviar ambas notificaciones
func SendBothNotifications(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	log.Println("Notificación recibida:", notification)

	// Enviar la primera notificación
	body1 := "<p>Tienes una nueva notificación:<br> " + notification.Message + "</p>"
	subject1 := "Nueva notificación"
	sendNotification(c, notification, subject1, body1)

	// Enviar la segunda notificación (recordatorio)
	body2 := "<p>Recuerda devolver el vehículo a tiempo para evitar cargos adicionales.</p>"
	subject2 := "Recordatorio de Devolución"
	sendNotification(c, notification, subject2, body2)
}

// GetUserNotifications recupera las notificaciones de un usuario
func GetUserNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	notifications, err := models.GetNotificationsByUserID(config.DB, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las notificaciones"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
