package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAdminCommand(_ *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "admin", Short: "Administrative commands"}
	cmd.AddCommand(&cobra.Command{
		Use:   "cleanup",
		Short: "Request cleanup of expired resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), `{"cleanup":"not_implemented"}`)
			return err
		},
	})
	return cmd
}
