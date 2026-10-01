# easy-notes

A small CLI tool for managing local notes in your terminal with your favorite editor. Notes are stored as plain Markdown files with YAML Frontmatter — no database, no cloud, no lock-in.

Works seamlessly with standard Unix tools such as [`grep`](https://www.gnu.org/software/grep/), [`rg`](https://github.com/BurntSushi/ripgrep), [`ag`](https://github.com/ggreer/the_silver_searcher), [`fzf`](https://github.com/junegunn/fzf), [`peco`](https://github.com/peco/peco), and any editor that can be started from the command line.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Create a new note](#create-a-new-note)
  - [List notes](#list-notes)
  - [View a note](#view-a-note)
  - [Edit a note](#edit-a-note)
  - [Manage categories](#manage-categories)
  - [Configuration](#configuration)
- [Environment Variables](#environment-variables)
- [Unix Workflows](#unix-workflows)
- [FAQ](#faq)
- [License](#license)

---

## Requirements

- **Go 1.24 or later** — required to build and install
- **A terminal editor** — `nano`, `vim`, `nvim`, or any editor launchable from the command line
- **Linux or macOS** — Windows is not tested

---

## Installation

**From source with Go:**

```sh
go install github.com/dariiolopeez/easy-notes/cmd/easy-notes@latest
```

The binary is placed in `~/go/bin/`. Make sure that directory is in your `$PATH`:

```sh
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc   # or ~/.zshrc
source ~/.bashrc
```

**Clone and build with Make:**

```sh
git clone https://github.com/dariiolopeez/easy-notes.git
cd easy-notes
make build          # compiles to bin/easy-notes
make install        # runs go install for global access
```

**Available Make targets:**

```
make build     compile binary to bin/easy-notes
make install   go install for global access
make test      run all tests with go test -v ./...
make clean     remove bin/ and test cache
```

**To uninstall:**

```sh
rm -rf ~/.local/share/easy-notes    # remove all notes
rm "$(which easy-notes)"            # remove the binary
```

---

## Quick Start

```sh
# 1. Create your first note (opens in your editor)
easy-notes new "My first note"

# 2. List all notes
easy-notes list

# 3. View a note in the terminal (use the 8-char ID from list)
easy-notes view 550e8400

# 4. Edit a note
easy-notes edit 550e8400

# 5. Check active configuration
easy-notes config
```

---

## Usage

### Create a new note

```sh
easy-notes new "Go Cobra tips" --category dev --tags go,cli,cobra
```

This creates a file at `~/.local/share/easy-notes/dev/7f3a1b2c_go-cobra-tips.md` and opens it in your editor immediately. The editor is detected automatically (see [Edit a note](#edit-a-note)).

Every note starts with a YAML Frontmatter header generated automatically:

```markdown
---
id: 7f3a1b2c-e29b-41d4-a716-446655440000
title: Go Cobra tips
tags:
    - go
    - cli
    - cobra
created_at: 2026-10-01 20:00
updated_at: 2026-10-01 20:00
---

Write your note body here in Markdown.
```

> **Do not remove the frontmatter block.** `easy-notes` uses it to identify and index notes. You may edit all fields freely.

The file name is derived from the first 8 characters of the UUID and a slug of the title: `<id8>_<slug>.md`. If you rename the `title` field from the editor, `easy-notes edit` detects the change and renames the file automatically.

If no category is specified, the note goes to `inbox/`.

```sh
easy-notes new "Quick thought"          # → inbox/
easy-notes new "Sprint retro" -c work   # → work/
easy-notes new "Docker tips" -c dev -t docker,containers
```

```
Flags:
  -c, --category string   Destination category (default "inbox")
  -t, --tags string       Comma-separated tags
```

---

### List notes

```sh
easy-notes list                    # list all notes
easy-notes list --category dev     # filter by category
easy-notes list --tag go           # filter by tag
```

Example output:

```
ID        TITLE                   CATEGORY   TAGS            UPDATED
──        ─────                   ────────   ────            ───────
7f3a1b2c  Go Cobra tips           dev        go, cli, cobra  2026-10-01 20:00
550e8400  What to buy this week   inbox      shopping        2026-10-01 18:00
f0e1d2c3  Trip to Lisbon          personal                   2026-09-30 09:15
```

```
Flags:
  -c, --category string   Filter by category
  -t, --tag string        Filter by tag
```

---

### View a note

Renders the note directly in your terminal with syntax highlighting via [glamour](https://github.com/charmbracelet/glamour). The output includes a styled metadata table followed by the rendered Markdown body. Automatically adapts to light and dark terminal themes.

Accepts the first 8 characters of the UUID or the full ID:

```sh
easy-notes view 7f3a1b2c
easy-notes view 7f3a1b2c-e29b-41d4-a716-446655440000
```

---

### Edit a note

```sh
easy-notes edit 7f3a1b2c
```

Opens the note in the configured editor. When you exit:

- `updated_at` in the frontmatter is updated automatically.
- If you changed the `title` field, the file is renamed to match the new slug (`<id8>_<new-slug>.md`).

The editor is resolved in this order:

| Priority | Source |
|---|---|
| 1 | `editor` field in `~/.config/easy-notes/config.yaml` |
| 2 | `$EDITOR` environment variable |
| 3 | `$VISUAL` environment variable |
| 4 | `nano` if available in `$PATH` |
| 5 | `vim` as last resort |

Editor commands with flags are fully supported — `code --wait`, `subl -w`, `nvim`, etc.

---

### Manage categories

Categories are plain subdirectories under the notes home. Create one explicitly:

```sh
easy-notes category add work
```

Or use `--category` when creating a note — the directory is created automatically if it does not exist.

Notes directory layout:

```
~/.local/share/easy-notes/
├── inbox/
│   └── 550e8400_what-to-buy-this-week.md
├── dev/
│   ├── 7f3a1b2c_go-cobra-tips.md
│   └── a1b2c3d4_dockerize-the-app.md
└── personal/
    └── f0e1d2c3_trip-to-lisbon.md
```

---

### Configuration

```sh
easy-notes config
```

Shows the active configuration:

```
base_dir  /home/you/.local/share/easy-notes
editor    nvim
```

The configuration file lives at `~/.config/easy-notes/config.yaml` and is created with defaults if it does not exist:

```yaml
base_dir: /home/you/.local/share/easy-notes
editor: nvim
```

To change the notes directory, edit `base_dir`. To set a preferred editor, edit `editor`. Both fields are optional — if omitted, defaults apply.

---

## Environment Variables

| Variable | Description |
|---|---|
| `$EDITOR` | Editor used when `editor` is not set in config |
| `$VISUAL` | Secondary editor fallback after `$EDITOR` |

> The notes home directory is configured via `base_dir` in `config.yaml`, not via an environment variable.

---

## Unix Workflows

Since notes are plain `.md` files in a known directory, you can use standard Unix tools directly on them.

**Search across all notes with ripgrep:**

```sh
rg "kubernetes" ~/.local/share/easy-notes/
```

**Search and filter by category:**

```sh
rg "cobra" ~/.local/share/easy-notes/dev/
```

**Interactively pick a note with fzf and open it in your editor:**

```sh
find ~/.local/share/easy-notes -name "*.md" | fzf | xargs -o $EDITOR
```

**Interactively pick and view with fzf preview:**

```sh
find ~/.local/share/easy-notes -name "*.md" \
  | fzf --preview 'cat {}' \
  | xargs -I{} easy-notes view {}
```

**Open the most recently modified note:**

```sh
$EDITOR "$(find ~/.local/share/easy-notes -name '*.md' -printf '%T@ %p\n' | sort -rn | head -1 | cut -d' ' -f2-)"
```

**Delete all notes in a category:**

```sh
rm ~/.local/share/easy-notes/inbox/*.md
```

**Count notes per category:**

```sh
find ~/.local/share/easy-notes -name "*.md" | sed 's|/[^/]*$||' | sort | uniq -c
```

---

## FAQ

**How do I set my preferred editor?**

Edit `~/.config/easy-notes/config.yaml`:

```yaml
editor: nvim
```

Or rely on the `$EDITOR` environment variable:

```sh
export EDITOR=nvim
```

**Can I use VS Code or another GUI editor?**

Yes. Use the `--wait` flag so the terminal blocks until you close the file:

```yaml
# ~/.config/easy-notes/config.yaml
editor: code --wait
```

Sublime Text:

```yaml
editor: subl -w
```

**How do I change the notes home directory?**

Edit `base_dir` in `~/.config/easy-notes/config.yaml`:

```yaml
base_dir: /path/to/my/notes
```

**How do I search inside notes?**

Use `grep`, `rg`, or `ag` directly on the notes directory:

```sh
grep -r "kubernetes" ~/.local/share/easy-notes/
rg "kubernetes" ~/.local/share/easy-notes/dev/
ag "kubernetes" ~/.local/share/easy-notes/
```

**What happens if I rename the title inside the editor?**

`easy-notes edit` re-reads the frontmatter when you exit. If the `title` field changed, the file is renamed automatically to `<id8>_<new-slug>.md`. The UUID and all other metadata are preserved.

**Can I version-control my notes with Git?**

Yes — the notes directory is just a directory of plain Markdown files:

```sh
cd ~/.local/share/easy-notes
git init
git add .
git commit -m "initial notes snapshot"
git remote add origin git@github.com:you/my-notes.git
git push -u origin main
```

**How do I run the tests?**

```sh
make test
# or directly:
go test -v ./...
```

---

## License

[MIT](LICENSE)
