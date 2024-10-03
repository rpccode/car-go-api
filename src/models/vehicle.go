package models

import "database/sql"

type Vehicle struct {
	ID           int     `json:"id"`
	LicensePlate string  `json:"license_plate"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	Status       string  `json:"status"` // disponible, reservado, en uso
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

// Actualizar la ubicación del vehículo
func (v *Vehicle) UpdateLocation(db *sql.DB, lat, long float64) error {
	query := `UPDATE vehicles SET latitude = $1, longitude = $2 WHERE id = $3`
	_, err := db.Exec(query, lat, long, v.ID)
	return err
}

// Cambiar el estado del vehículo
func (v *Vehicle) UpdateStatus(db *sql.DB, status string) error {
	query := `UPDATE vehicles SET status = $1 WHERE id = $2`
	_, err := db.Exec(query, status, v.ID)
	return err
}

// Obtener todos los vehículos
func GetAllVehicles(db *sql.DB) ([]Vehicle, error) {
	query := `SELECT id, license_plate, brand, model, status, latitude, longitude FROM vehicles`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(&v.ID, &v.LicensePlate, &v.Brand, &v.Model, &v.Status, &v.Latitude, &v.Longitude); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

// Obtener vehículo por ID
func (v *Vehicle) GetByID(db *sql.DB, id int) error {
	query := `SELECT license_plate, brand, model, status, latitude, longitude 
              FROM vehicles WHERE id = $1`
	return db.QueryRow(query, id).Scan(&v.LicensePlate, &v.Brand, &v.Model, &v.Status, &v.Latitude, &v.Longitude)
}
