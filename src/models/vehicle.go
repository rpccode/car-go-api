package models

import (
	"database/sql"
	"time"
)

// Define custom types for nullable values
type VehicleType sql.NullString
type FuelType sql.NullString
type VehicleStatus sql.NullString
type Rating sql.NullInt64 // Changed to sql.NullFloat64 to handle decimal/float rating

// Vehicle struct represents the vehicle model with all fields as nullable
type Vehicle struct {
	ID              sql.NullInt64   `json:"id"`
	Brand           sql.NullString  `json:"brand"`
	Model           sql.NullString  `json:"model"`
	LicensePlate    sql.NullString  `json:"license_plate"`
	Latitude        sql.NullFloat64 `json:"latitude"`
	Longitude       sql.NullFloat64 `json:"longitude"`
	Type            VehicleType     `json:"type"`
	FuelType        FuelType        `json:"fuel_type"`
	Distance        sql.NullFloat64 `json:"distance"`
	FuelEfficiency  sql.NullFloat64 `json:"fuel_efficiency"`
	FuelConsumption sql.NullFloat64 `json:"fuel_consumption"`
	PricePerMinute  sql.NullFloat64 `json:"price_per_minute"`
	PricePerMile    sql.NullFloat64 `json:"price_per_mile"`
	Status          VehicleStatus   `json:"status"`
	ImageURL        sql.NullString  `json:"image_url"`
	Rating          sql.NullInt64   `json:"rating"` // Updated to sql.NullFloat64 to handle float rating
	IsBooked        sql.NullBool    `json:"is_booked"`
	IsReserved      sql.NullBool    `json:"is_reserved"`
	IsAvailable     sql.NullBool    `json:"is_available"`
	IsRented        sql.NullBool    `json:"is_rented"`
	IsFavorited     sql.NullBool    `json:"is_favorited"`
	IsEconomic      sql.NullBool    `json:"is_economic"`
	IsLuxury        sql.NullBool    `json:"is_luxury"`
}

// UpdateLocation updates the vehicle's latitude and longitude
func (v *Vehicle) UpdateLocation(db *sql.DB, lat, long sql.NullFloat64) error {
	query := `UPDATE vehicles SET latitude = $1, longitude = $2 WHERE id = $3`
	_, err := db.Exec(query, lat, long, v.ID)
	return err
}

// UpdateStatus changes the status of the vehicle
func (v *Vehicle) UpdateStatus(db *sql.DB, status VehicleStatus) error {
	query := `UPDATE vehicles SET status = $1 WHERE id = $2`
	_, err := db.Exec(query, status, v.ID)
	return err
}

// GetAllVehicles retrieves all vehicles from the database
func GetAllVehicles(db *sql.DB) ([]Vehicle, error) {
	query := `SELECT id, brand, model, license_plate, latitude, longitude, 
                     type, fuel_type, distance, fuel_efficiency, 
                     fuel_consumption, price_per_minute, price_per_mile, 
                     status, image_url, rating, is_booked, is_reserved, 
                     is_available, is_rented, is_favorited, 
                     is_economic, is_luxury
              FROM vehicles`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		var typeStr, fuelTypeStr, statusStr sql.NullString
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate, &v.Latitude, &v.Longitude,
			&typeStr, &fuelTypeStr, &v.Distance, &v.FuelEfficiency, &v.FuelConsumption,
			&v.PricePerMinute, &v.PricePerMile, &statusStr, &v.ImageURL,
			&v.Rating, &v.IsBooked, &v.IsReserved, &v.IsAvailable,
			&v.IsRented, &v.IsFavorited, &v.IsEconomic, &v.IsLuxury,
		); err != nil {
			return nil, err
		}

		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

// GetByID retrieves a vehicle by its ID
func (v *Vehicle) GetByID(db *sql.DB, id int) error {
	query := `SELECT id, brand, model, license_plate, latitude, longitude, 
                     type, fuel_type, distance, fuel_efficiency, 
                     fuel_consumption, price_per_minute, price_per_mile, 
                     status, image_url, rating, is_booked, is_reserved, 
                     is_available, is_rented, is_favorited, 
                     is_economic, is_luxury 
              FROM vehicles WHERE id = $1`
	return db.QueryRow(query, id).Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate,
		&v.Latitude, &v.Longitude, &v.Type, &v.FuelType,
		&v.Distance, &v.FuelEfficiency, &v.FuelConsumption,
		&v.PricePerMinute, &v.PricePerMile, &v.Status, &v.ImageURL,
		&v.Rating, &v.IsBooked, &v.IsReserved, &v.IsAvailable,
		&v.IsRented, &v.IsFavorited, &v.IsEconomic, &v.IsLuxury,
	)
}

// GetAllAvailableVehicles retrieves available vehicles within a date range
func GetAllAvailableVehicles(db *sql.DB, startTime, endTime time.Time) ([]Vehicle, error) {
	query := `
        SELECT id, brand, model, license_plate, latitude, longitude, 
               type, fuel_type, distance, fuel_efficiency, 
               fuel_consumption, price_per_minute, price_per_mile, 
               status, image_url, rating, is_booked, is_reserved, 
               is_available, is_rented, is_favorited, 
               is_economic, is_luxury
        FROM vehicles v
        WHERE v.id NOT IN (
            SELECT vehicle_id 
            FROM reservations 
            WHERE status = 'activa' 
            AND ((start_time <= $1 AND end_time >= $1) 
            OR (start_time <= $2 AND end_time >= $2)
            OR ($1 <= start_time AND $2 >= start_time))
        )`

	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		var typeStr, fuelTypeStr, statusStr sql.NullString
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate, &v.Latitude, &v.Longitude,
			&typeStr, &fuelTypeStr, &v.Distance, &v.FuelEfficiency, &v.FuelConsumption,
			&v.PricePerMinute, &v.PricePerMile, &statusStr, &v.ImageURL,
			&v.Rating, &v.IsBooked, &v.IsReserved, &v.IsAvailable,
			&v.IsRented, &v.IsFavorited, &v.IsEconomic, &v.IsLuxury,
		); err != nil {
			return nil, err
		}

		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}
