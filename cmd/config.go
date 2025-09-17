package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long:  `Configuration management commands for Comical Tool.`,
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Show the current configuration values.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current configuration:")
		fmt.Println("=====================")

		// Server config
		fmt.Printf("Server Host: %s\n", viper.GetString("server.host"))
		fmt.Printf("Server Port: %d\n", viper.GetInt("server.port"))

		// Database config
		fmt.Printf("Database Host: %s\n", viper.GetString("database.host"))
		fmt.Printf("Database Port: %d\n", viper.GetInt("database.port"))
		fmt.Printf("Database User: %s\n", viper.GetString("database.user"))
		fmt.Printf("Database Name: %s\n", viper.GetString("database.db_name"))

		// Redis config
		fmt.Printf("Redis Host: %s\n", viper.GetString("redis.host"))
		fmt.Printf("Redis Port: %d\n", viper.GetInt("redis.port"))
		fmt.Printf("Redis DB: %d\n", viper.GetInt("redis.db"))

		// Short URL config
		fmt.Printf("Short URL Domain: %s\n", viper.GetString("short_url.domain"))
		fmt.Printf("Code Length: %d\n", viper.GetInt("short_url.code_length"))
		fmt.Printf("Default Expiry Hours: %d\n", viper.GetInt("short_url.default_expiry_hours"))
		fmt.Printf("Analytics Retention Days: %d\n", viper.GetInt("short_url.analytics_retention_days"))
	},
}

// configInitCmd represents the config init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Initialize a new configuration file with default values.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile := "config.yaml"
		if viper.ConfigFileUsed() != "" {
			configFile = viper.ConfigFileUsed()
		}

		// Check if config file already exists
		if _, err := os.Stat(configFile); err == nil {
			fmt.Printf("Configuration file %s already exists.\n", configFile)
			fmt.Println("Use --force flag to overwrite it.")
			return
		}

		// Write default config
		err := viper.WriteConfigAs(configFile)
		if err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			return
		}

		fmt.Printf("Configuration file created: %s\n", configFile)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)

	configInitCmd.Flags().Bool("force", false, "Overwrite existing config file")
}
