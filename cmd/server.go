package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	comical_tool "github.com/kiritoxkiriko/comical-tool/biz/handler/comical_tool"
	"github.com/kiritoxkiriko/comical-tool/biz/router"
	"github.com/kiritoxkiriko/comical-tool/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Comical Tool server",
	Long: `Start the Comical Tool server with all services including:
- Short URL service
- Analytics tracking
- Redis caching
- MySQL database

The server will listen on the configured host and port.`,
	Run: runServer,
}

var (
	host string
	port int
)

func init() {
	rootCmd.AddCommand(serverCmd)

	// Server-specific flags
	serverCmd.Flags().StringVar(&host, "host", "0.0.0.0", "server host")
	serverCmd.Flags().IntVar(&port, "port", 8080, "server port")
	serverCmd.Flags().Bool("dev", false, "development mode")
	serverCmd.Flags().Bool("debug", false, "debug mode")

	// Bind flags to viper
	viper.BindPFlag("server.host", serverCmd.Flags().Lookup("host"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}

func runServer(cmd *cobra.Command, args []string) {
	log.Println("Starting Comical Tool server...")

	// Initialize services
	err := comical_tool.InitServices()
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Create Hertz server
	h := server.Default(
		server.WithHostPorts(viper.GetString("server.host")+":"+viper.GetString("server.port")),
		server.WithMaxRequestBodySize(10*1024*1024), // 10MB
	)

	// Register routes - we need to call the register function from main package
	// For now, we'll use the generated router directly
	router.GeneratedRegister(h)

	// Start cleanup routine for old analytics
	go startCleanupRoutine()

	// Start server in a goroutine
	go func() {
		h.Spin()
	}()

	log.Printf("Server started on %s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	log.Println("Press Ctrl+C to stop the server")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// startCleanupRoutine starts a background routine to clean up old analytics data
func startCleanupRoutine() {
	analyticsService := services.NewAnalyticsService()
	ticker := time.NewTicker(24 * time.Hour) // Run daily
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		retentionDays := viper.GetInt("short_url.analytics_retention_days")
		if retentionDays == 0 {
			retentionDays = 30 // default
		}
		err := analyticsService.CleanupOldAnalytics(ctx, retentionDays)
		cancel()

		if err != nil {
			hlog.Errorf("Failed to cleanup old analytics: %v", err)
		} else {
			hlog.Info("Successfully cleaned up old analytics data")
		}
	}
}
