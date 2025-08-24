package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
	qrcode "github.com/skip2/go-qrcode" // Added QR code library
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

type RegisterUserRequest struct {
	PhoneNumber string `json:"phone_number"`
	SocialID    string `json:"social_id,omitempty"`
}

type RegisterUserResponse struct {
	UserID string `json:"user_id"`
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error) {
	// Check if user already exists by phone number
	existingUser, err := uc.userRepo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil && err.Error() != "user not found" {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with phone number %s already exists", req.PhoneNumber)
	}

	user := &domain.User{
		ID:                  uuid.New().String(),
		PhoneNumber:         req.PhoneNumber,
		SocialID:            req.SocialID,
		DiscountLevel:       0,                   // Initial discount level
		ProgressToNextLevel: 0.0,                 // Initial progress
		QRCode:              uuid.New().String(), // Generate a unique QR code
	}

	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &RegisterUserResponse{UserID: user.ID}, nil
}

type GetUserProfileResponse struct {
	ID                  string  `json:"id"`
	PhoneNumber         string  `json:"phone_number"`
	DiscountLevel       int     `json:"discount_level"`
	ProgressToNextLevel float64 `json:"progress_to_next_level"`
	QRCode              string  `json:"qr_code"`
}

func (uc *UserUseCase) GetUserProfile(ctx context.Context, userID string) (*GetUserProfileResponse, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &GetUserProfileResponse{
		ID:                  user.ID,
		PhoneNumber:         user.PhoneNumber,
		DiscountLevel:       user.DiscountLevel,
		ProgressToNextLevel: user.ProgressToNextLevel,
		QRCode:              user.QRCode,
	}, nil
}

// GetUserQRCode retrieves the QR code string for a user.
func (uc *UserUseCase) GetUserQRCode(ctx context.Context, userID string) (string, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user QR code: %w", err)
	}
	return user.QRCode, nil
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
	user, err := uc.userRepo.GetUserDiscountCard(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user discount card: %w", err)
	}

	return &GetUserProfileResponse{
		ID:                  user.ID,
		PhoneNumber:         user.PhoneNumber,
		DiscountLevel:       user.DiscountLevel,
		ProgressToNextLevel: user.ProgressToNextLevel,
		QRCode:              user.QRCode,
	}, nil
}

// UpdateUserDiscountCard updates the discount card level and progress for a user.
func (uc *UserUseCase) UpdateUserDiscountCard(ctx context.Context, userID string, level int, progress float64) error {
	if err := uc.userRepo.UpdateUserDiscountCard(ctx, userID, level, progress); err != nil {
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
	store, err := uc.storeRepo.GetStoreByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get store by ID: %w", err)
	}

	return &GetStoreByIDResponse{Store: store}, nil
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
	notification := &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
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
	notifications, err := uc.notificationRepo.GetNotificationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return &GetNotificationsResponse{Notifications: notifications}, nil
}
