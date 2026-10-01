package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dariiolopeez/easy-notes/internal/editor"
	"github.com/dariiolopeez/easy-notes/internal/note"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Open an existing note in the configured editor",
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

	originalTitle := n.Title
	originalPath := n.FilePath

	if err := editor.Open(cfg.ResolveEditor(), originalPath); err != nil {
		return fmt.Errorf("editor error: %w", err)
	}

	data, err := os.ReadFile(originalPath)
	if err != nil {
		return err
	}

	updated, err := note.Parse(data, n.Category, originalPath)
	if err != nil || updated == nil {
		return err
	}

	updated.UpdatedAt = note.Now()

	if updated.Title != originalTitle {
		shortID := updated.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		newFilename := fmt.Sprintf("%s_%s.md", shortID, note.Slug(updated.Title))
		newPath := filepath.Join(cfg.BaseDir, updated.Category, newFilename)
		if err := os.Rename(originalPath, newPath); err != nil {
			return fmt.Errorf("failed to rename note file: %w", err)
		}
		updated.FilePath = newPath
	}

	return store.Save(updated)
}
