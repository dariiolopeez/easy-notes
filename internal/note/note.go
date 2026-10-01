package note

import (
	"bytes"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	ID        string    `yaml:"id"`
	Title     string    `yaml:"title"`
	Tags      []string  `yaml:"tags"`
	CreatedAt time.Time `yaml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at"`
}

type Note struct {
	Frontmatter
	Category string
	Body     string
	FilePath string
}

func Slug(title string) string {
	var b strings.Builder
	prev := false
	for _, r := range strings.ToLower(title) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prev = false
		} else if !prev {
			b.WriteRune('-')
			prev = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func Parse(data []byte, category, filePath string) (*Note, error) {
	s := string(data)

	if !strings.HasPrefix(s, "---\n") {
		return nil, nil
	}

	parts := strings.SplitN(s[4:], "\n---", 2)
	if len(parts) != 2 {
		return nil, nil
	}

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(parts[0]), &fm); err != nil {
		return nil, err
	}

	body := strings.TrimPrefix(parts[1], "\n")
	body = strings.TrimSpace(body)

	return &Note{
		Frontmatter: fm,
		Category:    category,
		Body:        body,
		FilePath:    filePath,
	}, nil
}

func (n *Note) Marshal() ([]byte, error) {
	fm, err := yaml.Marshal(&n.Frontmatter)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(fm)
	buf.WriteString("---\n")
	if n.Body != "" {
		buf.WriteString("\n")
		buf.WriteString(n.Body)
		buf.WriteString("\n")
	}

	return buf.Bytes(), nil
}
