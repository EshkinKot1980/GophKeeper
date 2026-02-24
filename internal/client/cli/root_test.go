package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd_Structure(t *testing.T) {
	if rootCmd.Use != "gophkeeper" {
		t.Errorf(`Expected command name: "gophkeeper", got: "%s"`, rootCmd.Use)
	}

	tests := []struct {
		name           string
		cmd            *cobra.Command
		wantSubcommand map[string]bool
	}{
		{
			name: "root_subcommands",
			cmd:  rootCmd,
			wantSubcommand: map[string]bool{
				"register": false,
				"login":    false,
				"add":      false,
				"get":      false,
				"list":     false,
			},
		}, {
			name: "add_subcommands",
			cmd:  addCmd,
			wantSubcommand: map[string]bool{
				"credentials": false,
				"file":        false,
				"text":        false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, cmd := range test.cmd.Commands() {
				if _, exists := test.wantSubcommand[cmd.Name()]; exists {
					test.wantSubcommand[cmd.Name()] = true
				} else {
					t.Errorf("Unexpected  %q  command registered in rootCmd", cmd.Name())
				}
			}

			for name, found := range test.wantSubcommand {
				if !found {
					t.Errorf("Command %q is not registered in rootCmd", name)
				}
			}
		})
	}

}
