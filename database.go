package database

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/allyelvis/erp-pos-hotel-system-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initializes the database connection.
// It attempts to connect to a PostgreSQL database using environment variables
// suitable for a Google Cloud SQL environment. If the required environment
// variables are not set, it falls back to an in-memory SQLite database
// for local development.
func Init() {
	var err error

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

	// If environment variables are set, connect to PostgreSQL (for Cloud Run)
	if dbUser != "" && dbPassword != "" && dbName != "" && instanceConnectionName != "" {
		slog.Info("Connecting to PostgreSQL database...")
		// DSN format for Cloud SQL with Unix socket
		dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s",
			dbUser,
			dbPassword,
			dbName,
			instanceConnectionName,
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			slog.Error("failed to connect to postgres database", "error", err)
			panic("failed to connect database")
		}
	} else {
		// Fallback to SQLite for local development
		slog.Info("PostgreSQL environment variables not set, falling back to local SQLite database.")
		DB, err = gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
		if err != nil {
			slog.Error("failed to connect to sqlite database", "error", err)
			panic("failed to connect database")
		}
	}

	slog.Info("Database connection established.")

	// Auto-migrate the schema
	DB.AutoMigrate(&models.Room{}, &models.Booking{})
	slog.Info("Database migrated.")
}
