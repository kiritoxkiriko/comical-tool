package command

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

func newImageCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "image", Short: "Manage image hosting"}
	cmd.AddCommand(newAssetUpload(state, "upload", "/api/images", false))
	cmd.AddCommand(newAssetList(state, "list", "/api/images"))
	cmd.AddCommand(newAssetDelete(state, "delete", "/api/images/"))
	return cmd
}

func newFileCommand(state *appState) *cobra.Command {
	cmd := &cobra.Command{Use: "file", Short: "Manage temporary files"}
	cmd.AddCommand(newAssetUpload(state, "upload", "/api/files", true))
	cmd.AddCommand(newAssetList(state, "list", "/api/files"))
	cmd.AddCommand(newFileDownload(state))
	cmd.AddCommand(newAssetDelete(state, "delete", "/api/files/"))
	return cmd
}

func newAssetUpload(state *appState, use string, path string, accessPolicy bool) *cobra.Command {
	var filePath, ttl, password string
	var maxVisits int
	var link bool
	cmd := &cobra.Command{
		Use:   use,
		Short: "Upload an object",
		RunE: func(cmd *cobra.Command, args []string) error {
			values := map[string]string{"ttl": ttl}
			if link {
				values["link"] = "true"
			}
			if accessPolicy {
				values["password"] = password
				if maxVisits > 0 {
					values["max_visits"] = strconv.Itoa(maxVisits)
				}
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
	if accessPolicy {
		cmd.Flags().StringVar(&password, "password", "", "download password")
		cmd.Flags().IntVar(&maxVisits, "max-visits", 0, "max successful downloads")
	}
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

func newFileDownload(state *appState) *cobra.Command {
	var output, password string
	cmd := &cobra.Command{
		Use:   "download <id>",
		Short: "Download a temporary file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := output
			if target == "" {
				target = filepath.Base(args[0])
			}
			query := url.Values{}
			if password != "" {
				query.Set("password", password)
			}
			path := "/api/assets/" + args[0]
			if encoded := query.Encode(); encoded != "" {
				path += "?" + encoded
			}
			return state.client().Download(path, target)
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	cmd.Flags().StringVar(&password, "password", "", "download password")
	return cmd
}
