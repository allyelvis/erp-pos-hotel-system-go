package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Role     string `json:"role"`
	Password string `json:"-"`
}

type Room struct {
	gorm.Model
	Number string  `json:"number"`
	Type   string  `json:"type"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

type Booking struct {
	gorm.Model
	GuestName string    `json:"guest_name"`
	RoomID    uint      `json:"room_id"`
	Room      Room      `json:"room,omitempty"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
}

type MenuItem struct {
	gorm.Model
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type Order struct {
	gorm.Model
	TableNumber int         `json:"table_number"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Total       float64     `json:"total"`
	Status      string      `json:"status"`
}

type OrderItem struct {
	gorm.Model
	OrderID    uint     `json:"order_id"`
	MenuItemID uint     `json:"menu_item_id"`
	MenuItem   MenuItem `json:"menu_item,omitempty"`
	Quantity   int      `json:"quantity"`
}

type InventoryItem struct {
	gorm.Model
	Name     string  `json:"name"`
	SKU      string  `gorm:"uniqueIndex" json:"sku"`
	Quantity int     `json:"quantity"`
	Unit     string  `json:"unit"`
	Reorder  int     `json:"reorder_level"`
}

var db *gorm.DB

func initDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open("erp_pos_hotel.db"), &gorm.Config{})
	}
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&User{}, &Room{}, &Booking{}, &MenuItem{}, &Order{}, &OrderItem{}, &InventoryItem{})
	seedData()
}

func seedData() {
	var count int64
	db.Model(&Room{}).Count(&count)
	if count == 0 {
		db.Create(&Room{Number: "101", Type: "Single", Price: 100.0, Status: "Available"})
		db.Create(&MenuItem{Name: "Burger", Category: "Main", Price: 12.50})
		db.Create(&InventoryItem{Name: "Beef Patty", SKU: "BEEF-001", Quantity: 100, Unit: "units", Reorder: 20})
	}
}

func main() {
	initDB()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"status": "running"}) })
	api := r.Group("/api")
	{
		api.GET("/hotel/rooms", func(c *gin.Context) {
			var rooms []Room
			db.Find(&rooms)
			c.JSON(200, rooms)
		})
		api.GET("/pos/menu", func(c *gin.Context) {
			var menu []MenuItem
			db.Find(&menu)
			c.JSON(200, menu)
		})
		api.GET("/erp/inventory", func(c *gin.Context) {
			var items []InventoryItem
			db.Find(&items)
			c.JSON(200, items)
		})
	}
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	r.Run(":" + port)
}
