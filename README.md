# Prometheus MCP 服务器

这是一个简化版的Prometheus MCP服务器，用于收集和暴露MCP服务器的指标。

## 功能特性

ai 友好化分析指标

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