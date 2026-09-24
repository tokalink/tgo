package models

import (
	"context"
	"time"

	"github.com/tgo-framework/tgo/pkg/database"
)

// Test represents the entity model
type Test struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TestRepository handles database queries for Test
type TestRepository struct {
	db database.DBEngine
}

func NewTestRepository(db database.DBEngine) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) FindByID(ctx context.Context, id string) (*Test, error) {
	scopedDB, err := r.db.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	row := scopedDB.Conn().QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM tests WHERE id = ?", id)
	var m Test
	if err := row.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}
