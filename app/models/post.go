package models

import (
	"context"
	"time"

	"github.com/tgo-framework/tgo/pkg/database"
)

// Post represents the entity model
type Post struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PostRepository handles database queries for Post
type PostRepository struct {
	db database.DBEngine
}

func NewPostRepository(db database.DBEngine) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) FindByID(ctx context.Context, id string) (*Post, error) {
	scopedDB, err := r.db.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	row := scopedDB.Conn().QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM posts WHERE id = ?", id)
	var m Post
	if err := row.Scan(&m.ID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}
