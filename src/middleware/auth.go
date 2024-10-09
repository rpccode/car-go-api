package middlewares

import (
	"go-auth-api/src/models"
	"os"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Verificar si el encabezado Authorization está presente
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere autenticación"})
			c.Abort()
			return
		}

		// Dividir el encabezado por el espacio y comprobar si el formato es "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de autorización inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &models.Claims{}

		// Parsear el token JWT usando la clave secreta
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// Verificar si el token es válido
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Guardar el nombre de usuario en el contexto para que esté disponible en los controladores
		c.Set("username", claims.Username)
		c.Next()
	}
}
