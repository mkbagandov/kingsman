package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mkbagandov/kingsman/backend/app/internal/domain"
)

type PostgreSQLStoreRepository struct {
	db *sql.DB
}

func NewPostgreSQLStoreRepository(db *sql.DB) *PostgreSQLStoreRepository {
	return &PostgreSQLStoreRepository{db: db}
}

func (r *PostgreSQLStoreRepository) GetStores(ctx context.Context) ([]*domain.Store, error) {
	query := `SELECT id, name, address, location, phone FROM stores`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get stores: %w", err)
	}
	defer rows.Close()

	var stores []*domain.Store
	for rows.Next() {
		store := &domain.Store{}
		if err := rows.Scan(&store.ID, &store.Name, &store.Address, &store.Location, &store.Phone); err != nil {
			return nil, fmt.Errorf("failed to scan store: %w", err)
		}
		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return stores, nil
}

func (r *PostgreSQLStoreRepository) GetStoreByID(ctx context.Context, id string) (*domain.Store, error) {
	store := &domain.Store{}
	query := `SELECT id, name, address, location, phone FROM stores WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&store.ID, &store.Name, &store.Address, &store.Location, &store.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("store not found")
		}
		return nil, fmt.Errorf("failed to get store by ID: %w", err)
	}
	return store, nil
}
