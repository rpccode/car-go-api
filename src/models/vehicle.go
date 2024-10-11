package models

import (
	"database/sql"
	"time"
)

// Define custom types
type VehicleType *string
type FuelType *string
type VehicleStatus *string
type Rating float64

// Vehicle struct represents the vehicle model
type Vehicle struct {
	ID              int            `json:"id"`
	Brand           string         `json:"brand"`
	Model           string         `json:"model"`
	LicensePlate    string         `json:"license_plate"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	Type            VehicleType    `json:"type"`
	FuelType        FuelType       `json:"fuel_type"`
	Distance        float64        `json:"distance"`
	FuelEfficiency  float64        `json:"fuel_efficiency"`
	FuelConsumption float64        `json:"fuel_consumption"`
	PricePerMinute  float64        `json:"price_per_minute"`
	PricePerMile    float64        `json:"price_per_mile"`
	Status          VehicleStatus  `json:"status"`
	Rating          Rating         `json:"rating"`
	IsBooked        bool           `json:"is_booked"`
	IsReserved      bool           `json:"is_reserved"`
	IsAvailable     bool           `json:"is_available"`
	IsRented        bool           `json:"is_rented"`
	IsFavorited     bool           `json:"is_favorited"`
	IsEconomic      bool           `json:"is_economic"`
	IsLuxury        bool           `json:"is_luxury"`
	Images          []VehicleImage `json:"images"` // New field for multiple images
}

// VehicleImage struct represents an image of a vehicle
type VehicleImage struct {
	ID        int    `json:"id"`
	VehicleID int    `json:"vehicle_id"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

// GetVehicleImages retrieves all images for a specific vehicle
func GetVehicleImages(db *sql.DB, vehicleID int) ([]VehicleImage, error) {
	query := `SELECT id, vehicle_id, image_url, is_primary FROM vehicle_images WHERE vehicle_id = $1`
	rows, err := db.Query(query, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []VehicleImage
	for rows.Next() {
		var img VehicleImage
		if err := rows.Scan(&img.ID, &img.VehicleID, &img.ImageURL, &img.IsPrimary); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

// GetAllVehicles retrieves all vehicles from the database, including their images
func GetAllVehicles(db *sql.DB) ([]Vehicle, error) {
	query := `SELECT v.id, b.name, m.name, v.license_plate, v.latitude, v.longitude, 
                     ft.type, v.distance, v.fuel_efficiency, v.fuel_consumption, 
                     p.price_per_minute, p.price_per_mile, v.status, v.rating, 
                     v.is_booked, v.is_reserved, v.is_available, v.is_rented, 
                     v.is_favorited, v.is_economic, v.is_luxury
              FROM vehicles v
              JOIN brand b ON v.brand_id = b.id
              JOIN model m ON v.model_id = m.id
              JOIN fuel_type ft ON v.fuel_type_id = ft.id
              JOIN pricing p ON p.vehicle_id = v.id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		var fuelTypeStr, statusStr sql.NullString
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate, &v.Latitude, &v.Longitude,
			&fuelTypeStr, &v.Distance, &v.FuelEfficiency, &v.FuelConsumption,
			&v.PricePerMinute, &v.PricePerMile, &statusStr, &v.Rating,
			&v.IsBooked, &v.IsReserved, &v.IsAvailable, &v.IsRented,
			&v.IsFavorited, &v.IsEconomic, &v.IsLuxury,
		); err != nil {
			return nil, err
		}

		// Set nullable fields
		if fuelTypeStr.Valid {
			v.FuelType = &fuelTypeStr.String
		}
		if statusStr.Valid {
			v.Status = &statusStr.String
		}

		// Get vehicle images
		images, err := GetVehicleImages(db, v.ID)
		if err != nil {
			return nil, err
		}
		v.Images = images

		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

// GetByID retrieves a vehicle by its ID, including its images
func (v *Vehicle) GetByID(db *sql.DB, id int) error {
	query := `SELECT v.id, b.name, m.name, v.license_plate, v.latitude, v.longitude, 
                     ft.type, v.distance, v.fuel_efficiency, v.fuel_consumption, 
                     p.price_per_minute, p.price_per_mile, v.status, v.rating, 
                     v.is_booked, v.is_reserved, v.is_available, v.is_rented, 
                     v.is_favorited, v.is_economic, v.is_luxury
              FROM vehicles v
              JOIN brand b ON v.brand_id = b.id
              JOIN model m ON v.model_id = m.id
              JOIN fuel_type ft ON v.fuel_type_id = ft.id
              JOIN pricing p ON p.vehicle_id = v.id
              WHERE v.id = $1`
	err := db.QueryRow(query, id).Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate,
		&v.Latitude, &v.Longitude, &v.FuelType, &v.Distance,
		&v.FuelEfficiency, &v.FuelConsumption, &v.PricePerMinute,
		&v.PricePerMile, &v.Status, &v.Rating, &v.IsBooked,
		&v.IsReserved, &v.IsAvailable, &v.IsRented, &v.IsFavorited,
		&v.IsEconomic, &v.IsLuxury)
	if err != nil {
		return err
	}

	// Get vehicle images
	v.Images, err = GetVehicleImages(db, v.ID)
	return err
}

func (v *Vehicle) UpdateLocation(db *sql.DB, lat, long float64) error {
	query := `UPDATE vehicles SET latitude = $1, longitude = $2 WHERE id = $3`
	_, err := db.Exec(query, lat, long, v.ID)
	return err
}

func (v *Vehicle) UpdateStatus(db *sql.DB, status string) error {
	query := `UPDATE vehicles SET status = $1 WHERE id = $2`
	_, err := db.Exec(query, status, v.ID)
	return err
}
func GetAllAvailableVehicles(db *sql.DB, startTime, endTime time.Time) ([]Vehicle, error) {
	query := `
        SELECT id, brand, model, license_plate, latitude, longitude, status
        FROM vehicles v
        WHERE v.id NOT IN (
            SELECT vehicle_id 
            FROM reservations 
            WHERE (start_time <= $1 AND end_time >= $1) 
            OR (start_time <= $2 AND end_time >= $2)
            OR ($1 <= start_time AND $2 >= start_time)
        ) AND v.status = 'available'`

	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.LicensePlate, &v.Latitude, &v.Longitude, &v.Status); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}
