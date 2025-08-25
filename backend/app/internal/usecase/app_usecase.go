package usecase

import (
	"context"
	"fmt"
	"strconv" // Added for string to int conversion
	"time"

	"github.com/golang-jwt/jwt/v5" // Added for JWT token generation
	"github.com/google/uuid"
	"github.com/mkbagandov/kingsman/backend/app/internal/domain" // Update with your actual project path
	qrcode "github.com/skip2/go-qrcode"                          // Added QR code library
	"golang.org/x/crypto/bcrypt"                                 // Added for password hashing
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
	Username            string                    `json:"username"`
	Email               string                    `json:"email"`
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
		Username:            user.Username,
		PhoneNumber:         user.PhoneNumber,
		Email:               user.Email,
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

type UpdateUserProfileRequest struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type UpdateUserProfileResponse struct {
	Message string `json:"message"`
}

func (uc *UserUseCase) UpdateUserProfile(ctx context.Context, req *UpdateUserProfileRequest) (*UpdateUserProfileResponse, error) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user.Username = req.Username
	user.PhoneNumber = req.PhoneNumber
	user.Email = req.Email

	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return &UpdateUserProfileResponse{Message: "User profile updated successfully"}, nil
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
			Username:            user.Username,
			PhoneNumber:         user.PhoneNumber,
			DiscountLevel:       user.DiscountLevel,
			ProgressToNextLevel: user.ProgressToNextLevel,
			QRCode:              user.QRCode,
		},
		nil
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

type CartUseCase struct {
	cartRepo            domain.CartRepository
	cartItemRepo        domain.CartItemRepository
	productRepo         domain.ProductRepository
	orderRepo           domain.OrderRepository
	orderItemRepo       domain.OrderItemRepository
	loyaltyUseCase      *LoyaltyUseCase
	notificationUseCase *NotificationUseCase
	userRepo            domain.UserRepository
}

func NewCartUseCase(
	cartRepo domain.CartRepository,
	cartItemRepo domain.CartItemRepository,
	productRepo domain.ProductRepository,
	orderRepo domain.OrderRepository,
	orderItemRepo domain.OrderItemRepository,
	loyaltyUseCase *LoyaltyUseCase,
	notificationUseCase *NotificationUseCase,
	userRepo domain.UserRepository,
) *CartUseCase {
	return &CartUseCase{
		cartRepo:            cartRepo,
		cartItemRepo:        cartItemRepo,
		productRepo:         productRepo,
		orderRepo:           orderRepo,
		orderItemRepo:       orderItemRepo,
		loyaltyUseCase:      loyaltyUseCase,
		notificationUseCase: notificationUseCase,
		userRepo:            userRepo,
	}
}

func (uc *CartUseCase) GetCartByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	cart, err := uc.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart by user ID: %w", err)
	}

	// If cart doesn't exist, create one
	if cart == nil {
		newCart := &domain.Cart{UserID: userID, IsPaid: false} // Initialize IsPaid
		err := uc.cartRepo.CreateCart(ctx, newCart)
		if err != nil {
			return nil, fmt.Errorf("failed to create new cart for user: %w", err)
		}
		cart = newCart
	}

	return cart, nil
}

type AddItemToCartRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (uc *CartUseCase) AddItemToCart(ctx context.Context, req *AddItemToCartRequest) (*domain.CartItem, error) {
	cart, err := uc.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// Check if the product exists
	productIDInt, err := strconv.Atoi(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}

	product, err := uc.productRepo.GetProductByID(ctx, productIDInt)
	if err != nil || product == nil {
		return nil, fmt.Errorf("product with ID %s not found: %w", req.ProductID, err)
	}

	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	// Check if item already in cart
	cartItems, err := uc.cartItemRepo.GetCartItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	for _, item := range cartItems {
		if item.ProductID == req.ProductID {
			// Update quantity if item already exists
			item.Quantity += req.Quantity
			err := uc.cartItemRepo.UpdateCartItem(ctx, item)
			if err != nil {
				return nil, fmt.Errorf("failed to update cart item quantity: %w", err)
			}
			return item, nil
		}
	}

	// Add new item to cart
	cartItem := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := uc.cartItemRepo.CreateCartItem(ctx, cartItem); err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	return cartItem, nil
}

type UpdateCartItemRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (uc *CartUseCase) UpdateCartItem(ctx context.Context, req *UpdateCartItemRequest) error {
	cart, err := uc.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	cartItems, err := uc.cartItemRepo.GetCartItemsByCartID(ctx, cart.ID)
	if err != nil {
		return fmt.Errorf("failed to get cart items: %w", err)
	}

	found := false
	for _, item := range cartItems {
		if item.ProductID == req.ProductID {
			item.Quantity = req.Quantity
			err := uc.cartItemRepo.UpdateCartItem(ctx, item)
			if err != nil {
				return fmt.Errorf("failed to update cart item: %w", err)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("product with ID %s not found in cart", req.ProductID)
	}

	return nil
}

type RemoveCartItemRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
}

func (uc *CartUseCase) RemoveCartItem(ctx context.Context, req *RemoveCartItemRequest) error {
	cart, err := uc.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	err = uc.cartItemRepo.DeleteCartItemByCartIDAndProductID(ctx, cart.ID, req.ProductID)
	if err != nil {
		return fmt.Errorf("failed to remove item from cart: %w", err)
	}

	return nil
}

type GetCartResponse struct {
	Cart      *domain.Cart       `json:"cart"`
	CartItems []*domain.CartItem `json:"cart_items"`
}

func (uc *CartUseCase) GetUserCart(ctx context.Context, userID string) (*GetCartResponse, error) {
	cart, err := uc.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cartItems, err := uc.cartItemRepo.GetCartItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items for user: %w", err)
	}

	return &GetCartResponse{Cart: cart, CartItems: cartItems}, nil
}

func (uc *CartUseCase) ClearCart(ctx context.Context, userID string) error {
	cart, err := uc.GetCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	err = uc.cartRepo.DeleteCart(ctx, cart.ID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

type PlaceOrderRequest struct {
	UserID string `json:"user_id"`
}

type PlaceOrderResponse struct {
	OrderID int    `json:"order_id"`
	Message string `json:"message"`
}

func (uc *CartUseCase) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*PlaceOrderResponse, error) {
	cart, err := uc.cartRepo.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user cart: %w", err)
	}
	if cart == nil || cart.IsPaid {
		return nil, fmt.Errorf("cart is empty or already paid")
	}

	cartItems, err := uc.cartItemRepo.GetCartItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		productIDInt, err := strconv.Atoi(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID format: %w", err)
		}
		product, err := uc.productRepo.GetProductByID(ctx, productIDInt)
		if err != nil || product == nil {
			return nil, fmt.Errorf("product with ID %s not found: %w", item.ProductID, err)
		}
		totalAmount += product.Price * float64(item.Quantity)
	}

	// Create order
	order := &domain.Order{
		UserID:        req.UserID,
		OrderDate:     time.Now().Format(time.RFC3339),
		TotalAmount:   totalAmount,
		Status:        "completed", // Assuming successful payment
		PaymentStatus: "paid",
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}
	if err := uc.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Move cart items to order items
	for _, item := range cartItems {
		productIDInt, err := strconv.Atoi(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID format: %w", err)
		}
		product, err := uc.productRepo.GetProductByID(ctx, productIDInt)
		if err != nil || product == nil {
			return nil, fmt.Errorf("product with ID %s not found: %w", item.ProductID, err)
		}

		orderItem := &domain.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Store current product price at time of order
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		if err := uc.orderItemRepo.CreateOrderItem(ctx, orderItem); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Mark cart as paid
	cart.IsPaid = true
	if err := uc.cartRepo.UpdateCart(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to mark cart as paid: %w", err)
	}

	// Clear cart items (or delete the cart itself if preferred)
	// For simplicity, we'll delete the cart items for now
	for _, item := range cartItems {
		if err := uc.cartItemRepo.DeleteCartItem(ctx, item.ID); err != nil {
			return nil, fmt.Errorf("failed to delete cart item after order: %w", err)
		}
	}

	// Send notification to user
	userIDInt, err := strconv.Atoi(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format for notification: %w", err)
	}

	notificationReq := &SendNotificationRequest{
		UserID:  req.UserID,
		Type:    "purchase_confirmation",
		Title:   "Заказ успешно оплачен!",
		Message: fmt.Sprintf("Ваш заказ #%d на сумму $%.2f успешно оплачен и принят в обработку.", order.ID, totalAmount),
	}
	if err := uc.notificationUseCase.SendNotification(ctx, notificationReq); err != nil {
		return nil, fmt.Errorf("failed to send purchase confirmation notification: %w", err)
	}

	// Accrue loyalty points (e.g., 1 point per $10 spent)
	pointsToAccrue := int(totalAmount / 10)
	if pointsToAccrue > 0 {
		if err := uc.loyaltyUseCase.AddLoyaltyPoints(ctx, userIDInt, pointsToAccrue, "purchase"); err != nil {
			return nil, fmt.Errorf("failed to add loyalty points: %w", err)
		}
	}

	return &PlaceOrderResponse{OrderID: order.ID, Message: "Order placed successfully and cart marked as paid"}, nil
}

type OrderUseCase struct {
	orderRepo     domain.OrderRepository
	orderItemRepo domain.OrderItemRepository
	productRepo   domain.ProductRepository // To fetch product details for order items
}

func NewOrderUseCase(orderRepo domain.OrderRepository, orderItemRepo domain.OrderItemRepository, productRepo domain.ProductRepository) *OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo, orderItemRepo: orderItemRepo, productRepo: productRepo}
}

type GetOrdersResponse struct {
	Orders []*domain.Order `json:"orders"`
}

func (uc *OrderUseCase) GetOrdersByUserID(ctx context.Context, userID string) (*GetOrdersResponse, error) {
	paidStatus := "paid"
	orders, err := uc.orderRepo.GetOrdersByUserID(ctx, userID, &paidStatus) // Pass 'paid' status
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for user: %w", err)
	}

	// For each order, fetch its items
	for _, order := range orders {
		orderItems, err := uc.orderItemRepo.GetOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items for order %d: %w", order.ID, err)
		}
		// Fetch product details for each order item
		for _, item := range orderItems {
			productIDInt, err := strconv.Atoi(item.ProductID)
			if err != nil {
				return nil, fmt.Errorf("invalid product ID format for order item: %w", err)
			}
			_, err = uc.productRepo.GetProductByID(ctx, productIDInt) // Changed product to _
			if err != nil {
				// Log error but don't fail the entire order retrieval if product not found
				fmt.Printf("Warning: Product with ID %s not found for order item %d\n", item.ProductID, item.ID)
				continue
			}
			// Assign the product to the order item for a richer response
			// This requires a Product field in domain.OrderItem struct, which is not there.
			// For now, we'll just return the order items as is.
			// If we want to embed product details, we need to modify domain.OrderItem or create a response struct.
		}
		// Assign items to the order
		var plainOrderItems []domain.OrderItem
		for _, item := range orderItems {
			plainOrderItems = append(plainOrderItems, *item)
		}
		order.Items = plainOrderItems // Assign the converted slice
	}

	return &GetOrdersResponse{Orders: orders}, nil
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
