# Prometheus MCP 服务器

这是一个简化版的Prometheus MCP服务器，用于收集和暴露MCP服务器的指标。

## 功能特性

- 简化的指标收集：仅保留`mcp_requests_total`核心指标
- 自动化的指标注册：内置Prometheus指标注册表
- 轻量级设计：移除了非核心功能组件

## 快速开始

### 服务器端

1. 设置GitHub访问令牌：
```bash
export PROMETHEUS_URL=your_Prometheus_endpoint
```
2. 启动服务器：
```bash
go run cmd/server/main.go sse
```

服务器默认将在`:8081`端口暴露Prometheus指标。

### SSE服务器配置

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

使用SSE端点获取实时指标数据流。


## 指标说明

| 指标名称 | 类型 | 描述 |
|----------|------|------|
| mcp_requests_total | Counter | MCP请求总数，包含`tool`和`status`标签 |

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