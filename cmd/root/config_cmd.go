package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "cfg",
	Short: "Show the active configuration",
	Long: `Display the active base_dir and the resolved editor from
~/.config/easy-notes/config.yaml or environment fallbacks ($EDITOR, $VISUAL).`,
	Example: `  easy-notes cfg`,
	RunE:    runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	fmt.Printf("base_dir  %s\n", cfg.BaseDir)
	fmt.Printf("editor    %s\n", cfg.ResolveEditor())
	return nil
}
