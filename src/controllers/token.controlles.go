package controllers

import (
	"go-auth-api/src/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SecretKey clave secreta para firmar el token (debería ir en las variables de entorno)
var SecretKey = []byte("secret_key")

// GenerateToken genera un token JWT sin fecha de expiración
func GenerateToken(c *gin.Context) {
	// Obtiene el username del body de la solicitud o cualquier otro dato necesario
	username := c.PostForm("username")

	// Crear las reclamaciones (claims) del token
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()), // Tiempo de emisión
		},
	}

	// Crear el token con las reclamaciones
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	// Enviar el token en la respuesta
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
