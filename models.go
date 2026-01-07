package models

import (
	"time"

	"gorm.io/gorm"
)

// Room represents a hotel room
type Room struct {
	gorm.Model
	RoomNumber string  `gorm:"uniqueIndex;not null"`
	Type       string  `gorm:"not null"`
	Price      float64 `gorm:"not null"`
	IsBooked   bool    `gorm:"default:false"`
}

// Booking represents a reservation for a room
type Booking struct {
	gorm.Model
	RoomID    uint      `gorm:"not null"`
	GuestName string    `gorm:"not null"`
	CheckIn   time.Time `gorm:"not null"`
	CheckOut  time.Time `gorm:"not null"`
	Room      Room      `gorm:"foreignKey:RoomID"`
}
