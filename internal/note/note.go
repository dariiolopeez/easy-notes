package note

import (
	"bytes"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"
)

const timeFormat = "2006-01-02 15:04"

type Timestamp time.Time

func Now() Timestamp {
	return Timestamp(time.Now().Truncate(time.Minute))
}

func (t Timestamp) Format(layout string) string {
	return time.Time(t).Format(layout)
}

func (t Timestamp) Equal(other Timestamp) bool {
	return time.Time(t).Equal(time.Time(other))
}

func (t Timestamp) MarshalYAML() (interface{}, error) {
	return time.Time(t).Format(timeFormat), nil
}

func (t *Timestamp) UnmarshalYAML(value *yaml.Node) error {
	parsed, err := time.Parse(timeFormat, value.Value)
	if err != nil {
		return err
	}
	*t = Timestamp(parsed)
	return nil
}

type Frontmatter struct {
	ID        string    `yaml:"id"`
	Title     string    `yaml:"title"`
	Tags      []string  `yaml:"tags"`
	CreatedAt Timestamp `yaml:"created_at"`
	UpdatedAt Timestamp `yaml:"updated_at"`
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
