package note_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dariiolopeez/easy-notes/internal/note"
)

func TestSlug(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"hello-world", "hello-world"},
		{"Go 1.22 Features!!!", "go-1-22-features"},
		{"  leading and trailing  ", "leading-and-trailing"},
		{"Hello  double  spaces", "hello-double-spaces"},
		{"café au lait", "café-au-lait"},
		{"Hello, World!", "hello-world"},
		{"100% done", "100-done"},
		{"---", ""},
		{"", ""},
		{"!@#$%^&*()", ""},
		{"a", "a"},
		{"123", "123"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := note.Slug(tc.input)
			if got != tc.want {
				t.Errorf("Slug(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSlugNoLeadingOrTrailingDash(t *testing.T) {
	inputs := []string{"!hello", "hello!", "!hello!", "  spaces  "}
	for _, input := range inputs {
		got := note.Slug(input)
		if strings.HasPrefix(got, "-") || strings.HasSuffix(got, "-") {
			t.Errorf("Slug(%q) = %q: must not start or end with a dash", input, got)
		}
	}
}

func TestSlugConsecutiveSpecialChars(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"hello...world", "hello-world"},
		{"foo///bar", "foo-bar"},
		{"a  b  c", "a-b-c"},
	}
	for _, tc := range cases {
		got := note.Slug(tc.input)
		if got != tc.want {
			t.Errorf("Slug(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestMarshalParseRoundTrip(t *testing.T) {
	now := time.Date(2026, 10, 1, 20, 0, 0, 0, time.UTC)

	original := &note.Note{
		Frontmatter: note.Frontmatter{
			ID:        "550e8400-e29b-41d4-a716-446655440000",
			Title:     "Go Cobra Tips",
			Tags:      []string{"go", "cli", "cobra"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		Category: "dev",
		Body:     "Some **markdown** content.\n\nWith multiple paragraphs.",
		FilePath: "/tmp/test.md",
	}

	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Marshal() returned error: %v", err)
	}

	parsed, err := note.Parse(data, original.Category, original.FilePath)
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	if parsed == nil {
		t.Fatal("Parse() returned nil for valid data")
	}

	if parsed.ID != original.ID {
		t.Errorf("ID = %q, want %q", parsed.ID, original.ID)
	}
	if parsed.Title != original.Title {
		t.Errorf("Title = %q, want %q", parsed.Title, original.Title)
	}
	if parsed.Category != original.Category {
		t.Errorf("Category = %q, want %q", parsed.Category, original.Category)
	}
	if parsed.Body != original.Body {
		t.Errorf("Body = %q, want %q", parsed.Body, original.Body)
	}
	if !parsed.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", parsed.CreatedAt, original.CreatedAt)
	}
	if !parsed.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", parsed.UpdatedAt, original.UpdatedAt)
	}
	if len(parsed.Tags) != len(original.Tags) {
		t.Fatalf("Tags length = %d, want %d", len(parsed.Tags), len(original.Tags))
	}
	for i := range original.Tags {
		if parsed.Tags[i] != original.Tags[i] {
			t.Errorf("Tags[%d] = %q, want %q", i, parsed.Tags[i], original.Tags[i])
		}
	}
}

func TestMarshalEmptyBody(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			ID:    "abc-000",
			Title: "Empty Body",
		},
		Body: "",
	}

	data, err := n.Marshal()
	if err != nil {
		t.Fatalf("Marshal() returned error: %v", err)
	}

	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		t.Error("output must start with ---")
	}
	if !strings.Contains(content, "\n---\n") {
		t.Error("output must contain closing ---")
	}
}

func TestMarshalOutputStructure(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			ID:    "test-id",
			Title: "Structure Test",
			Tags:  []string{"a", "b"},
		},
		Body: "body text",
	}

	data, err := n.Marshal()
	if err != nil {
		t.Fatalf("Marshal() returned error: %v", err)
	}

	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		t.Error("output must start with ---")
	}
	if !strings.Contains(content, "title: Structure Test") {
		t.Error("output must contain title field")
	}
	if !strings.Contains(content, "body text") {
		t.Error("output must contain body")
	}
}

func TestParseInvalidInput(t *testing.T) {
	cases := []struct {
		name  string
		input []byte
	}{
		{"no frontmatter", []byte("Just plain text without frontmatter")},
		{"unclosed frontmatter", []byte("---\nid: abc\ntitle: test\n")},
		{"empty input", []byte{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := note.Parse(tc.input, "inbox", "/tmp/test.md")
			if err != nil {
				t.Logf("Parse() returned error (acceptable): %v", err)
				return
			}
			if n != nil {
				t.Errorf("expected nil note for invalid input %q, got non-nil", tc.name)
			}
		})
	}
}

func TestParsePreservesFilePath(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{ID: "abc", Title: "T"},
		Body:        "body",
	}
	data, _ := n.Marshal()

	const path = "/custom/path/abc_t.md"
	parsed, err := note.Parse(data, "cat", path)
	if err != nil || parsed == nil {
		t.Fatalf("Parse() failed: err=%v, note=%v", err, parsed)
	}
	if parsed.FilePath != path {
		t.Errorf("FilePath = %q, want %q", parsed.FilePath, path)
	}
}
