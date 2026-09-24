package models

import (
	"context"
	"time"

	"github.com/tgo-framework/tgo/pkg/database"
)

// Product represents the entity model
type Product struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProductRepository handles database queries for Product
type ProductRepository struct {
	db database.DBEngine
}

func NewProductRepository(db database.DBEngine) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*Product, error) {
	scopedDB, err := r.db.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	row := scopedDB.Conn().QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM products WHERE id = ?", id)
	var m Product
	if err := row.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}
