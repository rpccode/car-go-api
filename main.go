package main

import (
	"go-auth-api/src/config"
	routes "go-auth-api/src/router"
	"log"

	// Importar el nuevo archivo de rutas

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Agrega los orígenes permitidos
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Cabeceras permitidas
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache de preflight
	}))

	// Configurar las rutas
	routes.SetupRoutes(r)

	// Ejecutar el servidor en el puerto 3000
	log.Println("Iniciando servidor HTTPS en http://192.168.0.104:3000")
	r.Run(":3000")
}
