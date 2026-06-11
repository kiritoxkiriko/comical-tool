package command

import (
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

func newClipCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "clip", Short: "Manage temporary clipboard"}
	cmd.AddCommand(newClipPut(state), newClipGet(state), newClipDelete(state))
	return cmd
}

func newClipPut(state *appState) *cobra.Command {
	var content, password, ttl string
	var maxVisits int
	var link bool
	cmd := &cobra.Command{
		Use:   "put",
		Short: "Create a clipboard item",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{"content": content, "password": password, "ttl": ttl, "max_visits": maxVisits, "link": link}
			data, err := state.client().JSON(http.MethodPost, "/api/clip", body)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	cmd.Flags().StringVar(&content, "content", "", "clipboard content")
	cmd.Flags().StringVar(&password, "password", "", "read password")
	cmd.Flags().StringVar(&ttl, "ttl", "", "TTL duration")
	cmd.Flags().IntVar(&maxVisits, "max-visits", 0, "max visits")
	cmd.Flags().BoolVar(&link, "link", false, "create short link")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newClipGet(state *appState) *cobra.Command {
	var password string
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Read a clipboard item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := url.Values{"password": []string{password}}
			data, err := state.client().Get("/api/clip/"+args[0], query)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	cmd.Flags().StringVar(&password, "password", "", "read password")
	return cmd
}

func newClipDelete(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a clipboard item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := state.client().JSON(http.MethodDelete, "/api/clip/"+args[0], nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
}
