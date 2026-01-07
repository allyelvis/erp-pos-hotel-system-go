package models

import (
	"time"

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
	Name     string `json:"name"`
	SKU      string `gorm:"uniqueIndex" json:"sku"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
	Reorder  int    `json:"reorder_level"`
}
