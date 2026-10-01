package root

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available notes",
	RunE:  runList,
}

var (
	listCategory string
	listTag      string
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listCategory, "category", "c", "", "Filter by category")
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter by tag")
}

func runList(cmd *cobra.Command, args []string) error {
	notes, err := store.List(listCategory, listTag)
	if err != nil {
		return err
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tTITLE\tCATEGORY\tTAGS\tUPDATED")
	fmt.Fprintln(w, "──\t─────\t────────\t────\t───────")
	for _, n := range notes {
		shortID := n.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		tags := strings.Join(n.Tags, ", ")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			shortID,
			n.Title,
			n.Category,
			tags,
			n.UpdatedAt.Format("2006-01-02 15:04"),
		)
	}
	return w.Flush()
}
