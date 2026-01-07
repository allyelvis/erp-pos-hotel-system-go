package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/erp-pos-hotel/database"
	"github.com/user/erp-pos-hotel/models"
	"gorm.io/gorm"
)

func GetRooms(c *gin.Context) {
	var rooms []models.Room
	if err := database.DB.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation
	if room.Number == "" || room.Type == "" || room.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: number, type, and a valid price are required."})
		return
	}

	if err := database.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		log.Printf("Error creating room: %v", err)
		return
	}
	c.JSON(http.StatusCreated, room)
}

func GetRoom(c *gin.Context) {
	var room models.Room
	id := c.Param("id")
	if err := database.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func UpdateRoom(c *gin.Context) {
	var room models.Room
	id := c.Param("id")
	if err := database.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&room)
	c.JSON(http.StatusOK, room)
}

func DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	// The .Delete method automatically handles the "not found" case gracefully
	// by doing nothing if the record doesn't exist.
	result := database.DB.Delete(&models.Room{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	// Preload Room data to include it in the response
	if err := database.DB.Preload("Room").Order("check_in desc").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- Validation ---
	if booking.GuestName == "" || booking.RoomID == 0 || booking.CheckIn.IsZero() || booking.CheckOut.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: guest_name, room_id, check_in, and check_out are required."})
		return
	}
	if !booking.CheckOut.After(booking.CheckIn) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Check-out date must be after check-in date."})
		return
	}

	// Check for overlapping bookings
	var existingBooking models.Booking
	err := database.DB.Where("room_id = ? AND check_in < ? AND check_out > ?", booking.RoomID, booking.CheckOut, booking.CheckIn).First(&existingBooking).Error
	if err == nil {
		// If err is nil, a booking was found, so there is an overlap.
		c.JSON(http.StatusConflict, gin.H{"error": "This room is already booked for the selected dates."})
		return
	} else if err != gorm.ErrRecordNotFound {
		// Handle potential database errors other than "not found"
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while checking for existing bookings."})
		return
	}

	// Find the room to calculate the total price
	var room models.Room
	if err := database.DB.First(&room, booking.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Calculate total price
	duration := booking.CheckOut.Sub(booking.CheckIn).Hours() / 24
	booking.Total = float64(duration) * room.Price
	booking.Status = "Confirmed" // Default status

	if err := database.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}
	c.JSON(http.StatusCreated, booking)
}

func GetPosMenu(c *gin.Context) {
	var menu []models.MenuItem
	if err := database.DB.Find(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve menu items"})
		return
	}
	c.JSON(http.StatusOK, menu)
}

func GetErpInventory(c *gin.Context) {
	var items []models.InventoryItem
	if err := database.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventory items"})
		return
	}
	c.JSON(http.StatusOK, items)
}
