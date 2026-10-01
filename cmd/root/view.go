package root

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "Visualiza una nota directamente en el terminal",
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func runView(cmd *cobra.Command, args []string) error {
	n, err := store.FindByID(args[0])
	if err != nil {
		return err
	}

	divider := strings.Repeat("─", 52)
	shortID := n.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	fmt.Println(divider)
	fmt.Printf("  %s\n", n.Title)
	fmt.Println(divider)
	fmt.Printf("  ID         %s\n", shortID)
	fmt.Printf("  Categoría  %s\n", n.Category)
	if len(n.Tags) > 0 {
		fmt.Printf("  Tags       %s\n", strings.Join(n.Tags, ", "))
	}
	fmt.Printf("  Creado     %s\n", n.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("  Actualizado %s\n", n.UpdatedAt.Format("2006-01-02 15:04"))
	fmt.Println(divider)
	fmt.Println()
	fmt.Println(n.Body)
	fmt.Println()
	fmt.Println(divider)
	return nil
}
