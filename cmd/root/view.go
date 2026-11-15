package root

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "v <id>",
	Short: "Render a note in the terminal with syntax highlighting",
	Long: `Render a note's metadata and Markdown body directly in the terminal using
glamour. You can pass either the 8-character short ID or the full UUID.`,
	Example: `  easy-notes v 7f3a1b2c
  easy-notes v 7f3a1b2c-e29b-41d4-a716-446655440000`,
	Args: cobra.ExactArgs(1),
	RunE: runView,
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func runView(cmd *cobra.Command, args []string) error {
	n, err := store.FindByID(args[0])
	if err != nil {
		return err
	}

	shortID := n.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", n.Title)
	sb.WriteString("| Field | Value |\n|---|---|\n")
	fmt.Fprintf(&sb, "| ID | `%s` |\n", shortID)
	fmt.Fprintf(&sb, "| Category | `%s` |\n", n.Category)
	if len(n.Tags) > 0 {
		fmt.Fprintf(&sb, "| Tags | %s |\n", strings.Join(n.Tags, ", "))
	}
	fmt.Fprintf(&sb, "| Created | %s |\n", n.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Fprintf(&sb, "| Updated | %s |\n", n.UpdatedAt.Format("2006-01-02 15:04"))
	sb.WriteString("\n---\n\n")
	if n.Body != "" {
		sb.WriteString(n.Body)
		sb.WriteString("\n")
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return fmt.Errorf("renderer init error: %w", err)
	}

	out, err := renderer.Render(sb.String())
	if err != nil {
		return fmt.Errorf("render error: %w", err)
	}

	fmt.Print(out)
	return nil
}
