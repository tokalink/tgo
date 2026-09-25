package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tokalink/tgo/pkg/database"
)

var (
	tenantDataDir string
)

var TenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Manage multi-tenant provisioning and schemas",
}

var TenantCreateCmd = &cobra.Command{
	Use:     "create <slug>",
	Aliases: []string{"tenant:create"},
	Short:   "Provision a new tenant database schema and run migrations",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]
		log.Printf("[Craft] Provisioning new tenant: %q...\n", slug)

		if err := database.ValidateTenantSlug(slug); err != nil {
			log.Fatalf("[Craft] Error: %v", err)
		}

		ctx := context.Background()

		// Default to SQLite file-per-tenant isolator
		isolator, err := database.NewSQLiteTenantIsolator(tenantDataDir)
		if err != nil {
			log.Fatalf("[Craft] Failed to initialize tenant isolator: %v", err)
		}
		defer isolator.Close()

		// 1. Provision tenant
		if err := isolator.ProvisionTenant(ctx, slug); err != nil {
			log.Fatalf("[Craft] Failed to provision tenant %q: %v", slug, err)
		}

		// 2. Run migrations
		migrator := database.NewMigrator(nil)
		if err := migrator.MigrateTenant(ctx, isolator, slug); err != nil {
			log.Fatalf("[Craft] Failed to run migrations for tenant %q: %v", slug, err)
		}

		fmt.Println("--------------------------------------------------")
		log.Printf("[Craft] Tenant %q successfully provisioned and migrated!\n", slug)
		fmt.Println("--------------------------------------------------")
	},
}

// TenantCreateRootCmd allows executing `craft tenant:create <slug>` directly from root command
var TenantCreateRootCmd = &cobra.Command{
	Use:    "tenant:create <slug>",
	Short:  "Provision a new tenant database schema and run migrations",
	Hidden: false,
	Args:   cobra.ExactArgs(1),
	Run:    TenantCreateCmd.Run,
}

func init() {
	TenantCreateCmd.Flags().StringVar(&tenantDataDir, "data-dir", "./data/tenants", "Directory where tenant databases are stored (for SQLite)")
	TenantCreateRootCmd.Flags().StringVar(&tenantDataDir, "data-dir", "./data/tenants", "Directory where tenant databases are stored (for SQLite)")
	TenantCmd.AddCommand(TenantCreateCmd)
}
