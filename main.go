package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tokalink/tgo/cmd/craft/commands"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "craft",
		Short: "TGo Craft CLI",
		Long:  "Craft is the command line developer tool for TGo to generate code, run migrations, manage tenants, and serve the application.",
		Run: func(cmd *cobra.Command, args []string) {
			// Default action when executed without subcommands is to start the server
			commands.ServeCmd.Run(cmd, args)
		},
	}

	// Subcommands
	rootCmd.AddCommand(commands.ServeCmd)
	rootCmd.AddCommand(commands.MigrateCmd)
	rootCmd.AddCommand(commands.TenantCmd)
	rootCmd.AddCommand(commands.TenantCreateRootCmd)
	rootCmd.AddCommand(commands.MakeCmd)

	// Aliases for artisan-style colon syntax
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

	// Also support "craft" as first subcommand if someone runs "go run . craft <command>"
	rootCmd.AddCommand(&cobra.Command{
		Use:                "craft",
		Short:              "Craft CLI alias",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				commands.ServeCmd.Run(cmd, args)
				return
			}
			rootCmd.SetArgs(args)
			_ = rootCmd.Execute()
		},
	})

	// Handle case where user executes: go run . craft <subcommand>
	if len(os.Args) > 1 && os.Args[1] == "craft" {
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
