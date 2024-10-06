package controllers

import (
	"go-auth-api/src/config"
	"go-auth-api/src/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Obtener clave JWT desde variables de entorno

// Registrar nuevo usuario
func Register(c *gin.Context) {
	var user models.User

	// Validar que los datos JSON son correctos
	if err := c.ShouldBindJSON(&user); err != nil || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos o contraseña vacía"})
		return
	}

	// Intentar registrar al usuario
	if err := user.Register(config.DB); err != nil {
		log.Printf("Error registrando usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Usuario registrado exitosamente"})
}

// Iniciar sesión
func Login(c *gin.Context) {
	var user models.User
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Validar que los datos JSON son correctos
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	user.Username = credentials.Username
	log.Println("Intentando autenticar usuario con password proporcionada")

	// Autenticar usuario
	if err := user.Authenticate(config.DB, credentials.Password); err != nil {
		log.Printf("Error de autenticación para usuario %s: %v", credentials.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Generar token JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generando token para usuario %s: %v", credentials.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	// Devolver token al cliente
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"id":      user.ID,
		"email":   credentials.Username,
		"token":   tokenString,
	})
}
