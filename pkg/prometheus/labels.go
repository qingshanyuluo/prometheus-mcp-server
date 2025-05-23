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

func getMetricLabels(v1api v1.API) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_metric_labels",
			mcp.WithDescription("获取指定指标名称下的所有标签"),
			mcp.WithString("name", mcp.Description("指标名称，例如：'http_requests_total'")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'name' 不是字符串类型")
			}

			// 构建查询表达式
			matcher := fmt.Sprintf("{__name__=\"%s\"}", name)

			// 查询标签列表
			labels, _, err := v1api.LabelNames(ctx, []string{matcher}, time.Now().Add(-time.Hour), time.Now())
			if err != nil {
				return nil, fmt.Errorf("查询标签失败: %w", err)
			}

			// 过滤掉__name__标签
			filteredLabels := make([]string, 0, len(labels))
			for _, label := range labels {
				if label != "__name__" {
					filteredLabels = append(filteredLabels, label)
				}
			}

			// 将结果转换为 JSON
			r, err := json.Marshal(filteredLabels)
			if err != nil {
				return nil, fmt.Errorf("序列化结果失败: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func getMetricLabelValues(v1api v1.API) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_metric_label_values",
			mcp.WithDescription("获取指定指标名称下某个标签的所有值"),
			mcp.WithString("name", mcp.Description("指标名称，例如：'http_requests_total'")),
			mcp.WithString("label", mcp.Description("标签名称，例如：'method'")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, ok := request.Params.Arguments["name"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'name' 不是字符串类型")
			}

			label, ok := request.Params.Arguments["label"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'label' 不是字符串类型")
			}

			// 构建查询表达式
			matcher := fmt.Sprintf("{__name__=\"%s\"}", name)

			// 查询标签值
			values, _, err := v1api.LabelValues(ctx, label, []string{matcher}, time.Now().Add(-time.Hour), time.Now())
			if err != nil {
				return nil, fmt.Errorf("查询标签值失败: %w", err)
			}

			// 将结果转换为字符串切片
			labelValues := make([]string, len(values))
			for i, value := range values {
				labelValues[i] = string(value)
			}

			// 将结果转换为 JSON
			r, err := json.Marshal(labelValues)
			if err != nil {
				return nil, fmt.Errorf("序列化结果失败: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
