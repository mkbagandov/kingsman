package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain" // Update with your actual project path
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (user_id, total_amount, status, payment_status, order_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query, order.UserID, order.TotalAmount, order.Status, order.PaymentStatus, order.OrderDate, order.CreatedAt, order.UpdatedAt,
	).Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, orderID int) (*domain.Order, error) {
	query := `
		SELECT id, user_id, order_date, total_amount, status, payment_status, created_at, updated_at
		FROM orders WHERE id = $1
	`
	order := &domain.Order{}
	err := r.db.QueryRowContext(
		ctx, query, orderID,
	).Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalAmount, &order.Status, &order.PaymentStatus, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}
	return order, nil
}

func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID string, paymentStatus *string) ([]*domain.Order, error) {
	baseQuery := `
		SELECT id, user_id, order_date, total_amount, status, payment_status, created_at, updated_at
		FROM orders WHERE user_id = $1
	`
	args := []interface{}{userID}

	if paymentStatus != nil && *paymentStatus != "" {
		baseQuery += ` AND payment_status = $2`
		args = append(args, *paymentStatus)
	}

	baseQuery += ` ORDER BY order_date DESC`

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user ID: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order := &domain.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalAmount, &order.Status, &order.PaymentStatus, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return orders, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders
		SET total_amount = $1, status = $2, payment_status = $3, updated_at = $4
		WHERE id = $5
	`
	result, err := r.db.ExecContext(ctx, query, order.TotalAmount, order.Status, order.PaymentStatus, order.UpdatedAt, order.ID)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", order.ID)
	}
	return nil
}
