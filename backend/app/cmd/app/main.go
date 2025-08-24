package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"                     // Added golang-migrate
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL driver for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source for migrate
	"github.com/joho/godotenv"                                 // Added godotenv

	"github.com/mkbagandov/kingsman/backend/app/internal/delivery"
	"github.com/mkbagandov/kingsman/backend/app/internal/infrastructure"
	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, loading from environment: %v", err)
	}

	// Load configuration (e.g., from environment variables or a config file)
	dbConfig := &infrastructure.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432, // Default PostgreSQL port
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Initialize database connection
	db, err := infrastructure.NewPostgreSQLDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Construct database URL for migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	// Run database migrations
	m, err := migrate.New(
		"file://db/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations applied successfully!")

	// Initialize repositories
	userRepo := infrastructure.NewPostgreSQLUserRepository(db)
	storeRepo := infrastructure.NewPostgreSQLStoreRepository(db)
	notificationRepo := infrastructure.NewPostgreSQLNotificationRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	storeUseCase := usecase.NewStoreUseCase(storeRepo)
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo)

	// Initialize handlers
	userHandler := delivery.NewUserHandler(userUseCase)
	storeHandler := delivery.NewStoreHandler(storeUseCase)
	notificationHandler := delivery.NewNotificationHandler(notificationUseCase)

	// Setup HTTP router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Kingsman Backend!"))
	})

	// User routes
	r.Post("/register", userHandler.RegisterUser)
	r.Get("/users/{userID}", userHandler.GetUserProfile)
	r.Get("/users/{userID}/discount-card", userHandler.GetUserDiscountCard)    // New route
	r.Put("/users/{userID}/discount-card", userHandler.UpdateUserDiscountCard) // New route
	r.Get("/users/{userID}/qrcode", userHandler.GetUserQRCode)                 // New route

	// Store routes
	r.Get("/stores", storeHandler.GetStores)
	r.Get("/stores/{storeID}", storeHandler.GetStoreByID)

	// Notification routes
	r.Post("/notifications", notificationHandler.SendNotification)
	r.Get("/users/{userID}/notifications", notificationHandler.GetNotifications)

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
