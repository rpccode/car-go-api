package main

import (
	"fmt"
	"go-auth-api/src/config"
	routes "go-auth-api/src/router"
	"log"
	"net/http"

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

	http.HandleFunc("/", homePage)

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

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Bienvenido a mi API</title>
    </head>
    <body>
        <h1>Bienvenido a mi API</h1>
        <p>Esta es la página pública de la API.</p>
    </body>
    </html>
    `)
}
