package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
)

type PostgreSQLProductRepository struct {
	db *sql.DB
}

func NewPostgreSQLProductRepository(db *sql.DB) *PostgreSQLProductRepository {
	return &PostgreSQLProductRepository{db: db}
}

func (r *PostgreSQLProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (id, name, description, category_id, price, quantity, image_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	product.CreatedAt = time.Now().Format(time.RFC3339)
	product.UpdatedAt = time.Now().Format(time.RFC3339)
	err := r.db.QueryRowContext(ctx, query, product.ID, product.Name, product.Description, product.CategoryID, product.Price, product.Quantity, product.ImageURL, product.CreatedAt, product.UpdatedAt).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *PostgreSQLProductRepository) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	product := &domain.Product{}
	query := `SELECT id, name, description, category_id, price, quantity, image_url, created_at, updated_at FROM products WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price, &product.Quantity, &product.ImageURL, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}
	return product, nil
}

func (r *PostgreSQLProductRepository) GetProducts(ctx context.Context, categoryID *string, minPrice *float64, maxPrice *float64, sortBy *string, sortOrder *string, limit, offset int) ([]*domain.Product, error) {
	baseQuery := `SELECT id, name, description, category_id, price, quantity, image_url, created_at, updated_at FROM products`
	conditions := []string{`1=1`}
	args := []interface{}{}
	argCounter := 1

	if categoryID != nil && *categoryID != "" {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argCounter))
		args = append(args, *categoryID)
		argCounter++
	}
	if minPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argCounter))
		args = append(args, *minPrice)
		argCounter++
	}
	if maxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argCounter))
		args = append(args, *maxPrice)
		argCounter++
	}

	whereClause := " WHERE " + strings.Join(conditions, " AND ")

	orderByClause := ""
	if sortBy != nil && *sortBy != "" {
		order := "ASC"
		if sortOrder != nil && *sortOrder == "desc" {
			order = "DESC"
		}
		orderByClause = fmt.Sprintf(" ORDER BY %s %s", *sortBy, order)
	}

	paginationClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, limit, offset)

	query := baseQuery + whereClause + orderByClause + paginationClause
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price, &product.Quantity, &product.ImageURL, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over product rows: %w", err)
	}

	return products, nil
}

func (r *PostgreSQLProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	product.UpdatedAt = time.Now().Format(time.RFC3339)
	query := `UPDATE products SET name = $2, description = $3, category_id = $4, price = $5, quantity = $6, image_url = $7, updated_at = $8 WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Description, product.CategoryID, product.Price, product.Quantity, product.ImageURL, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *PostgreSQLProductRepository) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
