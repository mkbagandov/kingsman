package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
)

type PostgreSQLNotificationRepository struct {
	db *sql.DB
}

func NewPostgreSQLNotificationRepository(db *sql.DB) *PostgreSQLNotificationRepository {
	return &PostgreSQLNotificationRepository{db: db}
}

func (r *PostgreSQLNotificationRepository) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	query := `INSERT INTO notifications (id, user_id, type, title, message, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, notification.ID, notification.UserID, notification.Type, notification.Title, notification.Message, notification.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (r *PostgreSQLNotificationRepository) GetNotificationsByUserID(ctx context.Context, userID string) ([]*domain.Notification, error) {
	query := `SELECT id, user_id, type, title, message, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		notification := &domain.Notification{}
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &notification.Title, &notification.Message, &notification.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return notifications, nil
}
