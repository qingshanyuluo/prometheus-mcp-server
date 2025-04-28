package prometheus

import (
	"github.com/mark3labs/mcp-go/server"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func NewServer(v1api v1.API) *server.MCPServer {
	// Create a new MCP server
	s := server.NewMCPServer(
		"github-mcp-server",
		"0.0.1",
		server.WithResourceCapabilities(true, true),
		server.WithLogging())

	// Add GitHub tools
	s.AddTool(searchMetrics(v1api))
	s.AddTool(getMetricLabels(v1api))
	s.AddTool(getMetricLabelValues(v1api))
	s.AddTool(query(v1api))
	s.AddTool(queryRange(v1api))

	return s
}
