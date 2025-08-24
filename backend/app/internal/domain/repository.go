package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	UpdateUserDiscountCard(ctx context.Context, userID string, level int, progress float64) error
	GetUserDiscountCard(ctx context.Context, userID string) (*User, error) // Can return a User with only discount card fields populated
}

type StoreRepository interface {
	GetStores(ctx context.Context) ([]*Store, error)
	GetStoreByID(ctx context.Context, id string) (*Store, error)
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotificationsByUserID(ctx context.Context, userID string) ([]*Notification, error)
}
