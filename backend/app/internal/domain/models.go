package domain

type User struct {
	ID                  string  `json:"id"`
	PhoneNumber         string  `json:"phone_number"`
	SocialID            string  `json:"social_id,omitempty"` // For social logins
	DiscountLevel       int     `json:"discount_level"`
	ProgressToNextLevel float64 `json:"progress_to_next_level"`
	QRCode              string  `json:"qr_code"`
}

type Store struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Location string `json:"location"` // e.g., latitude, longitude
	Phone    string `json:"phone"`
}

type Notification struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"` // e.g., "new_arrival", "promotion"
	Title     string `json:"title"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
