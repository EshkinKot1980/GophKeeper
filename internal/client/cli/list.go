package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List of all stored secrets",
	Long:  "Retrieves and displays a list of all user secret keys stored in the system.",
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := secretService.InfoList()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tNAME\tCREATED AT")
		for _, item := range list {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				item.ID,
				item.DataType,
				item.Name,
				item.Created.Format("2006-01-02 15:04:05"),
			)
		}
		w.Flush()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
