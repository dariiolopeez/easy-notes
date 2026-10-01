package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage note categories",
}

var categoryAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new category",
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
		return fmt.Errorf("failed to create category '%s': %w", name, err)
	}
	fmt.Printf("category '%s' created at %s\n", name, cfg.BaseDir)
	return nil
}
