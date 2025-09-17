package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "1.0.0"
var buildTime = "unknown"
var gitCommit = "unknown"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Comical Tool",
	Long:  `Print the version number, build time, and git commit of Comical Tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Comical Tool version %s\n", version)
		fmt.Printf("Build time: %s\n", buildTime)
		fmt.Printf("Git commit: %s\n", gitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
