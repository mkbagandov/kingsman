package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error) // New method for authentication
	UpdateUser(ctx context.Context, user *User) error
	UpdateUserDiscountCard(ctx context.Context, userID int, level int, progress float64) error
	GetUserDiscountCard(ctx context.Context, userID int) (*User, error) // Can return a User with only discount card fields populated

	// Loyalty Program methods
	CreateLoyaltyPoint(ctx context.Context, point *LoyaltyPoint) error
	GetLoyaltyPointsByUserID(ctx context.Context, userID int) ([]*LoyaltyPoint, error)
	GetLoyaltyTierByID(ctx context.Context, id int) (*LoyaltyTier, error)
	GetLoyaltyTierByName(ctx context.Context, name string) (*LoyaltyTier, error)
	GetAllLoyaltyTiers(ctx context.Context) ([]*LoyaltyTier, error)
	CreateLoyaltyTier(ctx context.Context, tier *LoyaltyTier) error
	UpdateLoyaltyTier(ctx context.Context, tier *LoyaltyTier) error
	DeleteLoyaltyTier(ctx context.Context, id int) error
	CreateLoyaltyActivity(ctx context.Context, activity *LoyaltyActivity) error
	GetLoyaltyActivitiesByUserID(ctx context.Context, userID int) ([]*LoyaltyActivity, error)
	GetUserLoyalty(ctx context.Context, userID int) (*UserLoyalty, error)
	UpdateUserLoyalty(ctx context.Context, userLoyalty *UserLoyalty) error
}

type StoreRepository interface {
	GetStores(ctx context.Context) ([]*Store, error)
	GetStoreByID(ctx context.Context, id int) (*Store, error)
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotificationsByUserID(ctx context.Context, userID int) ([]*Notification, error)
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *Category) error
	GetCategoryByID(ctx context.Context, id int) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, id int) (*Product, error)
	GetProducts(ctx context.Context, categoryID *string, minPrice *float64, maxPrice *float64, sortBy *string, sortOrder *string, limit, offset int) ([]*Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type CartRepository interface {
	CreateCart(ctx context.Context, cart *Cart) error
	GetCartByUserID(ctx context.Context, userID string) (*Cart, error)
	UpdateCart(ctx context.Context, cart *Cart) error
	DeleteCart(ctx context.Context, cartID string) error
}

type CartItemRepository interface {
	CreateCartItem(ctx context.Context, cartItem *CartItem) error
	GetCartItemsByCartID(ctx context.Context, cartID string) ([]*CartItem, error)
	UpdateCartItem(ctx context.Context, cartItem *CartItem) error
	DeleteCartItem(ctx context.Context, cartItemID string) error
	DeleteCartItemByCartIDAndProductID(ctx context.Context, cartID, productID string) error
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, orderID int) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID string, paymentStatus *string) ([]*Order, error) // Added paymentStatus filter
	UpdateOrder(ctx context.Context, order *Order) error
}

type OrderItemRepository interface {
	CreateOrderItem(ctx context.Context, orderItem *OrderItem) error
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]*OrderItem, error)
	UpdateOrderItem(ctx context.Context, orderItem *OrderItem) error
	DeleteOrderItem(ctx context.Context, orderItemID int) error
}
