package commands

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/tokalink/tgo/pkg/database"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Craft: Starting migration process...")
		
		// In a real application, configuration should be loaded here
		// Using SQLite memory database for demonstration purposes
		adapter, err := database.NewSQLite(":memory:")
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer adapter.Close()
		
		migrator := database.NewMigrator(adapter)
		
		if err := migrator.Up(context.Background()); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		
		log.Println("Craft: Migration completed successfully!")
	},
}
