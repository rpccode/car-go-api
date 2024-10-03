package models

import (
	"database/sql"
	"log"
)

type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	Gender  string `json:"gender"`
	Species string `json:"species"`
}

// Guardar el personaje en la base de datos
func (c *Character) SaveCharacter(db *sql.DB) error {
	query := `INSERT INTO characters (id, name, imageurl, status, gender, species)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  ON CONFLICT (id) DO NOTHING`

	_, err := db.Exec(query, c.ID, c.Name, c.Image, c.Status, c.Gender, c.Species)
	if err != nil {
		log.Printf("Error al guardar el personaje en la base de datos: %v\n", err)
		return err
	}
	return nil
}

// Obtener personajes con paginación
func GetCharacters(db *sql.DB, limit int, offset int, search string) ([]Character, error) {
	query := `SELECT id, name, imageurl, status, gender, species
			  FROM characters
			  WHERE LOWER(name) LIKE LOWER($1)
			  LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []Character
	for rows.Next() {
		var character Character
		if err := rows.Scan(&character.ID, &character.Name, &character.Image, &character.Status, &character.Gender, &character.Species); err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}
