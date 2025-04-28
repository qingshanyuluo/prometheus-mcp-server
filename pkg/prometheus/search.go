package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func searchMetrics(v1api v1.API) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("search_metrics",
			mcp.WithDescription("search metrics by name pattern"),
			mcp.WithString("pattern", mcp.Description("Metric name pattern, supports regex, e.g. 'http_.*' matches all metrics starting with 'http_'")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pattern, ok := request.Params.Arguments["pattern"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'pattern' 不是字符串类型")
			}
			println("pattern:", []string{pattern}[0])
			pattern = "{__name__=~\"" + pattern + "\"}"
			// 执行标签查询
			labelValues, _, err := v1api.LabelValues(ctx, "__name__", []string{pattern}, time.Now().Add(-time.Hour), time.Now())
			if err != nil {
				return nil, fmt.Errorf("查询指标失败: %w", err)
			}

			// 将结果转换为字符串切片
			metrics := make([]string, len(labelValues))
			for i, value := range labelValues {
				metrics[i] = string(value)
			}

			// 将结果转换为 JSON
			r, err := json.Marshal(metrics)
			if err != nil {
				return nil, fmt.Errorf("序列化结果失败: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
