package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dariiolopeez/easy-notes/internal/editor"
	"github.com/dariiolopeez/easy-notes/internal/note"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <category> <title>",
	Short: "Create a new note in a category and open it in the editor",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runNew,
}

var (
	newTags string
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&newTags, "tags", "t", "", "Comma-separated tags")
}

func runNew(cmd *cobra.Command, args []string) error {
	category := args[0]
	title := strings.Join(args[1:], " ")
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

	now := note.Now()
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			ID:        id,
			Title:     title,
			Tags:      tags,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Category: category,
		Body:     "",
		FilePath: filepath.Join(cfg.BaseDir, category, filename),
	}

	if err := store.Save(n); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	editorCmd := cfg.ResolveEditor()
	if err := editor.Open(editorCmd, n.FilePath); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	data, err := os.ReadFile(n.FilePath)
	if err != nil {
		return err
	}
	updated, err := note.Parse(data, n.Category, n.FilePath)
	if err != nil || updated == nil {
		return nil
	}
	updated.UpdatedAt = note.Now()
	if err := store.Save(updated); err != nil {
		return err
	}

	fmt.Printf("note saved: %s\n", n.FilePath)
	return nil
}
