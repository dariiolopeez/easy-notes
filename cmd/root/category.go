package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

var categoryCmd = &cobra.Command{
	Use:   "cat",
	Short: "Manage note categories",
	Long:  `Manage note categories, which are stored as subdirectories under base_dir.`,
}

var categoryAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new category",
	Long:  `Create a new category subdirectory inside your configured base_dir.`,
	Example: `  easy-notes cat add work
  easy-notes cat add personal`,
	Args: cobra.ExactArgs(1),
	RunE: runCategoryAdd,
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
