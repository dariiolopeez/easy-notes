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
	Short: "Gestor de notas local basado en archivos Markdown",
	Long: `easy-notes es un gestor de notas local basado en archivos Markdown
plano y convenciones Unix. Cada nota es un archivo .md con YAML Frontmatter.`,
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
		fmt.Fprintf(os.Stderr, "error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	store, err = storage.New(cfg.BaseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error inicializando almacenamiento: %v\n", err)
		os.Exit(1)
	}
}
