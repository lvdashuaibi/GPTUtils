# GPTUtils 项目总结

## 📌 项目概述

GPTUtils 是一个功能完整、易于使用的通义千问大模型 API Go SDK。它提供了稳定可靠的 HTTP 客户端实现，支持基础对话、流式输出、多轮对话等功能，可以轻松集成到其他 Go 项目中。

## ✨ 核心特性

### 已实现功能

- ✅ **基础对话** - 单轮和多轮对话支持
- ✅ **流式输出** - 实时流式返回响应内容
- ✅ **参数配置** - 温度、Top-P、最大Token等参数调整
- ✅ **多模型支持** - qwen-plus、qwen-turbo、qwen-max 等
- ✅ **安全配置** - API Key 环境变量管理
- ✅ **原生HTTP实现** - 稳定可靠的 HTTP 客户端
- ✅ **交互式命令行** - 功能完整的对话工具
- ✅ **SDK 导出接口** - 易于集成到其他项目

### 规划功能

- 🚧 **工具调用** - Function Calling 功能
- 🚧 **联网搜索** - 联网搜索增强回答
- 🚧 **图像处理** - 支持图像输入和分析
- 🚧 **音频处理** - 支持音频输入和输出

## 📁 项目结构

```
GPTUtils/
├── client/                      # 核心客户端实现
│   ├── http_client.go          # HTTP 客户端（主要实现）
│   ├── types.go                # 数据类型定义
│   ├── client.go               # OpenAI SDK 客户端（备选）
│   ├── stream.go               # 流式输出支持
│   ├── tools.go                # 工具调用（规划中）
│   └── search.go               # 联网搜索（规划中）
├── config/                      # 配置管理
│   └── config.go               # 配置类和工厂函数
├── cmd/                         # 命令行工具
│   └── chat/main.go            # 交互式对话工具
├── examples/                    # 示例程序
│   ├── http_simple_chat.go     # 基础对话示例
│   ├── http_stream_chat.go     # 流式输出示例
│   ├── http_multi_turn.go      # 多轮对话示例
│   └── sdk_integration_example.go # SDK 集成示例
├── sdk.go                       # SDK 公共接口导出
├── README.md                    # 项目文档
├── SDK_USAGE.md                 # SDK 详细使用指南
├── INTEGRATION_GUIDE.md         # 集成指南
├── PROJECT_SUMMARY.md           # 本文件
├── .env.example                 # 环境变量示例
├── .gitignore                   # Git 忽略文件
├── go.mod                       # Go 模块定义
├── go.sum                       # 依赖校验和
└── VERSION                      # 版本号
```

## 🚀 快速开始

### 1. 安装

```bash
# 克隆项目
git clone <repository-url>
cd GPTUtils

# 下载依赖
go mod download
```

### 2. 配置

```bash
# 设置 API Key
export API_KEY="your-api-key-here"

# 或创建 .env 文件
cp .env.example .env
# 编辑 .env 文件，填入 API Key
```

### 3. 使用

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
)

func main() {
    client := gptutils.NewDefaultClient()
    response, _ := client.SimpleChat(context.Background(), "你好")
    fmt.Println(response)
}
```

## 📚 文档

| 文档 | 说明 |
|------|------|
| [README.md](./README.md) | 项目主文档，包含功能介绍和基本使用 |
| [SDK_USAGE.md](./SDK_USAGE.md) | SDK 详细使用指南，包含完整 API 参考 |
| [INTEGRATION_GUIDE.md](./INTEGRATION_GUIDE.md) | 集成指南，包含多个集成场景示例 |
| [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md) | 项目总结（本文件） |

## 🔧 API 参考

### 主要方法

```go
// 创建客户端
client := gptutils.NewDefaultClient()

// 简单对话
response, err := client.SimpleChat(ctx, "你好")

// 流式对话
err := client.SimpleChatStream(ctx, "你好", func(chunk string) error {
    fmt.Print(chunk)
    return nil
})

// 完整对话
req := gptutils.ChatRequest{
    Messages: []gptutils.Message{
        {Role: "user", Content: "你好"},
    },
}
resp, err := client.Chat(ctx, req)

// 流式完整对话
err := client.ChatStream(ctx, req, func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
```

### 数据类型

```go
// 消息
type Message struct {
    Role    string  // "user", "assistant", "system"
    Content string  // 消息内容
}

// 聊天请求
type ChatRequest struct {
    Model       string    // 模型名称
    Messages    []Message // 消息列表
    Stream      bool      // 是否流式输出
    Temperature *float64  // 采样温度
    TopP        *float64  // 核采样概率
    MaxTokens   *int      // 最大输出token数
}

// 聊天响应
type ChatResponse struct {
    ID      string      // 响应ID
    Model   string      // 使用的模型
    Choices []Choice    // 选择列表
    Usage   TokenUsage  // Token使用情况
}
```

## 💡 使用示例

### 示例 1: 基础对话

```bash
go run examples/http_simple_chat.go
```

### 示例 2: 流式输出

```bash
go run examples/http_stream_chat.go
```

### 示例 3: 多轮对话

```bash
go run examples/http_multi_turn.go
```

### 示例 4: SDK 集成

```bash
go run examples/sdk_integration_example.go
```

### 示例 5: 命令行工具

```bash
go run cmd/chat/main.go
```

## 🔐 安全性

- **API Key 管理**: 通过环境变量管理，不硬编码
- **`.env` 文件**: 已添加到 `.gitignore`，不会被提交
- **最佳实践**: 使用环境变量或配置文件管理敏感信息

## 📊 支持的模型

| 模型 | 描述 | 推荐场景 |
|------|------|--------|
| qwen-plus | 通用模型，性能均衡 | 通用对话（推荐） |
| qwen-turbo | 快速响应模型 | 对延迟敏感的应用 |
| qwen-max | 最强性能模型 | 复杂推理、创意写作 |
| qwen-long | 长文本模型 | 处理长文本输入 |

## 🛠️ 开发指南

### 项目依赖

```
github.com/openai/openai-go v0.1.0-alpha.62
```

### 编译

```bash
# 编译客户端库
go build ./client

# 编译命令行工具
go build -o gptutils-chat cmd/chat/main.go

# 编译示例
go build -o example examples/http_simple_chat.go
```

### 测试

```bash
# 运行示例
export API_KEY="your-api-key"
go run examples/sdk_integration_example.go
```

## 📈 性能指标

- **响应时间**: 通常 < 2 秒（取决于网络和模型）
- **吞吐量**: 支持并发请求
- **Token 消耗**: 根据输入输出长度计算

## 🤝 集成方式

### 方式 1: 本地开发

```bash
# 在 go.work 中添加
use (
    ./your-project
    ./GPTUtils
)
```

### 方式 2: 远程模块

```bash
go get github.com/yourusername/gptutils@latest
```

### 方式 3: 本地模块

```bash
# 在 go.mod 中添加
require gptutils v0.1.0
```

## 📝 版本历史

### v0.1.0 (2025-11-25)

- ✨ 初始版本发布
- ✅ 支持基础对话和流式输出
- ✅ 支持多轮对话
- ✅ 提供命令行工具
- ✅ 完整的 SDK 文档和示例
- ✅ 支持多个集成场景

## 🐛 已知问题

1. **OpenAI SDK BaseURL 处理**: OpenAI SDK v0.1.0-alpha.62 在处理 BaseURL 时可能有问题，已通过原生 HTTP 实现解决
2. **工具调用**: 目前规划中，尚未实现
3. **联网搜索**: 目前规划中，尚未实现

## 🔮 未来计划

- [ ] 实现 Function Calling 工具调用
- [ ] 实现联网搜索功能
- [ ] 支持图像输入和分析
- [ ] 支持音频输入和输出
- [ ] 添加缓存机制
- [ ] 添加速率限制
- [ ] 发布到 GitHub
- [ ] 发布到 pkg.go.dev

## 📞 支持

- 📖 [官方 API 文档](https://help.aliyun.com/zh/model-studio/developer-reference/qwen-api)
- 🐛 [提交 Issue](https://github.com/lvdashuaibi/GPTUtils/issues)
- 💬 [讨论区](https://github.com/lvdashuaibi/GPTUtils/discussions)
- 🔗 [GitHub 仓库](https://github.com/lvdashuaibi/GPTUtils)

## 📄 许可证

MIT License

## 🙏 致谢

- 感谢阿里云通义千问团队提供的 API
- 感谢 OpenAI Go SDK 的参考实现

---

**最后更新**: 2025-11-25
**版本**: 0.1.0
**作者**: GPTUtils Team
