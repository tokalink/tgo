package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tokalink/tgo/cmd/craft/commands"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "craft",
		Short: "TGo Craft CLI",
		Long:  "Craft is the command line developer tool for TGo to generate code, run migrations, and manage tenants.",
	}

	// Direct root commands
	rootCmd.AddCommand(commands.NewCmd)
	rootCmd.AddCommand(commands.ServeCmd)
	rootCmd.AddCommand(commands.MigrateCmd)
	rootCmd.AddCommand(commands.TenantCmd)
	rootCmd.AddCommand(commands.TenantCreateRootCmd)
	rootCmd.AddCommand(commands.MakeCmd)

	// Aliases for artisan-style colon syntax (craft make:model, craft make:service, etc.)
	makeModelAlias := &cobra.Command{
		Use:   "make:model <Name>",
		Short: "Generate a new model struct and repository",
		Args:  cobra.ExactArgs(1),
		Run:   commands.MakeModelCmd.Run,
	}
	makeModelAlias.Flags().AddFlagSet(commands.MakeModelCmd.Flags())
	rootCmd.AddCommand(makeModelAlias)

	makeMigrationAlias := &cobra.Command{
		Use:   "make:migration <Name>",
		Short: "Generate a new SQL migration up and down pair",
		Args:  cobra.ExactArgs(1),
		Run:   commands.MakeMigrationCmd.Run,
	}
	rootCmd.AddCommand(makeMigrationAlias)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:service <Name>",
		Short: "Generate a new business logic service action",
		Args:  cobra.ExactArgs(1),
		Run:   commands.MakeServiceCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:handler <Name>",
		Short: "Generate a new RPC & HTTP controller handler",
		Args:  cobra.ExactArgs(1),
		Run:   commands.MakeHandlerCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:frontend",
		Short: "Generate unified frontend starter (HTML/CSS/JS + Wails ready)",
		Run:   commands.MakeFrontendCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:booster <Name>",
		Short: "Generate a declarative CRUDBooster-style admin controller",
		Args:  cobra.ExactArgs(1),
		Run:   commands.MakeBoosterCmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "make:crud <Name>",
		Short: "Alias for make:booster to scaffold instant CRUD modules",
		Args:  cobra.ExactArgs(1),
		Run:   commands.MakeBoosterCmd.Run,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
