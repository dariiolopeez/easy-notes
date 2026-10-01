# easy-notes

A small CLI tool for managing local notes in your terminal with your favorite editor. Notes are stored as plain Markdown files with YAML Frontmatter — no database, no cloud, no lock-in.

Works seamlessly with standard Unix tools such as [`grep`](https://www.gnu.org/software/grep/), [`rg`](https://github.com/BurntSushi/ripgrep), [`ag`](https://github.com/ggreer/the_silver_searcher), [`fzf`](https://github.com/junegunn/fzf), [`peco`](https://github.com/peco/peco), and any editor that can be started from the command line.

## Table of Contents

- [Installation](#installation)
- [Basic Usage](#basic-usage)
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

## Installation

Install from source with the Go toolchain (1.24 or later):

```sh
go install github.com/dariiolopeez/easy-notes/cmd/easy-notes@latest
```

Or clone and build manually:

```sh
git clone https://github.com/dariiolopeez/easy-notes.git
cd easy-notes
go build -o easy-notes ./cmd/easy-notes
mv easy-notes ~/.local/bin/   # or any directory in $PATH
```

To uninstall:

```sh
rm -rf "$(easy-notes config | grep base_dir | awk '{print $2}')"   # remove all notes
rm "$(which easy-notes)"                                            # remove the binary
```

---

## Basic Usage

`easy-notes` provides subcommands to manage your Markdown notes.

- **Create** a new note with `easy-notes new <title> [--category <cat>] [--tags <t1,t2>]`.
- **List** all notes with `easy-notes list`. Filter by category or tag with `-c` and `-t`.
- **View** a note rendered in the terminal with `easy-notes view <id>`.
- **Edit** a note in your preferred editor with `easy-notes edit <id>`.

Notes are stored as `.md` files under an [XDG-compliant](https://wiki.archlinux.org/index.php/XDG_Base_Directory) home directory:

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

The home directory can be overridden with the `$EASY_NOTES_HOME` environment variable.

---

## Usage

### Create a new note

```sh
easy-notes new "Go Cobra tips" --category dev --tags go,cli,cobra
```

This creates a file at `~/.local/share/easy-notes/dev/7f3a1b2c_go-cobra-tips.md` and opens it in your editor immediately.

Every note starts with a YAML Frontmatter header generated automatically:

```markdown
---
id: 7f3a1b2c-e29b-41d4-a716-446655440000
title: Go Cobra tips
tags:
    - go
    - cli
    - cobra
created_at: 2026-10-01T20:00:00Z
updated_at: 2026-10-01T20:00:00Z
---

Write your note body here.
```

> **Do not remove the frontmatter block.** `easy-notes` uses it to identify and index notes. You may edit all fields freely.

The file name is derived from the first 8 characters of the UUID and a slug of the title: `<id8>_<slug>.md`. If you rename the title from the editor, `easy-notes edit` detects the change and renames the file automatically.

If no category is specified, the note goes to `inbox/`.

```sh
easy-notes new "Quick thought"          # → inbox/
easy-notes new "Sprint retro" -c work   # → work/
```

For more details: `easy-notes new --help`

---

### List notes

```sh
easy-notes list                   # list all notes
easy-notes list --category dev    # filter by category
easy-notes list --tag go          # filter by tag
```

Example output:

```
ID        TITLE                   CATEGORY   TAGS            UPDATED
────────  ──────────────────────  ─────────  ──────────────  ───────────────
7f3a1b2c  Go Cobra tips           dev        go, cli, cobra  2026-10-01 20:00
550e8400  What to buy this week   inbox      shopping        2026-10-01 18:00
f0e1d2c3  Trip to Lisbon          personal                   2026-09-30 09:15
```

Because `easy-notes list` prints file paths, you can combine it with standard Unix tools:

```sh
# Print full file paths instead (pipe-friendly)
easy-notes list -c dev | xargs grep "cobra"

# Open all notes matching a keyword in Vim
vim $(easy-notes list | xargs grep -l "cobra")

# Interactively pick a note with fzf and open it
easy-notes list | fzf | xargs $EDITOR
```

For more details: `easy-notes list --help`

---

### View a note

Renders the note directly in your terminal with syntax highlighting via [glamour](https://github.com/charmbracelet/glamour). Accepts the first 8 characters of the UUID or the full ID.

```sh
easy-notes view 7f3a1b2c
```

The output includes a styled metadata table (ID, category, tags, timestamps) followed by the rendered Markdown body. Automatically adapts to light and dark terminal themes.

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

Editor commands with flags are supported — `code --wait`, `subl -w`, `nvim`, etc.

---

### Manage categories

Categories are plain subdirectories. Create one explicitly:

```sh
easy-notes category add work
```

Or just use `--category` when creating a note — the directory is created automatically.

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

The configuration file lives at `~/.config/easy-notes/config.yaml`:

```yaml
base_dir: ~/.local/share/easy-notes
editor: nvim
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `$EASY_NOTES_HOME` | `~/.local/share/easy-notes` | Home directory for all notes |
| `$EDITOR` | — | Editor fallback when not set in config |
| `$VISUAL` | — | Secondary editor fallback |

---

## Unix Workflows

`easy-notes` is designed to compose with standard Unix tools.

**Search across all notes with ripgrep:**

```sh
rg "cobra" $(easy-notes list)
```

**Interactively pick and open a note with fzf:**

```sh
easy-notes list | fzf --preview 'cat {}' | xargs -o $EDITOR
```

**Open the most recently updated note:**

```sh
$EDITOR $(easy-notes list | head -1)
```

**Delete all notes in a category:**

```sh
rm $(easy-notes list -c inbox)
```

**Find notes by tag and view the first result:**

```sh
easy-notes list -t go | head -1 | xargs easy-notes view
```

---

## FAQ

**Can I change the notes home directory?**

Set `base_dir` in `~/.config/easy-notes/config.yaml`, or export `$EASY_NOTES_HOME` before running:

```sh
export EASY_NOTES_HOME=/path/to/my/notes
```

**Can I use VS Code or another GUI editor?**

Yes. Set your editor with the `--wait` flag so the terminal blocks until you close the file:

```yaml
# ~/.config/easy-notes/config.yaml
editor: code --wait
```

Or with Sublime Text:

```yaml
editor: subl -w
```

**How do I search inside notes?**

Combine `easy-notes list` with any grep-compatible tool:

```sh
grep -r "kubernetes" ~/.local/share/easy-notes/
rg "kubernetes" $(easy-notes list)
ag "kubernetes" $(easy-notes list -c dev)
```

**How do I interactively filter and open notes?**

```sh
easy-notes list | fzf | xargs -o $EDITOR
```

**How do I open the last modified note?**

```sh
$EDITOR $(ls -t ~/.local/share/easy-notes/**/*.md | head -1)
```

**What happens if I rename the title inside the editor?**

`easy-notes edit` detects the title change when you exit the editor and renames the file automatically to `<id8>_<new-slug>.md`. The UUID and all other metadata are preserved.

**Can I version-control my notes with Git?**

Yes — the notes directory is a standard directory of Markdown files. Initialize a Git repo there:

```sh
cd ~/.local/share/easy-notes
git init
git add .
git commit -m "initial notes snapshot"
```

Add a remote and push to back up your notes anywhere.

---

## License

[MIT](LICENSE)
