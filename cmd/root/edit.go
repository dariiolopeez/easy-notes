package root

import (
	"fmt"
	"os"
	"time"

	"github.com/dariiolopeez/easy-notes/internal/editor"
	"github.com/dariiolopeez/easy-notes/internal/note"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Abre una nota existente en el editor",
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	n, err := store.FindByID(args[0])
	if err != nil {
		return err
	}

	editorCmd := cfg.ResolveEditor()
	if err := editor.Open(editorCmd, n.FilePath); err != nil {
		return fmt.Errorf("error abriendo editor: %w", err)
	}

	data, err := os.ReadFile(n.FilePath)
	if err != nil {
		return err
	}

	updated, err := note.Parse(data, n.Category, n.FilePath)
	if err != nil || updated == nil {
		return err
	}

	updated.UpdatedAt = time.Now()
	return store.Save(updated)
}
