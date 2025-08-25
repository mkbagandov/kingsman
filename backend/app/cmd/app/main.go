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
	"github.com/go-chi/cors" // Import the CORS package

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

var jwtKey = []byte("YOUR_SECRET_KEY") // TODO: Load from environment variable

func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]
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
		ctx := context.WithValue(r.Context(), domain.UserContextKey, claims.UserID)
		r = r.WithContext(ctx)

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
	cartRepo := infrastructure.NewCartRepository(db)
	cartItemRepo := infrastructure.NewCartItemRepository(db)
	orderRepo := infrastructure.NewOrderRepository(db)         // Initialize OrderRepository
	orderItemRepo := infrastructure.NewOrderItemRepository(db) // Initialize OrderItemRepository

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	storeUseCase := usecase.NewStoreUseCase(storeRepo)
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	productUseCase := usecase.NewProductUseCase(productRepo)
	loyaltyUseCase := usecase.NewLoyaltyUseCase(userRepo)                                                                                               // Initialize LoyaltyUseCase
	cartUseCase := usecase.NewCartUseCase(cartRepo, cartItemRepo, productRepo, orderRepo, orderItemRepo, loyaltyUseCase, notificationUseCase, userRepo) // Updated CartUseCase
	orderUseCase := usecase.NewOrderUseCase(orderRepo, orderItemRepo, productRepo)                                                                      // Initialize OrderUseCase

	// Initialize handlers
	userHandler := delivery.NewUserHandler(userUseCase, loyaltyUseCase) // Pass loyaltyUseCase
	storeHandler := delivery.NewStoreHandler(storeUseCase)
	notificationHandler := delivery.NewNotificationHandler(notificationUseCase)
	categoryHandler := delivery.NewCategoryHandler(categoryUseCase)
	productHandler := delivery.NewProductHandler(productUseCase)
	cartHandler := delivery.NewCartHandler(cartUseCase)
	orderHandler := delivery.NewOrderHandler(orderUseCase) // Initialize OrderHandler

	// Setup HTTP router
	r := chi.NewRouter()

	// Configure CORS middleware
	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(corsHandler) // Apply the CORS middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Kingsman Backend!"))
	})

	// Public routes (registration and login)
	r.Post("/users/register", userHandler.RegisterUser)
	r.Post("/users/login", userHandler.LoginUser)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)

		// User routes
		r.Get("/users/profile", userHandler.GetUserProfile)
		r.Get("/users/discount-card", userHandler.GetUserDiscountCard)
		r.Put("/users/discount-card", userHandler.UpdateUserDiscountCard)
		r.Get("/users/qrcode", userHandler.GetUserQRCode)

		// Loyalty routes
		r.Get("/users/loyalty", userHandler.GetUserLoyaltyProfile)
		r.Post("/users/loyalty-points", userHandler.AddLoyaltyPoints)
		r.Post("/users/loyalty-activity", userHandler.AddLoyaltyActivity)
		r.Get("/loyalty-tiers", userHandler.GetLoyaltyTiers)

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
		r.Get("/users/notifications", notificationHandler.GetNotifications)

		// Cart routes
		r.Route("/cart", func(r chi.Router) {
			r.Post("/checkout", cartHandler.PlaceOrder) // New route for placing an order
			r.Post("/items", cartHandler.AddItemToCart)
			r.Put("/items", cartHandler.UpdateCartItem)
			r.Delete("/items/{productID}", cartHandler.RemoveCartItem)
			r.Get("/", cartHandler.GetUserCart)
			r.Delete("/clear", cartHandler.ClearCart)
		})

		// Order routes
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", orderHandler.GetUserOrders) // Route to get user's order history
		})
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
