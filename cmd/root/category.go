package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Gestiona categorías de notas",
}

var categoryAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Crea una nueva categoría",
	Args:  cobra.ExactArgs(1),
	RunE:  runCategoryAdd,
}

func init() {
	rootCmd.AddCommand(categoryCmd)
	categoryCmd.AddCommand(categoryAddCmd)
}

func runCategoryAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	if err := store.EnsureCategory(name); err != nil {
		return fmt.Errorf("error creando categoría '%s': %w", name, err)
	}
	fmt.Printf("categoría '%s' creada en %s\n", name, cfg.BaseDir)
	return nil
}
