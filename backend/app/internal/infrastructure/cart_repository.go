package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
)

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) domain.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	query := `INSERT INTO carts (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`
	cart.ID = uuid.New().String()
	cart.CreatedAt = time.Now().Format(time.RFC3339)
	cart.UpdatedAt = time.Now().Format(time.RFC3339)

	err := r.db.QueryRowContext(ctx, query, cart.UserID, cart.CreatedAt, cart.UpdatedAt).Scan(&cart.ID)
	if err != nil {
		return fmt.Errorf("failed to create cart: %w", err)
	}
	return nil
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`
	cart := &domain.Cart{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Cart not found
		}
		return nil, fmt.Errorf("failed to get cart by user ID: %w", err)
	}
	return cart, nil
}

func (r *cartRepository) UpdateCart(ctx context.Context, cart *domain.Cart) error {
	query := `UPDATE carts SET updated_at = $2 WHERE id = $1`
	cart.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, query, cart.ID, cart.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}
	return nil
}

func (r *cartRepository) DeleteCart(ctx context.Context, cartID string) error {
	query := `DELETE FROM carts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}
