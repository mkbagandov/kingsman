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
