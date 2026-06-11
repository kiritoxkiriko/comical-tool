package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConfigCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "Manage CLI config"}
	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Print default environment config",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(),
				"COMICAL_API_BASE_URL=%s\nCOMICAL_ADMIN_TOKEN=%s\n",
				state.cfg.BaseURL, state.cfg.AdminToken)
			return err
		},
	})
	return cmd
}
