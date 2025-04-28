# Prometheus MCP 服务器

这是一个简化版的Prometheus MCP服务器，用于收集和暴露MCP服务器的指标。

## 功能特性

### AI友好化分析指标

1. **JSON格式输出**：所有查询结果均以JSON格式返回，便于AI系统解析处理

### MCP工具

1. **search_metrics** - 按正则表达式搜索指标
   - 参数: `pattern` (指标名称模式，支持正则表达式)
   - 示例: `http_.*` 匹配所有以`http_`开头的指标
   - 返回: 匹配的指标名称列表(JSON格式)

2. **get_metric_labels** - 获取指标的所有标签
   - 参数: `name` (指标名称)
   - 示例: `http_requests_total`
   - 返回: 该指标的所有标签列表(JSON格式，不包括`__name__`标签)

3. **get_metric_label_values** - 获取指标标签的所有值
   - 参数: 
     - `name` (指标名称)
     - `label` (标签名称)
   - 示例: 获取`http_requests_total`指标的`method`标签的所有值
   - 返回: 指定标签的所有值列表(JSON格式)

### 服务器端

1. 设置Prometheus地址：
```bash
export PROMETHEUS_URL=your_Prometheus_endpoint
```
2. 启动服务器：
```bash
go run cmd/server/main.go sse
```

服务器默认将在`:8081`端口暴露Prometheus指标。

### mcp 客户端配置

配置示例：
```json
{
  "mcpServers": {
    "products-sse": {
      "url": "http://localhost:8081/sse"
    }
  }
}
```
## 项目结构

```
.
├── cmd/           # 命令行程序入口
│   └── server/   # 服务器代码
├── internal/      # 内部实现
├── pkg/           # 可复用包
│   └── prometheus # Prometheus相关实现
└── go.mod         # Go模块定义
```