package command

import (
	"fmt"

	"github.com/kiritoxkiriko/comical-tool/cli/internal/client"
	"github.com/kiritoxkiriko/comical-tool/cli/internal/config"
	"github.com/spf13/cobra"
)

type appState struct {
	cfg config.Config
}

func NewRoot() *cobra.Command {
	state := &appState{cfg: config.Default()}
	cmd := &cobra.Command{
		Use:   "comical-cli",
		Short: "CLI for comical-tool",
	}
	cmd.PersistentFlags().StringVar(&state.cfg.BaseURL, "base-url", state.cfg.BaseURL, "API base URL")
	cmd.PersistentFlags().StringVar(&state.cfg.AdminToken, "token", state.cfg.AdminToken, "admin token")
	cmd.PersistentFlags().StringVarP(&state.cfg.Output, "output", "o", state.cfg.Output, "output format")
	cmd.AddCommand(newConfigCommand(state), newShortCommand(state), newClipCommand(state))
	cmd.AddCommand(newImageCommand(state), newFileCommand(state), newAdminCommand(state))
	return cmd
}

func (s *appState) client() *client.Client {
	return client.New(s.cfg.BaseURL, s.cfg.AdminToken)
}

func printJSON(data []byte) error {
	_, err := fmt.Println(string(data))
	return err
}
