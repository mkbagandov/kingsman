package domain

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID                  int     `json:"id"`
	PhoneNumber         string  `json:"phone_number"`
	Email               string  `json:"email,omitempty"`     // New field for user email
	PasswordHash        string  `json:"-"`                   // New field for storing hashed password, omit from JSON
	SocialID            string  `json:"social_id,omitempty"` // For social logins
	DiscountLevel       int     `json:"discount_level"`
	ProgressToNextLevel float64 `json:"progress_to_next_level"`
	QRCode              *string `json:"qr_code,omitempty"`
}

type Store struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Location string `json:"location"` // e.g., latitude, longitude
	Phone    string `json:"phone"`
}

type Notification struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Type      string `json:"type"` // e.g., "new_arrival", "promotion"
	Title     string `json:"title"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	ImageURL    string  `json:"image_url,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
