package database

import (
	"context"
	"testing"
)

func TestSQLiteAdapter(t *testing.T) {
	// Use in-memory database for testing
	adapter, err := NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to create SQLite adapter: %v", err)
	}
	defer adapter.Close()

	ctx := context.Background()

	// Test Ping
	if err := adapter.Ping(ctx); err != nil {
		t.Errorf("Failed to ping SQLite database: %v", err)
	}

	// Test basic CRUD
	db := adapter.Conn()
	_, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec("INSERT INTO test (name) VALUES (?)", "TGo")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	var name string
	err = db.QueryRow("SELECT name FROM test WHERE id = 1").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}

	if name != "TGo" {
		t.Errorf("Expected 'TGo', got '%s'", name)
	}
}
