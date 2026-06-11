package command

import (
	"net/http"

	"github.com/spf13/cobra"
)

func newAdminCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "admin", Short: "Administrative commands"}
	cmd.AddCommand(&cobra.Command{
		Use:   "cleanup",
		Short: "Request cleanup of expired resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := state.client().JSON(http.MethodPost, "/api/admin/cleanup", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	})
	return cmd
}
