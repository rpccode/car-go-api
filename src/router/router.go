package routes

import (
	"go-auth-api/src/controllers"
	middlewares "go-auth-api/src/middleware"

	"github.com/gin-gonic/gin"
)

// Configurar las rutas de la API
func SetupRoutes(r *gin.Engine) {
	// Rutas públicas
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/token", controllers.GenerateToken)
	// Rutas protegidas con autenticación
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Ruta protegida de ejemplo
		protected.GET("/protected", func(c *gin.Context) {
			username := c.MustGet("username").(string)
			c.JSON(200, gin.H{"message": "Bienvenido " + username})
		})

		// Rutas de reservas
		protected.POST("/reservations", controllers.CreateReservation)
		protected.GET("/reservations/all", controllers.GetAllReservation)
		protected.GET("/reservations/:id", controllers.GetReservation)       // Obtener una reserva específica
		protected.PUT("/reservations/:id", controllers.UpdateReservation)    // Actualizar una reserva
		protected.DELETE("/reservations/:id", controllers.DeleteReservation) // Eliminar una reserva
		protected.POST("/reservations/check-availability", controllers.CheckVehicleAvailability)

		// Rutas de notificaciones
		protected.POST("/notifications", controllers.SendNotification)
		protected.POST("/notifications/reminder", controllers.SendNotificationReminder)
		protected.POST("/notifications/bot", controllers.SendBothNotifications)

		protected.GET("/notifications/:user_id", controllers.GetUserNotifications) // Obtener notificaciones del usuario

		// Rutas de vehículos
		protected.GET("/vehicles", controllers.ListVehicles)   // Listar vehículos disponibles
		protected.GET("/vehicles/:id", controllers.GetVehicle) // Obtener información de un vehículo específico
		protected.POST("/vehicles/check-availability", controllers.CheckVehicleAvailability)
		// Rutas de personajes
		protected.GET("/characters/fetch-all", controllers.FetchAndSaveAllCharacters) // Obtener y guardar todos los personajes
		protected.GET("/characters", controllers.GetPaginatedCharacters)              // Obtener personajes con paginación y búsqueda
	}
}
