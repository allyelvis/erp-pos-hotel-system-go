package handlers

import (
	"net/http"

	"github.com/allyelvis/erp-pos-hotel-system-go/database"
	"github.com/allyelvis/erp-pos-hotel-system-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRooms retrieves a list of all rooms.
func GetRooms(c *gin.Context) {
	var rooms []models.Room
	if err := database.DB.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// CreateRoom adds a new room to the database.
func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// GetRoom retrieves a single room by its ID.
func GetRoom(c *gin.Context) {
	var room models.Room
	id := c.Param("id")

	if err := database.DB.First(&room, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// UpdateRoom updates an existing room's details.
func UpdateRoom(c *gin.Context) {
	var room models.Room
	id := c.Param("id")

	if err := database.DB.First(&room, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room"})
		return
	}

	// Using a map for updates allows for partial updates (PATCH-like behavior)
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&room).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// DeleteRoom removes a room from the database.
func DeleteRoom(c *gin.Context) {
	id := c.Param("id")

	// The Delete function works even if the record is not found,
	// so we check the result's RowsAffected.
	result := database.DB.Delete(&models.Room{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

// --- Placeholder Handlers ---

// GetBookings is a placeholder for retrieving bookings.
func GetBookings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// CreateBooking is a placeholder for creating a booking.
func CreateBooking(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// GetPosMenu is a placeholder for retrieving the POS menu.
func GetPosMenu(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

// GetErpInventory is a placeholder for retrieving ERP inventory.
func GetErpInventory(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}
