# GPTUtils - 通义千问 API Go SDK

一个功能完整、易于使用的通义千问大模型API Go SDK，支持基础对话、流式输出、多轮对话等功能。

## 🌟 功能特性

- ✅ **基础对话**: 支持单轮和多轮对话
- ✅ **流式输出**: 实时流式返回响应内容，提升用户体验
- ✅ **参数配置**: 支持温度、Top-P、最大Token等参数调整
- ✅ **多模型支持**: 支持 qwen-plus、qwen-turbo、qwen-max 等模型
- ✅ **安全配置**: API Key 通过环境变量管理，不硬编码
- ✅ **原生HTTP实现**: 稳定可靠的HTTP客户端，完全兼容通义千问API
- ✅ **交互式命令行**: 功能完整的命令行对话工具
- ✅ **易于集成**: 可作为SDK集成到其他项目中

## 📦 安装

### 方式1: 本地开发

```bash
# 克隆项目
git clone <repository-url>
cd GPTUtils

# 下载依赖
go mod download
```

### 方式2: 作为SDK使用

在你的项目中引入：

```bash
go get github.com/lvdashuaibi/GPTUtils@latest
```

或在 `go.mod` 中添加：

```
require github.com/lvdashuaibi/GPTUtils v0.1.0
```

## 🚀 快速开始

### 1. 设置环境变量

```bash
export API_KEY="your-api-key-here"
```

或创建 `.env` 文件：

```bash
cp .env.example .env
# 编辑 .env 文件，填入你的 API Key
```

### 2. 基础对话

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
)

func main() {
    // 创建客户端
    client := gptutils.NewDefaultClient()

    // 发送消息
    ctx := context.Background()
    response, err := client.SimpleChat(ctx, "你好，请介绍一下你自己")
    if err != nil {
        panic(err)
    }

    fmt.Println("AI回复:", response)
}
```

### 3. 流式输出

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
)

func main() {
    client := gptutils.NewDefaultClient()
    ctx := context.Background()

    fmt.Print("AI: ")
    err := client.SimpleChatStream(ctx, "请用100字介绍人工智能", func(chunk string) error {
        fmt.Print(chunk)
        return nil
    })

    if err != nil {
        panic(err)
    }
    fmt.Println()
}
```

### 4. 多轮对话

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
)

func main() {
    client := gptutils.NewDefaultClient()
    ctx := context.Background()

    // 初始化消息历史
    messages := []gptutils.Message{
        {Role: "system", Content: "你是一个友好的AI助手"},
        {Role: "user", Content: "我想学习Go语言"},
    }

    // 第一轮对话
    req := gptutils.ChatRequest{
        Messages: messages,
    }

    resp1, _ := client.Chat(ctx, req)
    if len(resp1.Choices) > 0 {
        content := resp1.Choices[0].Message.Content
        fmt.Println("AI:", content)
        messages = append(messages, gptutils.Message{
            Role: "assistant",
            Content: content,
        })
    }

    // 第二轮对话
    messages = append(messages, gptutils.Message{
        Role: "user",
        Content: "Go语言有哪些优势？",
    })
    req.Messages = messages

    resp2, _ := client.Chat(ctx, req)
    if len(resp2.Choices) > 0 {
        fmt.Println("AI:", resp2.Choices[0].Message.Content)
    }
}
```

## 📚 详细文档

### SDK 使用指南

详见 [SDK_USAGE.md](./SDK_USAGE.md)，包含：

- 完整的API参考
- 高级用法示例
- 常见问题解答
- 错误处理方法

### 命令行工具

```bash
# 设置环境变量
export API_KEY="your-api-key-here"

# 运行命令行工具
go run cmd/chat/main.go

# 使用流式输出（默认开启）
go run cmd/chat/main.go -stream

# 指定模型
go run cmd/chat/main.go -model qwen-max

# 调整温度参数
go run cmd/chat/main.go -temperature 0.9
```

命令行工具支持的命令：
- `exit` 或 `quit`: 退出程序
- `clear`: 清空对话历史
- `history`: 查看对话历史

## 📖 示例程序

项目包含多个示例程序，位于 `examples/` 目录：

```bash
# 设置环境变量
export API_KEY="your-api-key-here"

# 基础对话
go run examples/http_simple_chat.go

# 流式输出
go run examples/http_stream_chat.go

# 多轮对话
go run examples/http_multi_turn.go
```

## 🔧 API 参考

### 创建客户端

```go
// 使用默认配置
client := gptutils.NewDefaultClient()

// 使用自定义配置
cfg := &gptutils.Config{
    APIKey:  "your-api-key",
    BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
    Model:   "qwen-plus",
}
client := gptutils.NewClient(cfg)
```

### 主要方法

#### SimpleChat - 简单对话

```go
response, err := client.SimpleChat(ctx, "你好")
```

#### SimpleChatStream - 流式对话

```go
err := client.SimpleChatStream(ctx, "你好", func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
```

#### Chat - 完整对话

```go
req := gptutils.ChatRequest{
    Messages: []gptutils.Message{
        {Role: "user", Content: "你好"},
    },
    Temperature: &temperature,
    MaxTokens: &maxTokens,
}
resp, err := client.Chat(ctx, req)
```

#### ChatStream - 流式完整对话

```go
req := gptutils.ChatRequest{
    Messages: messages,
}
err := client.ChatStream(ctx, req, func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
```

## 📋 支持的模型

| 模型名称 | 描述 | 推荐场景 |
|---------|------|--------|
| `qwen-plus` | 通用模型，性能均衡 | 通用对话、推荐使用 |
| `qwen-turbo` | 快速响应模型 | 对延迟敏感的应用 |
| `qwen-max` | 最强性能模型 | 复杂推理、创意写作 |
| `qwen-long` | 长文本模型 | 处理长文本输入 |

更多模型请参考[官方文档](https://help.aliyun.com/zh/model-studio/getting-started/models)

## 🔐 安全性

- **API Key 管理**:
  - 永远不要将 API Key 硬编码在代码中
  - 使用环境变量或配置文件管理
  - `.env` 文件已添加到 `.gitignore`

- **最佳实践**:
  ```bash
  # 设置环境变量
  export API_KEY="your-api-key-here"

  # 或使用 .env 文件
  cp .env.example .env
  # 编辑 .env 文件
  ```

## 💡 高级用法

### 自定义参数

```go
temperature := 0.9
topP := 0.8
maxTokens := 2000

req := gptutils.ChatRequest{
    Messages:    messages,
    Temperature: &temperature,
    TopP:        &topP,
    MaxTokens:   &maxTokens,
}
```

### 设置超时

```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := client.SimpleChat(ctx, "你好")
```

### 获取Token使用情况

```go
resp, _ := client.Chat(ctx, req)
fmt.Printf("输入Tokens: %d\n", resp.Usage.PromptTokens)
fmt.Printf("输出Tokens: %d\n", resp.Usage.CompletionTokens)
fmt.Printf("总计Tokens: %d\n", resp.Usage.TotalTokens)
```

## 📁 项目结构

```
GPTUtils/
├── client/              # 客户端实现
│   ├── http_client.go  # HTTP客户端
│   ├── types.go        # 数据类型定义
│   ├── client.go       # 基础客户端（OpenAI SDK）
│   ├── stream.go       # 流式输出
│   ├── tools.go        # 工具调用
│   └── search.go       # 联网搜索
├── config/             # 配置管理
│   └── config.go
├── cmd/                # 命令行工具
│   └── chat/
│       └── main.go
├── examples/           # 示例程序
│   ├── http_simple_chat.go
│   ├── http_stream_chat.go
│   └── http_multi_turn.go
├── sdk.go              # SDK 导出接口
├── SDK_USAGE.md        # SDK 使用指南
├── .env.example        # 环境变量示例
├── .gitignore          # Git 忽略文件
├── go.mod              # Go 模块文件
├── VERSION             # 版本号
└── README.md           # 项目文档
```

## 🤝 集成到其他项目

### 步骤1: 添加依赖

在你的项目 `go.mod` 中：

```
require gptutils v0.1.0
```

### 步骤2: 导入并使用

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

### 步骤3: 设置环境变量

```bash
export API_KEY="your-api-key-here"
go run main.go
```

## ❓ 常见问题

### Q: 如何切换模型？

```go
cfg := &gptutils.Config{
    APIKey: os.Getenv("API_KEY"),
    Model:  "qwen-max",
}
client := gptutils.NewClient(cfg)
```

### Q: 如何处理错误？

```go
response, err := client.SimpleChat(ctx, "你好")
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}
```

### Q: 如何提高创意性？

增加温度参数：

```go
temperature := 1.5  // 范围 [0, 2)
req := gptutils.ChatRequest{
    Messages:    messages,
    Temperature: &temperature,
}
```

### Q: 如何降低成本？

1. 使用 `qwen-turbo` 模型
2. 设置 `MaxTokens` 限制输出长度
3. 使用流式输出避免超时重试

## 📝 许可证

MIT License

## 🙏 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

- 📖 [官方文档](https://help.aliyun.com/zh/model-studio/developer-reference/qwen-api)
- 🐛 [提交Issue](https://github.com/yourusername/gptutils/issues)
- 💬 [讨论区](https://github.com/yourusername/gptutils/discussions)

## 更新日志

### v0.1.0 (2025-11-24)

- ✨ 初始版本发布
- ✅ 支持基础对话和流式输出
- ✅ 支持多轮对话
- ✅ 提供命令行工具
- ✅ 完整的SDK文档和示例
