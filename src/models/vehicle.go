package models

import (
	"database/sql"
	"time"
)

// Define custom types
type VehicleType string
type FuelType string
type VehicleStatus string
type Rating float64

// Vehicle struct represents the vehicle model
type Vehicle struct {
	ID              int           `json:"id"`
	Brand           string        `json:"brand"`
	Model           string        `json:"model"`
	LicensePlate    string        `json:"license_plate"`
	Latitude        float64       `json:"latitude"`
	Longitude       float64       `json:"longitude"`
	Type            VehicleType   `json:"type"`
	FuelType        FuelType      `json:"fuel_type"`
	Distance        float64       `json:"distance"`
	FuelEfficiency  float64       `json:"fuel_efficiency"`
	FuelConsumption float64       `json:"fuel_consumption"`
	PricePerMinute  float64       `json:"price_per_minute"`
	PricePerMile    float64       `json:"price_per_mile"`
	Status          VehicleStatus `json:"status"`
	ImageURL        string        `json:"image_url"`
	Rating          Rating        `json:"rating"`
	IsBooked        bool          `json:"is_booked"`
	IsReserved      bool          `json:"is_reserved"`
	IsAvailable     bool          `json:"is_available"` // corrected to 'is_available'
	IsRented        bool          `json:"is_rented"`
	IsFavorited     bool          `json:"is_favorited"`
	IsEconomic      bool          `json:"is_economic"`
	IsLuxury        bool          `json:"is_luxury"`
	IsOpen          bool          `json:"is_open"`
	IsClosed        bool          `json:"is_closed"`
	IsRental        bool          `json:"is_rental"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// UpdateLocation updates the vehicle's latitude and longitude
func (v *Vehicle) UpdateLocation(db *sql.DB, lat, long float64) error {
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
                     is_economic, is_luxury, is_open, is_closed, 
                     is_rental, created_at, updated_at 
              FROM vehicles`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate, &v.Latitude, &v.Longitude,
			&v.Type, &v.FuelType, &v.Distance, &v.FuelEfficiency, &v.FuelConsumption,
			&v.PricePerMinute, &v.PricePerMile, &v.Status, &v.ImageURL,
			&v.Rating, &v.IsBooked, &v.IsReserved, &v.IsAvailable,
			&v.IsRented, &v.IsFavorited, &v.IsEconomic, &v.IsLuxury,
			&v.IsOpen, &v.IsClosed, &v.IsRental, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

// GetByID retrieves a vehicle by its ID
func (v *Vehicle) GetByID(db *sql.DB, id int) error { // Changed id type to string
	query := `SELECT id, brand, model, license_plate, latitude, longitude, 
                     type, fuel_type, distance, fuel_efficiency, 
                     fuel_consumption, price_per_minute, price_per_mile, 
                     status, image_url, rating, is_booked, is_reserved, 
                     is_available, is_rented, is_favorited, 
                     is_economic, is_luxury, is_open, is_closed, 
                     is_rental, created_at, updated_at 
              FROM vehicles WHERE id = $1`
	return db.QueryRow(query, id).Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate,
		&v.Latitude, &v.Longitude, &v.Type, &v.FuelType,
		&v.Distance, &v.FuelEfficiency, &v.FuelConsumption,
		&v.PricePerMinute, &v.PricePerMile, &v.Status, &v.ImageURL,
		&v.Rating, &v.IsBooked, &v.IsReserved, &v.IsAvailable,
		&v.IsRented, &v.IsFavorited, &v.IsEconomic, &v.IsLuxury,
		&v.IsOpen, &v.IsClosed, &v.IsRental, &v.CreatedAt, &v.UpdatedAt)
}

// GetAllAvailableVehicles retrieves available vehicles within a date range
func GetAllAvailableVehicles(db *sql.DB, startTime, endTime time.Time) ([]Vehicle, error) {
	query := `
        SELECT id, brand, model, license_plate, latitude, longitude, 
               type, fuel_type, distance, fuel_efficiency, 
               fuel_consumption, price_per_minute, price_per_mile, 
               status, image_url, rating, is_booked, is_reserved, 
               is_available, is_rented, is_favorited, 
               is_economic, is_luxury, is_open, is_closed, 
               is_rental, created_at, updated_at 
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
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate, &v.Latitude, &v.Longitude,
			&v.Type, &v.FuelType, &v.Distance, &v.FuelEfficiency, &v.FuelConsumption,
			&v.PricePerMinute, &v.PricePerMile, &v.Status, &v.ImageURL,
			&v.Rating, &v.IsBooked, &v.IsReserved, &v.IsAvailable,
			&v.IsRented, &v.IsFavorited, &v.IsEconomic, &v.IsLuxury,
			&v.IsOpen, &v.IsClosed, &v.IsRental, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}
