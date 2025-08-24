package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
)

type PostgreSQLUserRepository struct {
	db *sql.DB
}

func NewPostgreSQLUserRepository(db *sql.DB) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{db: db}
}

func (r *PostgreSQLUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.PhoneNumber, user.Email, user.PasswordHash, user.SocialID, user.DiscountLevel, user.ProgressToNextLevel, user.QRCode).Scan(&user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			// Check if the unique violation is for phone_number or email
			if strings.Contains(pqErr.Detail, "phone_number") {
				return fmt.Errorf("user with phone number %s already exists", user.PhoneNumber)
			} else if strings.Contains(pqErr.Detail, "email") {
				return fmt.Errorf("user with email %s already exists", user.Email)
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.SocialID, &user.DiscountLevel, &user.ProgressToNextLevel, &user.QRCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code FROM users WHERE phone_number = $1`
	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.SocialID, &user.DiscountLevel, &user.ProgressToNextLevel, &user.QRCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by phone number: %w", err)
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.SocialID, &user.DiscountLevel, &user.ProgressToNextLevel, &user.QRCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET phone_number = $2, email = $3, password_hash = $4, social_id = $5, discount_level = $6, progress_to_next_level = $7, qr_code = $8 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.PhoneNumber, user.Email, user.PasswordHash, user.SocialID, user.DiscountLevel, user.ProgressToNextLevel, user.QRCode)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) UpdateUserDiscountCard(ctx context.Context, userID int, level int, progress float64) error {
	query := `UPDATE users SET discount_level = $2, progress_to_next_level = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID, level, progress)
	if err != nil {
		return fmt.Errorf("failed to update user discount card: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) GetUserDiscountCard(ctx context.Context, userID int) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT discount_level, progress_to_next_level, qr_code FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.DiscountLevel, &user.ProgressToNextLevel, &user.QRCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user discount card: %w", err)
	}
	return user, nil
}
