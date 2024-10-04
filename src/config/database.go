package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	// Cargar el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Obtener la cadena de conexión desde las variables de entorno
	connStr := os.Getenv("CONNECTIONSTRING")

	// Conectar a la base de datos
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar la conexión
	err = DB.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	}

	fmt.Println("Connected to database!")
}
