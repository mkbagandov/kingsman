package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
)

type PostgreSQLCategoryRepository struct {
	db *sql.DB
}

func NewPostgreSQLCategoryRepository(db *sql.DB) *PostgreSQLCategoryRepository {
	return &PostgreSQLCategoryRepository{db: db}
}

func (r *PostgreSQLCategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (id, name) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, category.ID, category.Name).Scan(&category.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return fmt.Errorf("category with name %s already exists", category.Name)
		}
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *PostgreSQLCategoryRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	category := &domain.Category{}
	query := `SELECT id, name FROM categories WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}
	return category, nil
}

func (r *PostgreSQLCategoryRepository) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		category := &domain.Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over category rows: %w", err)
	}

	return categories, nil
}

func (r *PostgreSQLCategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	query := `UPDATE categories SET name = $2 WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, category.ID, category.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return fmt.Errorf("category with name %s already exists", category.Name)
		}
		return fmt.Errorf("failed to update category: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

func (r *PostgreSQLCategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}
