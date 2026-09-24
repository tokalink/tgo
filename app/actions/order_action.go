package actions

import (
	"context"
	"fmt"

	"github.com/tgo-framework/tgo/pkg/database"
)

// OrderAction encapsulates domain business logic
type OrderAction struct {
	db database.DBEngine
}

func NewOrderAction(db database.DBEngine) *OrderAction {
	return &OrderAction{db: db}
}

func (a *OrderAction) Execute(ctx context.Context, id string) error {
	scopedDB, err := a.db.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("resolving tenant db: %w", err)
	}
	_ = scopedDB
	return nil
}
