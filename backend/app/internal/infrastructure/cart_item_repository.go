package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
)

type cartItemRepository struct {
	db *sql.DB
}

func NewCartItemRepository(db *sql.DB) domain.CartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) CreateCartItem(ctx context.Context, cartItem *domain.CartItem) error {
	query := `INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	cartItem.ID = uuid.New().String()
	cartItem.CreatedAt = time.Now().Format(time.RFC3339)
	cartItem.UpdatedAt = time.Now().Format(time.RFC3339)

	err := r.db.QueryRowContext(ctx, query, cartItem.CartID, cartItem.ProductID, cartItem.Quantity, cartItem.CreatedAt, cartItem.UpdatedAt).Scan(&cartItem.ID)
	if err != nil {
		return fmt.Errorf("failed to create cart item: %w", err)
	}
	return nil
}

func (r *cartItemRepository) GetCartItemsByCartID(ctx context.Context, cartID string) ([]*domain.CartItem, error) {
	query := `SELECT id, cart_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE cart_id = $1`
	rows, err := r.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items by cart ID: %w", err)
	}
	defer rows.Close()

	var cartItems []*domain.CartItem
	for rows.Next() {
		cartItem := &domain.CartItem{}
		err := rows.Scan(
			&cartItem.ID,
			&cartItem.CartID,
			&cartItem.ProductID,
			&cartItem.Quantity,
			&cartItem.CreatedAt,
			&cartItem.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

func (r *cartItemRepository) UpdateCartItem(ctx context.Context, cartItem *domain.CartItem) error {
	query := `UPDATE cart_items SET quantity = $2, updated_at = $3 WHERE id = $1`
	cartItem.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, query, cartItem.ID, cartItem.Quantity, cartItem.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}
	return nil
}

func (r *cartItemRepository) DeleteCartItem(ctx context.Context, cartItemID string) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, cartItemID)
	if err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}
	return nil
}

func (r *cartItemRepository) DeleteCartItemByCartIDAndProductID(ctx context.Context, cartID, productID string) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	_, err := r.db.ExecContext(ctx, query, cartID, productID)
	if err != nil {
		return fmt.Errorf("failed to delete cart item by cart ID and product ID: %w", err)
	}
	return nil
}
