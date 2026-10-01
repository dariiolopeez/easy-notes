package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the active configuration",
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	fmt.Printf("base_dir  %s\n", cfg.BaseDir)
	fmt.Printf("editor    %s\n", cfg.ResolveEditor())
	return nil
}
