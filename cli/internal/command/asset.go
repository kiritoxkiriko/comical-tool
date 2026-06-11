package command

import (
	"net/http"

	"github.com/spf13/cobra"
)

func newImageCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "image", Short: "Manage image hosting"}
	cmd.AddCommand(newAssetUpload(state, "upload", "/api/images"))
	cmd.AddCommand(newAssetList(state, "list", "/api/images"))
	cmd.AddCommand(newAssetDelete(state, "delete", "/api/images/"))
	return cmd
}

func newFileCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "file", Short: "Manage temporary files"}
	cmd.AddCommand(newAssetUpload(state, "upload", "/api/files"))
	cmd.AddCommand(newAssetList(state, "list", "/api/files"))
	cmd.AddCommand(newAssetDelete(state, "delete", "/api/files/"))
	return cmd
}

func newAssetUpload(state *appState, use string, path string) *cobra.Command {
	var filePath, ttl string
	var link bool
	cmd := &cobra.Command{
		Use:   use,
		Short: "Upload an object",
		RunE: func(cmd *cobra.Command, args []string) error {
			values := map[string]string{"ttl": ttl}
			if link {
				values["link"] = "true"
			}
			data, err := state.client().Upload(path, filePath, values)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "file path")
	cmd.Flags().StringVar(&ttl, "ttl", "", "TTL duration")
	cmd.Flags().BoolVar(&link, "link", false, "create short link")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newAssetList(state *appState, use string, path string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: "List objects",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := state.client().Get(path, nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
}

func newAssetDelete(state *appState, use string, path string) *cobra.Command {
	return &cobra.Command{
		Use:   use + " <id>",
		Short: "Delete an object",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := state.client().JSON(http.MethodDelete, path+args[0], nil)
			if err != nil {
				return err
			}
			return printJSON(data)
		},
	}
}
