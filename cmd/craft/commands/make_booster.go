package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var MakeBoosterCmd = &cobra.Command{
	Use:   "make:booster <Name>",
	Short: "Generate a declarative CRUDBooster-style admin controller",
	Long:  "Scaffolds a declarative CRUD controller with auto grid columns, form fields, and lifecycle hooks.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.TrimSpace(args[0])
		if name == "" {
			fmt.Println("❌ Error: Module name cannot be empty.")
			os.Exit(1)
		}

		pascalName := strings.Title(name)
		snakeName := strings.ToLower(name)
		pluralSnake := snakeName
		if !strings.HasSuffix(pluralSnake, "s") {
			pluralSnake += "s"
		}

		targetDir := "app/controllers/admin"
		_ = os.MkdirAll(targetDir, 0755)

		filePath := filepath.Join(targetDir, fmt.Sprintf("%s_controller.go", snakeName))
		content := fmt.Sprintf(`package admin

import (
	"github.com/tokalink/tgo/pkg/booster"
)

// %sController manages CRUD for %s
type %sController struct {
	booster.Controller
}

func New%sController() *%sController {
	c := &%sController{}

	// 1. Module Configuration
	c.Title = "%s Management"
	c.Table = "%s"
	c.PrimaryKey = "id"
	c.OrderBy = "id DESC"

	// 2. Data Grid Columns
	c.Columns = []booster.Column{
		{Label: "ID", Name: "id", Type: booster.TypeText, Sortable: true},
		{Label: "Name", Name: "name", Type: booster.TypeText, Searchable: true, Sortable: true},
		{Label: "Status", Name: "status", Type: booster.TypeBadge},
	}

	// 3. Form Input Fields (Create & Edit)
	c.Forms = []booster.Field{
		{Label: "Name", Name: "name", Type: booster.InputText, Required: true, Placeholder: "Enter name..."},
		{Label: "Status", Name: "status", Type: booster.InputSelect, Required: true, Options: []booster.Option{
			{Value: "active", Label: "Active"},
			{Value: "inactive", Label: "Inactive"},
		}},
	}

	// 4. Lifecycle Hooks
	c.HookBeforeAdd = func(ctx *booster.Context, data map[string]interface{}) error {
		// Custom business logic before insert
		return nil
	}

	return c
}
`, pascalName, pluralSnake, pascalName, pascalName, pascalName, pascalName, pascalName, pluralSnake)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Printf("❌ Failed to write controller: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Generated Booster Controller: %s\n", filePath)
		fmt.Printf("\nMount to your TGo application container:\n")
		fmt.Printf("  admin := booster.NewEngine()\n")
		fmt.Printf("  admin.Register(adminCtrl.New%sController())\n", pascalName)
		fmt.Printf("  admin.Mount(application.Server(), \"/admin\")\n\n")
	},
}
