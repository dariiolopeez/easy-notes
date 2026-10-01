package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dariiolopeez/easy-notes/internal/note"
)

type Storage struct {
	BaseDir string
}

func New(baseDir string) (*Storage, error) {
	if err := os.MkdirAll(filepath.Join(baseDir, "inbox"), 0755); err != nil {
		return nil, err
	}
	return &Storage{BaseDir: baseDir}, nil
}

func (s *Storage) List(category, tag string) ([]*note.Note, error) {
	root := s.BaseDir
	if category != "" {
		root = filepath.Join(s.BaseDir, category)
	}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, nil
	}

	var notes []*note.Note

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(s.BaseDir, path)
		cat := filepath.Dir(rel)

		n, err := note.Parse(data, cat, path)
		if err != nil || n == nil {
			return nil
		}

		if tag != "" && !hasTag(n.Tags, tag) {
			return nil
		}

		notes = append(notes, n)
		return nil
	})

	return notes, err
}

func (s *Storage) FindByID(id string) (*note.Note, error) {
	notes, err := s.List("", "")
	if err != nil {
		return nil, err
	}

	for _, n := range notes {
		if strings.HasPrefix(n.ID, id) {
			return n, nil
		}
	}

	return nil, fmt.Errorf("nota no encontrada: %s", id)
}

func (s *Storage) Save(n *note.Note) error {
	dir := filepath.Join(s.BaseDir, n.Category)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := n.Marshal()
	if err != nil {
		return err
	}

	return os.WriteFile(n.FilePath, data, 0644)
}

func (s *Storage) EnsureCategory(name string) error {
	return os.MkdirAll(filepath.Join(s.BaseDir, name), 0755)
}

func hasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
