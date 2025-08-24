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

	// "github.com/golang-migrate/migrate/v4"                     // Added golang-migrate
	// _ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL driver for migrate
	// _ "github.com/golang-migrate/migrate/v4/source/file"       // File source for migrate
	"github.com/golang-jwt/jwt/v5"                               // Added golang-jwt
	"github.com/mkbagandov/kingsman/backend/app/internal/domain" // Added domain

	"github.com/joho/godotenv" // Added godotenv

	"github.com/mkbagandov/kingsman/backend/app/internal/delivery"
	"github.com/mkbagandov/kingsman/backend/app/internal/infrastructure"
	"github.com/mkbagandov/kingsman/backend/app/internal/usecase"
)

var jwtKey = []byte("YOUR_ULTRA_SECURE_SECRET_KEY") // TODO: Load from environment variable

func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix

		claims := &domain.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context for subsequent handlers
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, loading from environment: %v", err)
	}

	// Load configuration (e.g., from environment variables or a config file)
	dbConfig := &infrastructure.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432, // Default PostgreSQL port
		User:     "postgres",
		Password: "a6fbnmod",
		DBName:   "kingsman",
		SSLMode:  "disable",
	}

	// Initialize database connection
	db, err := infrastructure.NewPostgreSQLDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Construct database URL for migrations
	// dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
	// 	dbConfig.User,
	// 	dbConfig.Password,
	// 	dbConfig.Host,
	// 	dbConfig.Port,
	// 	dbConfig.DBName,
	// 	dbConfig.SSLMode,
	// )

	// Run database migrations
	// m, err := migrate.New(
	// 	"file://db/migrations",
	// 	dbURL,
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to create migrate instance: %v", err)
	// }
	// if err = m.Up(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatalf("Failed to run database migrations: %v", err)
	// }
	// log.Println("Database migrations applied successfully!")

	// Initialize repositories
	userRepo := infrastructure.NewPostgreSQLUserRepository(db)
	storeRepo := infrastructure.NewPostgreSQLStoreRepository(db)
	notificationRepo := infrastructure.NewPostgreSQLNotificationRepository(db)
	categoryRepo := infrastructure.NewPostgreSQLCategoryRepository(db)
	productRepo := infrastructure.NewPostgreSQLProductRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	storeUseCase := usecase.NewStoreUseCase(storeRepo)
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	productUseCase := usecase.NewProductUseCase(productRepo)

	// Initialize handlers
	userHandler := delivery.NewUserHandler(userUseCase)
	storeHandler := delivery.NewStoreHandler(storeUseCase)
	notificationHandler := delivery.NewNotificationHandler(notificationUseCase)
	categoryHandler := delivery.NewCategoryHandler(categoryUseCase)
	productHandler := delivery.NewProductHandler(productUseCase)

	// Setup HTTP router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Kingsman Backend!"))
	})

	// Public routes (registration and login)
	r.Post("/register", userHandler.RegisterUser)
	r.Post("/login", userHandler.LoginUser)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)

		// User routes
		r.Get("/users/{userID}", userHandler.GetUserProfile)
		r.Get("/users/{userID}/discount-card", userHandler.GetUserDiscountCard)
		r.Put("/users/{userID}/discount-card", userHandler.UpdateUserDiscountCard)
		r.Get("/users/{userID}/qrcode", userHandler.GetUserQRCode)

		// Store routes
		r.Get("/stores", storeHandler.GetStores)
		r.Get("/stores/{storeID}", storeHandler.GetStoreByID)

		// Category routes
		r.Get("/categories", categoryHandler.GetCategories)

		// Product routes
		r.Get("/products", productHandler.GetProductCatalog)
		r.Get("/products/{productID}", productHandler.GetProductByID)

		// Notification routes
		r.Post("/notifications", notificationHandler.SendNotification)
		r.Get("/users/{userID}/notifications", notificationHandler.GetNotifications)
	})

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
