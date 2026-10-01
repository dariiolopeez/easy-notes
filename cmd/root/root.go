package root

import (
	"fmt"
	"os"

	"github.com/dariiolopeez/easy-notes/internal/config"
	"github.com/dariiolopeez/easy-notes/internal/storage"
	"github.com/spf13/cobra"
)

var cfg *config.Config
var store *storage.Storage

var rootCmd = &cobra.Command{
	Use:   "easy-notes",
	Short: "Local note manager based on plain Markdown files",
	Long: `easy-notes is a local note manager based on plain Markdown files
and Unix conventions. Each note is a .md file with YAML Frontmatter.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDeps)
}

func initDeps() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	store, err = storage.New(cfg.BaseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
}
