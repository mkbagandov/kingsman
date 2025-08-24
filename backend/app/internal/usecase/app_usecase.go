package usecase

import (
	"context"
	"fmt"
	"strconv" // Added for string to int conversion
	"time"

	"github.com/golang-jwt/jwt/v5" // Added for JWT token generation
	"github.com/google/uuid"
	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
	qrcode "github.com/skip2/go-qrcode" // Added QR code library
	"golang.org/x/crypto/bcrypt"        // Added for password hashing
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// LoyaltyUseCase handles loyalty program related business logic.
type LoyaltyUseCase struct {
	userRepo domain.UserRepository // Reusing UserRepository for loyalty data
}

// NewLoyaltyUseCase creates a new LoyaltyUseCase.
func NewLoyaltyUseCase(userRepo domain.UserRepository) *LoyaltyUseCase {
	return &LoyaltyUseCase{userRepo: userRepo}
}

type RegisterUserRequest struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	SocialID    string `json:"social_id,omitempty"`
}

type RegisterUserResponse struct {
	UserID string `json:"user_id"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error) {
	// Check if user already exists by phone number
	fmt.Println(req)
	existingUser, err := uc.userRepo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	fmt.Println("1111111111111111111")
	fmt.Println(err)
	fmt.Println("111111111")
	fmt.Println(existingUser)
	fmt.Println("2222222222")
	if err != nil && err.Error() != "user not found" {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with phone number %s already exists", req.PhoneNumber)
	}

	// Check if user already exists by email
	existingUser, err = uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && err.Error() != "user not found" {
		return nil, fmt.Errorf("failed to check existing user by email: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// The database will generate the ID, so we set it to 0 initially
	// Generate a unique QR code string
	qrCodeString := uuid.New().String()
	fmt.Println(req)
	user := &domain.User{
		Username:            req.Username,
		PhoneNumber:         req.PhoneNumber,
		Email:               req.Email,
		PasswordHash:        string(hashedPassword),
		SocialID:            req.SocialID,
		DiscountLevel:       0,
		ProgressToNextLevel: 0.0,
		QRCode:              &qrCodeString,
		LoyaltyStatus:       "Bronze",
		CurrentPoints:       0,
	}

	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Initialize user loyalty entry
	userLoyalty := &domain.UserLoyalty{
		UserID:         user.ID,
		CurrentPoints:  0,
		CurrentTierID:  1, // Will be updated when tiers are defined and assigned
		LastActivityAt: time.Now().Format(time.RFC3339),
	}
	if err := uc.userRepo.UpdateUserLoyalty(ctx, userLoyalty); err != nil {
		return nil, fmt.Errorf("failed to initialize user loyalty: %w", err)
	}

	return &RegisterUserResponse{UserID: strconv.Itoa(user.ID)}, nil
}

// LoginUser authenticates a user and generates a JWT token.
func (uc *UserUseCase) LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		UserID: strconv.Itoa(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("YOUR_SECRET_KEY")) // TODO: Use a secure secret key from config
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginUserResponse{Token: tokenString}, nil
}

type GetUserProfileResponse struct {
	ID                  string                    `json:"id"`
	PhoneNumber         string                    `json:"phone_number"`
	DiscountLevel       int                       `json:"discount_level"`
	ProgressToNextLevel float64                   `json:"progress_to_next_level"`
	QRCode              *string                   `json:"qr_code"` // Changed to *string to match domain.User
	LoyaltyStatus       string                    `json:"loyalty_status"`
	CurrentPoints       int                       `json:"current_points"`
	CurrentTier         *LoyaltyTierResponse      `json:"current_tier,omitempty"`
	LoyaltyActivities   []*domain.LoyaltyActivity `json:"loyalty_activities,omitempty"`
}

type LoyaltyTierResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MinPoints   int    `json:"min_points"`
	Description string `json:"description,omitempty"`
	Benefits    string `json:"benefits,omitempty"`
}

func (uc *UserUseCase) GetUserProfile(ctx context.Context, userID string) (*GetUserProfileResponse, error) {
	// Convert userID string to int for repository call
	id, err := strconv.Atoi(userID)
	fmt.Println("1131313131313")
	fmt.Println(id)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format: %w", err)
	}
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Get user loyalty information
	userLoyalty, err := uc.userRepo.GetUserLoyalty(ctx, id)
	if err != nil && err.Error() != "user loyalty not found" {
		return nil, fmt.Errorf("failed to get user loyalty information: %w", err)
	}

	var currentTier *LoyaltyTierResponse
	if userLoyalty != nil && userLoyalty.CurrentTierID != 0 {
		tier, err := uc.userRepo.GetLoyaltyTierByID(ctx, userLoyalty.CurrentTierID)
		if err != nil && err.Error() != "loyalty tier not found" {
			return nil, fmt.Errorf("failed to get loyalty tier: %w", err)
		}
		if tier != nil {
			currentTier = &LoyaltyTierResponse{
				ID:          tier.ID,
				Name:        tier.Name,
				MinPoints:   tier.MinPoints,
				Description: tier.Description,
				Benefits:    tier.Benefits,
			}
		}
	}

	// Get loyalty activities
	activities, err := uc.userRepo.GetLoyaltyActivitiesByUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get loyalty activities: %w", err)
	}

	response := &GetUserProfileResponse{
		ID:                  strconv.Itoa(user.ID),
		PhoneNumber:         user.PhoneNumber,
		DiscountLevel:       user.DiscountLevel,
		ProgressToNextLevel: user.ProgressToNextLevel,
		QRCode:              user.QRCode,
		LoyaltyStatus:       user.LoyaltyStatus,
		CurrentPoints:       user.CurrentPoints,
		CurrentTier:         currentTier,
		LoyaltyActivities:   activities,
	}

	return response, nil
}

// GetUserQRCode retrieves the QR code string for a user.
func (uc *UserUseCase) GetUserQRCode(ctx context.Context, userID string) (string, error) {
	// Convert userID string to int for repository call
	id, err := strconv.Atoi(userID)
	if err != nil {
		return "", fmt.Errorf("invalid userID format: %w", err)
	}
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get user QR code: %w", err)
	}
	if user.QRCode == nil {
		return "", fmt.Errorf("QR code not found for user %s", userID)
	}
	return *user.QRCode, nil
}

// GenerateQRCodeImage generates a QR code image (PNG) for a given user's QR code string.
func (uc *UserUseCase) GenerateQRCodeImage(ctx context.Context, userID string) ([]byte, error) {
	qrCodeString, err := uc.GetUserQRCode(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get QR code string: %w", err)
	}

	p, err := qrcode.Encode(qrCodeString, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code image: %w", err)
	}
	return p, nil
}

// GetUserDiscountCard retrieves the discount card information for a user.
func (uc *UserUseCase) GetUserDiscountCard(ctx context.Context, userID string) (*GetUserProfileResponse, error) {
	// Convert userID string to int for repository call
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format: %w", err)
	}
	user, err := uc.userRepo.GetUserDiscountCard(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user discount card: %w", err)
	}

	return &GetUserProfileResponse{
		ID:                  strconv.Itoa(user.ID),
		PhoneNumber:         user.PhoneNumber,
		DiscountLevel:       user.DiscountLevel,
		ProgressToNextLevel: user.ProgressToNextLevel,
		QRCode:              user.QRCode,
	}, nil
}

// UpdateUserDiscountCard updates the discount card level and progress for a user.
func (uc *UserUseCase) UpdateUserDiscountCard(ctx context.Context, userID string, level int, progress float64) error {
	// Convert userID string to int for repository call
	id, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format: %w", err)
	}
	if err := uc.userRepo.UpdateUserDiscountCard(ctx, id, level, progress); err != nil {
		return fmt.Errorf("failed to update user discount card: %w", err)
	}
	return nil
}

type StoreUseCase struct {
	storeRepo domain.StoreRepository
}

func NewStoreUseCase(storeRepo domain.StoreRepository) *StoreUseCase {
	return &StoreUseCase{storeRepo: storeRepo}
}

type GetStoresResponse struct {
	Stores []*domain.Store `json:"stores"`
}

func (uc *StoreUseCase) GetStores(ctx context.Context) (*GetStoresResponse, error) {
	stores, err := uc.storeRepo.GetStores(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stores: %w", err)
	}

	return &GetStoresResponse{Stores: stores}, nil
}

type GetStoreByIDResponse struct {
	Store *domain.Store `json:"store"`
}

func (uc *StoreUseCase) GetStoreByID(ctx context.Context, id string) (*GetStoreByIDResponse, error) {
	// Convert id string to int for repository call
	storeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid store ID format: %w", err)
	}
	store, err := uc.storeRepo.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store by ID: %w", err)
	}

	return &GetStoreByIDResponse{Store: store}, nil
}

// New Category Use Case
type CategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUseCase(categoryRepo domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: categoryRepo}
}

type GetCategoriesResponse struct {
	Categories []*domain.Category `json:"categories"`
}

func (uc *CategoryUseCase) GetCategories(ctx context.Context) (*GetCategoriesResponse, error) {
	categories, err := uc.categoryRepo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return &GetCategoriesResponse{Categories: categories}, nil
}

// New Product Use Case
type ProductUseCase struct {
	productRepo domain.ProductRepository
}

func NewProductUseCase(productRepo domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

type GetProductCatalogRequest struct {
	CategoryID *string  `json:"category_id,omitempty"`
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	SortBy     *string  `json:"sort_by,omitempty"`
	SortOrder  *string  `json:"sort_order,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
}

type GetProductCatalogResponse struct {
	Products []*domain.Product `json:"products"`
}

func (uc *ProductUseCase) GetProductCatalog(ctx context.Context, req *GetProductCatalogRequest) (*GetProductCatalogResponse, error) {
	products, err := uc.productRepo.GetProducts(ctx, req.CategoryID, req.MinPrice, req.MaxPrice, req.SortBy, req.SortOrder, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get product catalog: %w", err)
	}

	return &GetProductCatalogResponse{Products: products}, nil
}

type GetProductByIDResponse struct {
	Product *domain.Product `json:"product"`
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, productID string) (*GetProductByIDResponse, error) {
	// Convert productID string to int for repository call
	id, err := strconv.Atoi(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}
	product, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return &GetProductByIDResponse{Product: product}, nil
}

type NotificationUseCase struct {
	notificationRepo domain.NotificationRepository
}

func NewNotificationUseCase(notificationRepo domain.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{notificationRepo: notificationRepo}
}

type SendNotificationRequest struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (uc *NotificationUseCase) SendNotification(ctx context.Context, req *SendNotificationRequest) error {
	// Convert req.UserID string to int
	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid UserID format for notification: %w", err)
	}

	notification := &domain.Notification{
		ID:        0, // ID will be auto-generated by the database
		UserID:    userID,
		Type:      req.Type,
		Title:     req.Title,
		Message:   req.Message,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := uc.notificationRepo.CreateNotification(ctx, notification); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}

type GetNotificationsResponse struct {
	Notifications []*domain.Notification `json:"notifications"`
}

func (uc *NotificationUseCase) GetNotifications(ctx context.Context, userID string) (*GetNotificationsResponse, error) {
	// Convert userID string to int for repository call
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format: %w", err)
	}
	notifications, err := uc.notificationRepo.GetNotificationsByUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return &GetNotificationsResponse{Notifications: notifications}, nil
}

// AddLoyaltyPoints adds loyalty points to a user and updates their tier if necessary.
func (uc *LoyaltyUseCase) AddLoyaltyPoints(ctx context.Context, userID int, points int, pointType string) error {
	// Create loyalty point record
	point := &domain.LoyaltyPoint{
		UserID: userID,
		Points: points,
		Type:   pointType,
	}
	if err := uc.userRepo.CreateLoyaltyPoint(ctx, point); err != nil {
		return fmt.Errorf("failed to create loyalty point: %w", err)
	}

	// Get user loyalty and update current points
	userLoyalty, err := uc.userRepo.GetUserLoyalty(ctx, userID)
	if err != nil && err.Error() == "user loyalty not found" {
		// Initialize if not found
		userLoyalty = &domain.UserLoyalty{
			UserID:         userID,
			CurrentPoints:  points,
			LastActivityAt: time.Now().Format(time.RFC3339),
		}
	} else if err != nil {
		return fmt.Errorf("failed to get user loyalty: %w", err)
	} else {
		userLoyalty.CurrentPoints += points
		userLoyalty.LastActivityAt = time.Now().Format(time.RFC3339)
	}

	// Determine new tier
	tiers, err := uc.userRepo.GetAllLoyaltyTiers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all loyalty tiers: %w", err)
	}

	var newTier *domain.LoyaltyTier
	for _, tier := range tiers {
		if userLoyalty.CurrentPoints >= tier.MinPoints {
			newTier = tier
		} else {
			break // Tiers are sorted by min_points, so we can stop.
		}
	}

	if newTier != nil && (userLoyalty.CurrentTierID == 0 || userLoyalty.CurrentTierID != newTier.ID) {
		userLoyalty.CurrentTierID = newTier.ID
		// Update user's loyalty status in the main user table as well
		user, err := uc.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user for tier update: %w", err)
		}
		user.LoyaltyStatus = newTier.Name
		user.CurrentPoints = userLoyalty.CurrentPoints
		if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
			return fmt.Errorf("failed to update user loyalty status: %w", err)
		}
		// Add loyalty activity for tier upgrade
		activity := &domain.LoyaltyActivity{
			UserID:      userID,
			Type:        "tier_upgrade",
			Description: fmt.Sprintf("Upgraded to %s tier", newTier.Name),
		}
		if err := uc.userRepo.CreateLoyaltyActivity(ctx, activity); err != nil {
			return fmt.Errorf("failed to create tier upgrade activity: %w", err)
		}
	}

	if err := uc.userRepo.UpdateUserLoyalty(ctx, userLoyalty); err != nil {
		return fmt.Errorf("failed to update user loyalty: %w", err)
	}

	return nil
}

// AddLoyaltyActivity records a loyalty-related activity for a user.
func (uc *LoyaltyUseCase) AddLoyaltyActivity(ctx context.Context, userID int, activityType, description string) error {
	activity := &domain.LoyaltyActivity{
		UserID:      userID,
		Type:        activityType,
		Description: description,
	}
	if err := uc.userRepo.CreateLoyaltyActivity(ctx, activity); err != nil {
		return fmt.Errorf("failed to create loyalty activity: %w", err)
	}
	return nil
}

// GetUserLoyaltyProfile retrieves a comprehensive loyalty profile for a user.
func (uc *LoyaltyUseCase) GetUserLoyaltyProfile(ctx context.Context, userID int) (*GetUserProfileResponse, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	userLoyalty, err := uc.userRepo.GetUserLoyalty(ctx, userID)
	if err != nil && err.Error() != "user loyalty not found" {
		return nil, fmt.Errorf("failed to get user loyalty information: %w", err)
	}

	var currentTier *LoyaltyTierResponse
	if userLoyalty != nil && userLoyalty.CurrentTierID != 0 {
		tier, err := uc.userRepo.GetLoyaltyTierByID(ctx, userLoyalty.CurrentTierID)
		if err != nil && err.Error() != "loyalty tier not found" {
			return nil, fmt.Errorf("failed to get loyalty tier: %w", err)
		}
		if tier != nil {
			currentTier = &LoyaltyTierResponse{
				ID:          tier.ID,
				Name:        tier.Name,
				MinPoints:   tier.MinPoints,
				Description: tier.Description,
				Benefits:    tier.Benefits,
			}
		}
	}

	activities, err := uc.userRepo.GetLoyaltyActivitiesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loyalty activities: %w", err)
	}

	response := &GetUserProfileResponse{
		ID:                  strconv.Itoa(user.ID),
		PhoneNumber:         user.PhoneNumber,
		DiscountLevel:       user.DiscountLevel,
		ProgressToNextLevel: user.ProgressToNextLevel,
		QRCode:              user.QRCode,
		LoyaltyStatus:       user.LoyaltyStatus,
		CurrentPoints:       user.CurrentPoints,
		CurrentTier:         currentTier,
		LoyaltyActivities:   activities,
	}

	return response, nil
}

// GetLoyaltyTiers retrieves all defined loyalty tiers.
func (uc *LoyaltyUseCase) GetLoyaltyTiers(ctx context.Context) ([]*LoyaltyTierResponse, error) {
	tiers, err := uc.userRepo.GetAllLoyaltyTiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all loyalty tiers: %w", err)
	}

	var tierResponses []*LoyaltyTierResponse
	for _, tier := range tiers {
		tierResponses = append(tierResponses, &LoyaltyTierResponse{
			ID:          tier.ID,
			Name:        tier.Name,
			MinPoints:   tier.MinPoints,
			Description: tier.Description,
			Benefits:    tier.Benefits,
		})
	}
	return tierResponses, nil
}
