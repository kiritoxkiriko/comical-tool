package command

import (
	"net/http"

	"github.com/spf13/cobra"
)

func newShortCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "short", Short: "Manage short links"}
	cmd.AddCommand(newShortCreate(state), newShortRevoke(state))
	return cmd
}

func newShortCreate(state *appState) *cobra.Command {
	var targetURL, slug, ttl string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a short link",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]string{"target_url": targetURL, "custom_slug": slug, "ttl": ttl}
			data, err := state.client().JSON(http.MethodPost, "/api/short-links", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	cmd.Flags().StringVar(&targetURL, "url", "", "target URL")
	cmd.Flags().StringVar(&slug, "slug", "", "custom slug")
	cmd.Flags().StringVar(&ttl, "ttl", "", "TTL duration")
	_ = cmd.MarkFlagRequired("url")
	return cmd
}

func newShortRevoke(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "revoke <slug>",
		Short: "Revoke a short link",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := state.client().JSON(http.MethodPost, "/api/short-links/"+args[0]+"/revoke", nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
}
