package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func query(v1api v1.API) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("query",
			mcp.WithDescription("执行PromQL即时查询"),
			mcp.WithString("query", mcp.Description("PromQL查询语句")),
			mcp.WithString("time", mcp.Description("查询时间戳，格式为RFC3339或Unix时间戳，可选")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			queryStr, ok := request.Params.Arguments["query"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'query' 不是字符串类型")
			}

			var queryTime time.Time
			if timeStr, ok := request.Params.Arguments["time"].(string); ok {
				var err error
				queryTime, err = time.Parse(time.RFC3339, timeStr)
				if err != nil {
					return nil, fmt.Errorf("解析时间参数失败: %w", err)
				}
			} else {
				queryTime = time.Now()
			}

			result, _, err := v1api.Query(ctx, queryStr, queryTime)
			if err != nil {
				return nil, fmt.Errorf("查询失败: %w", err)
			}

			r, err := json.Marshal(result)
			if err != nil {
				return nil, fmt.Errorf("序列化结果失败: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func queryRange(v1api v1.API) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("query_range",
			mcp.WithDescription("执行PromQL范围查询"),
			mcp.WithString("query", mcp.Description("PromQL查询语句")),
			mcp.WithString("start", mcp.Description("开始时间，格式为RFC3339或Unix时间戳(毫秒)")),
			mcp.WithString("end", mcp.Description("结束时间，格式为RFC3339或Unix时间戳(毫秒)")),
			mcp.WithString("step", mcp.Description("查询步长，例如'15s'、'1m'")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			queryStr, ok := request.Params.Arguments["query"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'query' 不是字符串类型")
			}

			startStr, ok := request.Params.Arguments["start"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'start' 不是字符串类型")
			}
			startTime, err := parseTime(startStr)
			if err != nil {
				return nil, fmt.Errorf("解析开始时间失败: %w", err)
			}

			endStr, ok := request.Params.Arguments["end"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'end' 不是字符串类型")
			}
			endTime, err := parseTime(endStr)
			if err != nil {
				return nil, fmt.Errorf("解析结束时间失败: %w", err)
			}

			stepStr, ok := request.Params.Arguments["step"].(string)
			if !ok {
				return nil, fmt.Errorf("参数 'step' 不是字符串类型")
			}
			stepDuration, err := time.ParseDuration(stepStr)
			if err != nil {
				return nil, fmt.Errorf("解析步长失败: %w", err)
			}

			r := v1.Range{
				Start: startTime,
				End:   endTime,
				Step:  stepDuration,
			}

			result, _, err := v1api.QueryRange(ctx, queryStr, r)
			if err != nil {
				return nil, fmt.Errorf("范围查询失败: %w", err)
			}

			// 限制返回的序列数不超过5条
			if result != nil && result.Type() == model.ValMatrix {
				matrix := result.(model.Matrix)
				if len(matrix) > 5 {
					result = matrix[:5]
				}
			}

			jsonResult, err := json.Marshal(result)
			if err != nil {
				return nil, fmt.Errorf("序列化结果失败: %w", err)
			}

			return mcp.NewToolResultText(string(jsonResult)), nil
		}
}
