package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims estructura de las reclamaciones del token
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
