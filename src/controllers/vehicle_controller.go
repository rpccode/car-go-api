package controllers

import (
	"go-auth-api/src/config"
	"go-auth-api/src/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Actualizar la ubicación del vehículo
func UpdateVehicleLocation(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Actualizar ubicación
	if err := vehicle.UpdateLocation(config.DB, vehicle.Latitude, vehicle.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la ubicación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ubicación actualizada correctamente"})
}

// Cambiar estado del vehículo
func UpdateVehicleStatus(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Actualizar estado
	if err := vehicle.UpdateStatus(config.DB, vehicle.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el estado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado del vehículo actualizado"})
}
func GetAvailableVehicles(c *gin.Context) {
	// Parse the start_time and end_time from the request query parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	// Convert the string parameters to time.Time
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha de inicio inválida", "err": err.Error()})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha de finalización inválida", "err": err, "endTime": endTime})
		return
	}

	// Lógica para obtener los vehículos disponibles
	vehicles, err := models.GetAllAvailableVehicles(config.DB, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los vehículos disponibles", "err": err})
		return
	}

	// Retornar la lista de vehículos disponibles
	c.JSON(http.StatusOK, vehicles)
}

// ListVehicles retrieves all vehicles
func ListVehicles(c *gin.Context) {
	vehicles, err := models.GetAllVehicles(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los vehículos", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// GetVehicle retrieves a specific vehicle by ID
func GetVehicle(c *gin.Context) {
	id := c.Param("id")
	var vehicle models.Vehicle
	vehicleID, err := strconv.Atoi(id) // Convert string ID to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido", "err": err.Error()})
		return
	}

	if err = vehicle.GetByID(config.DB, vehicleID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado", "err": err})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}
