package database

import (
	"log"
	"os"

	"github.com/user/erp-pos-hotel/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		DB, err = gorm.Open(sqlite.Open("erp_pos_hotel.db"), &gorm.Config{})
	}
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate all models
	err = DB.AutoMigrate(&models.User{}, &models.Room{}, &models.Booking{}, &models.MenuItem{}, &models.Order{}, &models.OrderItem{}, &models.InventoryItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedData()
}

func seedData() {
	var count int64
	DB.Model(&models.Room{}).Count(&count)
	if count == 0 {
		DB.Create(&models.Room{Number: "101", Type: "Single", Price: 100.0, Status: "Available"})
		DB.Create(&models.MenuItem{Name: "Burger", Category: "Main", Price: 12.50})
		DB.Create(&models.InventoryItem{Name: "Beef Patty", SKU: "BEEF-001", Quantity: 100, Unit: "units", Reorder: 20})
	}
}
