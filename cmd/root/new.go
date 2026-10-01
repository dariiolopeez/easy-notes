package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dariiolopeez/easy-notes/internal/editor"
	"github.com/dariiolopeez/easy-notes/internal/note"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Crea una nueva nota y la abre en el editor",
	Args:  cobra.ExactArgs(1),
	RunE:  runNew,
}

var (
	newCategory string
	newTags     string
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&newCategory, "category", "c", "inbox", "Categoría destino")
	newCmd.Flags().StringVarP(&newTags, "tags", "t", "", "Tags separados por coma")
}

func runNew(cmd *cobra.Command, args []string) error {
	title := args[0]
	id := uuid.New().String()
	slug := note.Slug(title)
	filename := fmt.Sprintf("%s_%s.md", id[:8], slug)

	var tags []string
	if newTags != "" {
		for _, t := range strings.Split(newTags, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	now := time.Now()
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			ID:        id,
			Title:     title,
			Tags:      tags,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Category: newCategory,
		Body:     "",
		FilePath: filepath.Join(cfg.BaseDir, newCategory, filename),
	}

	if err := store.Save(n); err != nil {
		return fmt.Errorf("error guardando nota: %w", err)
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
		return nil
	}
	updated.UpdatedAt = time.Now()
	if err := store.Save(updated); err != nil {
		return err
	}

	fmt.Printf("nota guardada: %s\n", n.FilePath)
	return nil
}
