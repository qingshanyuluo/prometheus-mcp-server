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
)

var (
	logFile       string
	prometheusURL string
	baseURL       string

	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Prometheus MCP Server",
		Long:  `A Prometheus MCP server that handles various tools and resources.`,
	}

	sseCmd = &cobra.Command{
		Use:   "sse",
		Short: "Start sse server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if prometheusURL == "" {
				fmt.Println("Error: --prometheus-url is required")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
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
	// Add global flags that will be shared by all commands
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Path to log file")
	rootCmd.PersistentFlags().StringVar(&prometheusURL, "prometheus-url", "", "Prometheus server URL")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "http://localhost:8081", "Base URL for the SSE server")

	// Add subcommands
	rootCmd.AddCommand(sseCmd)
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

	// 创建Prometheus客户端
	promClient, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		return fmt.Errorf("failed to create prometheus client: %w", err)
	}

	// 创建v1 API客户端
	v1api := v1.NewAPI(promClient)

	// Create server
	promServer := prometheus.NewServer(v1api)
	sseServer := server.NewSSEServer(promServer, server.WithBaseURL(baseURL))

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
