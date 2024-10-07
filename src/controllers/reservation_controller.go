package controllers

import (
	"go-auth-api/src/config"
	"go-auth-api/src/models"
	"go-auth-api/src/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateReservation handles creating a new reservation
func CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	reservation.UserID = userID.(int)

	// Set start and end time for the reservation
	reservation.StartTime = time.Now()
	reservation.EndTime = time.Now().Add(2 * time.Hour)
	// Verificar la disponibilidad del vehículo
	isAvailable, err := reservation.IsVehicleAvailable(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar la disponibilidad del vehículo"})
		return
	}
	if !isAvailable {
		c.JSON(http.StatusConflict, gin.H{"error": "El vehículo no está disponible en el rango de tiempo solicitado"})
		return
	}
	// Create the reservation in the database
	if err := reservation.Create(config.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send reservation confirmation notification
	notification := models.Notification{
		UserID:  reservation.UserID,
		Message: "Su reserva ha sido confirmada.",
	}
	if err := notification.Send(config.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo enviar la notificación de confirmación"})
		return
	}

	// Canal para comunicar el resultado del envío del correo
	emailResult := make(chan error)

	// Fetch the user's email
	var userEmail string
	queryUsuario := `SELECT email FROM users WHERE id = $1`
	if err := config.DB.QueryRow(queryUsuario, reservation.UserID).Scan(&userEmail); err != nil {
		log.Printf("No se pudo obtener el correo del usuario: %v", err)
	} else {
		// Crear el cuerpo del correo
		emailBody := "Su reserva ha sido confirmada desde " +
			reservation.StartTime.Format("2006-01-02 15:04:05") + " hasta " +
			reservation.EndTime.Format("2006-01-02 15:04:05") + ". Detalles de la reserva:" +
			"<p>Vehículo: XYZ</p>"

		// Enviar el correo en segundo plano usando una goroutine
		go func() {
			subject := "Confirmación de reserva de vehículo"
			err := utils.SendEmail(userEmail, subject, emailBody)
			emailResult <- err // Enviar el resultado (nil si es exitoso, error si falla)
		}()
	}

	// Esperar el resultado del envío de correo por un tiempo limitado
	select {
	case err := <-emailResult:
		if err != nil {
			// Si hay un error en el envío del correo, lo informamos
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Reserva creada, pero falló el envío del correo de confirmación: " + err.Error()})
		} else {
			// Si el correo fue exitoso
			c.JSON(http.StatusOK, gin.H{"message": "Reserva creada exitosamente y correo enviado", "reservation": reservation})
		}
	case <-time.After(5 * time.Second): // Si el envío tarda más de 5 segundos, se considera un timeout
		log.Println("Advertencia: el envío del correo está tardando demasiado")
		c.JSON(http.StatusOK, gin.H{"message": "Reserva creada exitosamente, pero el correo está tardando más de lo esperado", "reservation": reservation})
	}
}

func GetReservation(c *gin.Context) {
	id := c.Param("id")
	reservationID, err := strconv.Atoi(id) // Convert id to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var reservation models.Reservation
	if err := reservation.GetByID(config.DB, reservationID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}
func GetAllReservation(c *gin.Context) {
	var reservation models.Reservation
	reservations, err := reservation.GetAll(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las reservas"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

// Controlador para verificar la disponibilidad de un vehículo
func CheckVehicleAvailability(c *gin.Context) {
	var input struct {
		VehicleID int       `json:"vehicle_id" binding:"required"`
		StartTime time.Time `json:"start_time" binding:"required"`
		EndTime   time.Time `json:"end_time" binding:"required"`
	}

	// Verificar si los datos proporcionados son válidos
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "err": err})
		return
	}

	// Crear una reserva temporal para hacer la verificación
	reservation := models.Reservation{
		VehicleID: input.VehicleID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
	}

	// Verificar la disponibilidad del vehículo
	isAvailable, err := reservation.IsVehicleAvailable(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar la disponibilidad del vehículo"})
		return
	}

	// Responder si está disponible o no
	if !isAvailable {
		c.JSON(http.StatusOK, gin.H{"available": false, "message": "El vehículo no está disponible en el rango de tiempo solicitado"})
	} else {
		c.JSON(http.StatusOK, gin.H{"available": true, "message": "El vehículo está disponible"})
	}
}

// UpdateReservation updates a specific reservation by ID
func UpdateReservation(c *gin.Context) {
	id := c.Param("id")
	reservationID, err := strconv.Atoi(id) // Convert id to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := reservation.Update(config.DB, reservationID); err != nil { // Pass reservationID as int
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la reserva"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reserva actualizada correctamente"})
}

// DeleteReservation deletes a specific reservation by ID
func DeleteReservation(c *gin.Context) {
	id := c.Param("id")
	reservationID, err := strconv.Atoi(id) // Convert id to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := models.DeleteReservation(config.DB, reservationID); err != nil { // Pass reservationID as int
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la reserva"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reserva eliminada correctamente"})
}
