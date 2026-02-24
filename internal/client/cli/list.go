package cli

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List of all stored secrets",
	Long:  "Retrieves and displays a list of all user secrets stored in the system.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return list(os.Stdout)
	},
}

func list(out io.Writer) error {
	list, err := secretService.InfoList()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tType\tNname\tFileName\tCreated")
	for _, item := range list {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			item.ID,
			item.DataType,
			item.Name,
			getFileName(item.Meta),
			item.Created.Format("2006-01-02 15:04:05"),
		)
	}
	w.Flush()

	return nil
}

func getFileName(meta []dto.MetaData) string {
	for _, item := range meta {
		if item.Name == MetaFileName {
			return item.Value
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(listCmd)
}
