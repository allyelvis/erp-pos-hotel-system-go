package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/user/erp-pos-hotel/database"
	"github.com/user/erp-pos-hotel/handlers"
)

func main() {
	database.Init()

	r := gin.Default()

	// Use CORS middleware to allow requests from the frontend
	r.Use(cors.Default())

	// Serve static files (HTML, CSS, JS) from the 'frontend' directory
	r.Static("/ui", "./frontend") // Note: The UI will be at /ui/index.html

	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"status": "running"}) })

	api := r.Group("/api")
	{
		hotel := api.Group("/hotel")
		{
			hotel.GET("/rooms", handlers.GetRooms)
			hotel.POST("/rooms", handlers.CreateRoom)
			hotel.GET("/rooms/:id", handlers.GetRoom)
			hotel.PUT("/rooms/:id", handlers.UpdateRoom)
			hotel.DELETE("/rooms/:id", handlers.DeleteRoom)
			hotel.GET("/bookings", handlers.GetBookings)
			hotel.POST("/bookings", handlers.CreateBooking)
		}
		api.GET("/pos/menu", handlers.GetPosMenu)       // Kept for POS
		api.GET("/erp/inventory", handlers.GetErpInventory) // Kept for ERP
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
