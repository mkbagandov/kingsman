package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain" // Update with your actual project path
)

type orderItemRepository struct {
	db *sql.DB
}

func NewOrderItemRepository(db *sql.DB) domain.OrderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) CreateOrderItem(ctx context.Context, orderItem *domain.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price, orderItem.CreatedAt, orderItem.UpdatedAt,
	).Scan(&orderItem.ID)

	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	return nil
}

func (r *orderItemRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]*domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, price, created_at, updated_at
		FROM order_items WHERE order_id = $1 ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items by order ID: %w", err)
	}
	defer rows.Close()

	var orderItems []*domain.OrderItem
	for rows.Next() {
		orderItem := &domain.OrderItem{}
		err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price, &orderItem.CreatedAt, &orderItem.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		orderItems = append(orderItems, orderItem)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return orderItems, nil
}

func (r *orderItemRepository) UpdateOrderItem(ctx context.Context, orderItem *domain.OrderItem) error {
	query := `
		UPDATE order_items
		SET quantity = $1, price = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query, orderItem.Quantity, orderItem.Price, orderItem.UpdatedAt, orderItem.ID)
	if err != nil {
		return fmt.Errorf("failed to update order item: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order item with ID %d not found", orderItem.ID)
	}
	return nil
}

func (r *orderItemRepository) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	query := `
		DELETE FROM order_items WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order item with ID %d not found", orderItemID)
	}
	return nil
}
