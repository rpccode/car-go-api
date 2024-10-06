package main

import (
	"go-auth-api/src/config"
	routes "go-auth-api/src/router"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	config.ConnectDB()

	// Initialize the Gin router
	r := gin.Default()

	// Set up CORS middleware to allow cross-origin requests
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins for simplicity
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // HTTP methods allowed
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		AllowCredentials: true,                                                // Allow credentials like cookies or auth headers
		MaxAge:           12 * time.Hour,                                      // Cache the preflight response for 12 hours
	}))

	// Register the homepage route using Gin
	r.GET("/", homePage) // Use Gin's method to handle requests to the root path

	// Set up all the other routes from the router package
	routes.SetupRoutes(r)

	// Start the server on port 3000
	log.Println("Starting HTTPS server at http://192.168.0.104:3000")
	r.Run(":3000")
}

// homePage is a simple route handler for the root path that returns an HTML welcome message.
func homePage(c *gin.Context) {
	// Respond with an HTML page
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bienvenido a NuestraAPI</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;600&display=swap');
        
        body {
            font-family: 'Poppins', sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
        }
        .container {
            background-color: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            max-width: 800px;
            width: 90%;
        }
        h1 {
            color: #4a4a4a;
            text-align: center;
            margin-bottom: 10px;
            font-size: 2.5em;
            font-weight: 600;
        }
			 h4 {
            color: #4a4a4a;
            text-align: center;
            margin-bottom: 20px;
            font-size: 1.5em;
            font-weight: 400;
        }
        p {
            text-align: center;
            font-size: 1.1em;
            color: #666;
            margin-bottom: 30px;
        }
        .features {
            display: flex;
            justify-content: space-around;
            flex-wrap: wrap;
            margin-bottom: 40px;
        }
        .feature {
            flex-basis: 30%;
            text-align: center;
            margin-bottom: 20px;
        }
        .feature i {
            font-size: 2.5em;
            color: #764ba2;
            margin-bottom: 10px;
        }
        .feature h3 {
            font-size: 1.2em;
            color: #4a4a4a;
            margin-bottom: 10px;
        }
        .button {
            display: inline-block;
            background-color: #764ba2;
            color: #fff;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 30px;
            transition: all 0.3s ease;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .button:hover {
            background-color: #5a3d82;
            transform: translateY(-3px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        .center {
            text-align: center;
            margin-top: 30px;
        }
        .api-version {
            text-align: center;
            font-size: 0.9em;
            color: #888;
            margin-top: 20px;
        }
    </style>
    <script src="https://kit.fontawesome.com/your-fontawesome-kit.js" crossorigin="anonymous"></script>
</head>
<body>
    <div class="container">
        <h1>Bienvenido Golan Car API</h1>
        <h4>Creada por Rudy Alexander Perez casilla 1-15-1080</h4>


        <p>Descubre el poder de GOLAND CAR API API y lleva tus aplicaciones al siguiente nivel. Accede a datos en tiempo real, integra servicios avanzados y optimiza tu flujo de trabajo.</p>
        
        <div class="features">
            <div class="feature">
                <i class="fas fa-bolt"></i>
                <h3>Rápida</h3>
                <p>Respuestas en milisegundos para una experiencia fluida.</p>
            </div>
            <div class="feature">
                <i class="fas fa-lock"></i>
                <h3>Segura</h3>
                <p>Autenticación robusta y cifrado de datos.</p>
            </div>
            <div class="feature">
                <i class="fas fa-code"></i>
                <h3>Flexible</h3>
                <p>Fácil de integrar con múltiples lenguajes y frameworks.</p>
            </div>
        </div>
        
        <div class="center">
            <a href="#" class="button">Comenzar ahora</a>
        </div>
        
        <p class="api-version">Versión actual: v2.1.0</p>
    </div>
</body>
</html>
	`))
}
