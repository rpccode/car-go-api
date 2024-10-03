package models

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"` // Agrega este campo
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
}

// Registrar un nuevo usuario
func (u *User) Register(db *sql.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost) // Usa Password aquí
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, password, email, phone) VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(query, u.Username, string(hashedPassword), u.Email, u.Phone).Scan(&u.ID)
	return err
}

// Autenticar usuario
// Autenticar usuario
func (u *User) Authenticate(db *sql.DB, password string) error {
	// Cambiar el nombre del campo a 'password' en lugar de 'password_hash'
	query := `SELECT id, password FROM users WHERE username = $1 or email = $1`
	err := db.QueryRow(query, u.Username).Scan(&u.ID, &u.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("usuario no encontrado")
		}
		return err
	}

	// Imprimir para depurar
	log.Println("Usuario en base de datos:", u.Username)
	log.Println("Contraseña en base de datos (encriptada):", u.PasswordHash)
	log.Println("Contraseña recibida (sin encriptar):", password)

	// Comparar la contraseña ingresada con la contraseña encriptada
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return errors.New("la contraseña no coincide")
	}

	return nil
}
