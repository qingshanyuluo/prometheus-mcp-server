package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/qingshanyuluo/prometheus-mcp-server/pkg/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Prometheus MCP Server",
		Long:  `A Prometheus MCP server that handles various tools and resources.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Bind flag to viper
			viper.BindPFlag("log-file", cmd.PersistentFlags().Lookup("log-file"))
		},
	}

	sseCmd = &cobra.Command{
		Use:   "sse",
		Short: "Start sse server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(cmd *cobra.Command, args []string) {
			logFile := viper.GetString("log-file")
			logger, err := initLogger(logFile)
			if err != nil {
				stdlog.Fatal("Failed to initialize logger:", err)
			}
			if err := runSseServer(logger); err != nil {
				stdlog.Fatal("failed to run sse server:", err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Add global flags that will be shared by all commands
	rootCmd.PersistentFlags().String("log-file", "", "Path to log file")

	// Add subcommands
	rootCmd.AddCommand(sseCmd)
}

func initConfig() {
	// Initialize Viper configuration
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
}

func initLogger(outPath string) (*log.Logger, error) {
	if outPath == "" {
		return log.New(), nil
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New()
	logger.SetOutput(file)

	return logger, nil
}

func runSseServer(logger *log.Logger) error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create Prometheus client
	prometheusUrl := os.Getenv("PROMETHEUS_URL")
	if prometheusUrl == "" {
		logger.Fatal("PROMETHEUS_URL not set")
	}
	// 创建Prometheus客户端
	promClient, err := api.NewClient(api.Config{
		Address: prometheusUrl,
	})
	if err != nil {
		return fmt.Errorf("failed to create prometheus client: %w", err)
	}

	// 创建v1 API客户端
	v1api := v1.NewAPI(promClient)

	// Create server
	promServer := prometheus.NewServer(v1api)
	sseServer := server.NewSSEServer(promServer, server.WithBaseURL("http://localhost:8081"))

	if err := sseServer.Start(":8081"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		logger.Infof("shutting down server...")
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
