package controllers

import (
	"go-auth-api/src/config"
	"go-auth-api/src/models"
	"net/http"
	"strconv"

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

// ListVehicles retrieves all vehicles
func ListVehicles(c *gin.Context) {
	vehicles, err := models.GetAllVehicles(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los vehículos"})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// GetVehicle retrieves a specific vehicle by ID
func GetVehicle(c *gin.Context) {
	id := c.Param("id")
	var vehicle models.Vehicle
	vehicleID, err := strconv.Atoi(id)
	if err = vehicle.GetByID(config.DB, vehicleID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehículo no encontrado"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}
