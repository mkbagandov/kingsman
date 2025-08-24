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
	query := `INSERT INTO users (phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code, loyalty_status, current_points) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.PhoneNumber, user.Email, user.PasswordHash, user.SocialID, user.DiscountLevel, user.ProgressToNextLevel, user.QRCode, user.LoyaltyStatus, user.CurrentPoints).Scan(&user.ID)
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
	query := `SELECT id, phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code, loyalty_status, current_points FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.SocialID, &user.DiscountLevel, &user.ProgressToNextLevel, &user.QRCode, &user.LoyaltyStatus, &user.CurrentPoints)
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
	query := `SELECT id, phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code, loyalty_status, current_points FROM users WHERE phone_number = $1`
	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.SocialID, &user.DiscountLevel, &user.ProgressToNextLevel, &user.QRCode, &user.LoyaltyStatus, &user.CurrentPoints)
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
	query := `SELECT id, phone_number, email, password_hash, social_id, discount_level, progress_to_next_level, qr_code, loyalty_status, current_points FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.PasswordHash, &user.SocialID, &user.DiscountLevel, &user.ProgressToNextLevel, &user.QRCode, &user.LoyaltyStatus, &user.CurrentPoints)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET phone_number = $2, email = $3, password_hash = $4, social_id = $5, discount_level = $6, progress_to_next_level = $7, qr_code = $8, loyalty_status = $9, current_points = $10 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.PhoneNumber, user.Email, user.PasswordHash, user.SocialID, user.DiscountLevel, user.ProgressToNextLevel, user.QRCode, user.LoyaltyStatus, user.CurrentPoints)
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

// Loyalty Program Implementations

func (r *PostgreSQLUserRepository) CreateLoyaltyPoint(ctx context.Context, point *domain.LoyaltyPoint) error {
	query := `INSERT INTO loyalty_points (user_id, points, type) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, point.UserID, point.Points, point.Type).Scan(&point.ID, &point.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create loyalty point: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) GetLoyaltyPointsByUserID(ctx context.Context, userID int) ([]*domain.LoyaltyPoint, error) {
	query := `SELECT id, user_id, points, type, created_at FROM loyalty_points WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loyalty points by user ID: %w", err)
	}
	defer rows.Close()

	var points []*domain.LoyaltyPoint
	for rows.Next() {
		point := &domain.LoyaltyPoint{}
		if err := rows.Scan(&point.ID, &point.UserID, &point.Points, &point.Type, &point.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan loyalty point: %w", err)
		}
		points = append(points, point)
	}
	return points, nil
}

func (r *PostgreSQLUserRepository) GetLoyaltyTierByID(ctx context.Context, id int) (*domain.LoyaltyTier, error) {
	tier := &domain.LoyaltyTier{}
	query := `SELECT id, name, min_points, description, benefits FROM loyalty_tiers WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&tier.ID, &tier.Name, &tier.MinPoints, &tier.Description, &tier.Benefits)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loyalty tier not found")
		}
		return nil, fmt.Errorf("failed to get loyalty tier by ID: %w", err)
	}
	return tier, nil
}

func (r *PostgreSQLUserRepository) GetLoyaltyTierByName(ctx context.Context, name string) (*domain.LoyaltyTier, error) {
	tier := &domain.LoyaltyTier{}
	query := `SELECT id, name, min_points, description, benefits FROM loyalty_tiers WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&tier.ID, &tier.Name, &tier.MinPoints, &tier.Description, &tier.Benefits)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loyalty tier not found")
		}
		return nil, fmt.Errorf("failed to get loyalty tier by name: %w", err)
	}
	return tier, nil
}

func (r *PostgreSQLUserRepository) GetAllLoyaltyTiers(ctx context.Context) ([]*domain.LoyaltyTier, error) {
	query := `SELECT id, name, min_points, description, benefits FROM loyalty_tiers ORDER BY min_points ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all loyalty tiers: %w", err)
	}
	defer rows.Close()

	var tiers []*domain.LoyaltyTier
	for rows.Next() {
		tier := &domain.LoyaltyTier{}
		if err := rows.Scan(&tier.ID, &tier.Name, &tier.MinPoints, &tier.Description, &tier.Benefits); err != nil {
			return nil, fmt.Errorf("failed to scan loyalty tier: %w", err)
		}
		tiers = append(tiers, tier)
	}
	return tiers, nil
}

func (r *PostgreSQLUserRepository) CreateLoyaltyTier(ctx context.Context, tier *domain.LoyaltyTier) error {
	query := `INSERT INTO loyalty_tiers (name, min_points, description, benefits) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, tier.Name, tier.MinPoints, tier.Description, tier.Benefits).Scan(&tier.ID)
	if err != nil {
		return fmt.Errorf("failed to create loyalty tier: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) UpdateLoyaltyTier(ctx context.Context, tier *domain.LoyaltyTier) error {
	query := `UPDATE loyalty_tiers SET name = $2, min_points = $3, description = $4, benefits = $5 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tier.ID, tier.Name, tier.MinPoints, tier.Description, tier.Benefits)
	if err != nil {
		return fmt.Errorf("failed to update loyalty tier: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) DeleteLoyaltyTier(ctx context.Context, id int) error {
	query := `DELETE FROM loyalty_tiers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete loyalty tier: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) CreateLoyaltyActivity(ctx context.Context, activity *domain.LoyaltyActivity) error {
	query := `INSERT INTO loyalty_activities (user_id, type, description) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, activity.UserID, activity.Type, activity.Description).Scan(&activity.ID, &activity.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create loyalty activity: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) GetLoyaltyActivitiesByUserID(ctx context.Context, userID int) ([]*domain.LoyaltyActivity, error) {
	query := `SELECT id, user_id, type, description, created_at FROM loyalty_activities WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loyalty activities by user ID: %w", err)
	}
	defer rows.Close()

	var activities []*domain.LoyaltyActivity
	for rows.Next() {
		activity := &domain.LoyaltyActivity{}
		if err := rows.Scan(&activity.ID, &activity.UserID, &activity.Type, &activity.Description, &activity.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan loyalty activity: %w", err)
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (r *PostgreSQLUserRepository) GetUserLoyalty(ctx context.Context, userID int) (*domain.UserLoyalty, error) {
	userLoyalty := &domain.UserLoyalty{}
	query := `SELECT user_id, current_points, current_tier_id, last_activity_at FROM user_loyalty WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&userLoyalty.UserID, &userLoyalty.CurrentPoints, &userLoyalty.CurrentTierID, &userLoyalty.LastActivityAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user loyalty not found")
		}
		return nil, fmt.Errorf("failed to get user loyalty: %w", err)
	}
	return userLoyalty, nil
}

func (r *PostgreSQLUserRepository) UpdateUserLoyalty(ctx context.Context, userLoyalty *domain.UserLoyalty) error {
	query := `INSERT INTO user_loyalty (user_id, current_points, current_tier_id, last_activity_at) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET current_points = EXCLUDED.current_points, current_tier_id = EXCLUDED.current_tier_id, last_activity_at = EXCLUDED.last_activity_at`
	_, err := r.db.ExecContext(ctx, query, userLoyalty.UserID, userLoyalty.CurrentPoints, userLoyalty.CurrentTierID, userLoyalty.LastActivityAt)
	if err != nil {
		return fmt.Errorf("failed to update user loyalty: %w", err)
	}
	return nil
}
