package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/allyelvis/erp-pos-hotel-system-go/database"
	"github.com/allyelvis/erp-pos-hotel-system-go/handlers"
)

func main() {
	// Use structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	database.Init()

	r := gin.Default()

	// Configure CORS middleware from environment variables for security
	corsConfig := cors.DefaultConfig()
	// Example: ALLOWED_ORIGINS=http://localhost:3000,https://your-frontend.com
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		corsConfig.AllowOrigins = []string{origins}
	} else {
		// Default for local development
		corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Serve static files from the 'frontend' directory
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
		api.GET("/pos/menu", handlers.GetPosMenu)           // Kept for POS
		api.GET("/erp/inventory", handlers.GetErpInventory) // Kept for ERP_
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Goroutine for graceful shutdown
	go func() {
		slog.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("server exited gracefully")
}
