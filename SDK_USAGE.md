# GPTUtils SDK 使用指南

GPTUtils 是一个功能完整的通义千问大模型API Go SDK，可以轻松集成到你的项目中。

## 安装

### 方式1: 作为本地模块使用

如果你的项目和 GPTUtils 在同一个工作区，可以直接引用：

```bash
# 在你的项目 go.mod 中添加
require gptutils v0.1.0
```

然后在 `go.work` 中添加：

```
use (
    ./your-project
    ./GPTUtils
)
```

### 方式2: 作为远程模块使用

如果 GPTUtils 发布到 GitHub 或其他 Git 仓库：

```bash
go get github.com/lvdashuaibi/GPTUtils@latest
```

## 快速开始

### 1. 基础对话

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
)

func main() {
    // 创建客户端（自动从 API_KEY 环境变量读取）
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

### 2. 流式输出

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

### 3. 多轮对话

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

### 4. 自定义配置

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
    "gptutils/config"
)

func main() {
    // 创建自定义配置
    cfg := config.DefaultConfig()
    cfg.Model = "qwen-max"  // 使用更强大的模型

    // 创建客户端
    client := gptutils.NewClient(cfg)

    // 使用自定义参数
    ctx := context.Background()
    temperature := 0.9  // 更高的创造性

    req := gptutils.ChatRequest{
        Messages: []gptutils.Message{
            {Role: "user", Content: "写一首关于春天的诗"},
        },
        Temperature: &temperature,
    }

    resp, _ := client.Chat(ctx, req)
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## API 参考

### 客户端方法

#### SimpleChat
```go
func (c *HTTPClient) SimpleChat(ctx context.Context, message string) (string, error)
```
发送单条消息并获取回复。

**参数:**
- `ctx`: 上下文
- `message`: 用户消息

**返回:**
- 回复内容
- 错误信息

#### SimpleChatStream
```go
func (c *HTTPClient) SimpleChatStream(ctx context.Context, message string, handler StreamHandler) error
```
流式发送单条消息。

**参数:**
- `ctx`: 上下文
- `message`: 用户消息
- `handler`: 流式处理函数，每收到一个数据块就调用一次

#### Chat
```go
func (c *HTTPClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
```
发送完整的聊天请求。

**参数:**
- `ctx`: 上下文
- `req`: 聊天请求对象

**返回:**
- 完整的聊天响应
- 错误信息

#### ChatStream
```go
func (c *HTTPClient) ChatStream(ctx context.Context, req ChatRequest, handler StreamHandler) error
```
流式发送完整的聊天请求。

### 数据结构

#### Message
```go
type Message struct {
    Role    string  // "user", "assistant", "system"
    Content string  // 消息内容
}
```

#### ChatRequest
```go
type ChatRequest struct {
    Model       string    // 模型名称，默认 "qwen-plus"
    Messages    []Message // 消息列表
    Stream      bool      // 是否流式输出
    Temperature *float64  // 采样温度 [0, 2)
    TopP        *float64  // 核采样概率 (0, 1.0]
    MaxTokens   *int      // 最大输出token数
}
```

#### ChatResponse
```go
type ChatResponse struct {
    ID      string      // 响应ID
    Model   string      // 使用的模型
    Choices []Choice    // 选择列表
    Usage   TokenUsage  // Token使用情况
}
```

## 环境变量配置

SDK 需要以下环境变量：

```bash
# 必需
export API_KEY="your-api-key-here"

# 可选
export GPTUTILS_MODEL="qwen-plus"  # 默认模型
export GPTUTILS_BASE_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
```

## 支持的模型

- `qwen-plus`: 通用模型，性能均衡（推荐）
- `qwen-turbo`: 快速响应模型
- `qwen-max`: 最强性能模型
- `qwen-long`: 长文本模型

更多模型请参考[官方文档](https://help.aliyun.com/zh/model-studio/getting-started/models)

## 错误处理

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

    response, err := client.SimpleChat(ctx, "你好")
    if err != nil {
        // 处理错误
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Println(response)
}
```

## 高级用法

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

### 自定义HTTP客户端

```go
import "net/http"

cfg := config.DefaultConfig()
client := gptutils.NewClient(cfg)

// 可以通过修改底层HTTP客户端来自定义行为
// client.httpClient 是 *http.Client 类型
```

## 常见问题

### Q: 如何切换模型？

```go
cfg := config.DefaultConfig()
cfg.Model = "qwen-max"
client := gptutils.NewClient(cfg)
```

### Q: 如何处理长文本？

使用 `qwen-long` 模型：

```go
cfg := config.DefaultConfig()
cfg.Model = "qwen-long"
client := gptutils.NewClient(cfg)
```

### Q: 如何降低成本？

1. 使用 `qwen-turbo` 模型
2. 设置 `MaxTokens` 限制输出长度
3. 使用流式输出避免超时重试

```go
maxTokens := 500
req := gptutils.ChatRequest{
    Messages: messages,
    MaxTokens: &maxTokens,
}
```

### Q: 如何提高创意性？

增加温度参数：

```go
temperature := 1.5  // 范围 [0, 2)
req := gptutils.ChatRequest{
    Messages: messages,
    Temperature: &temperature,
}
```

## 完整示例

查看 `examples/` 目录获取更多完整示例：

- `http_simple_chat.go` - 基础对话
- `http_stream_chat.go` - 流式输出
- `http_multi_turn.go` - 多轮对话

## 许可证

MIT License

## 支持

如有问题，请提交 Issue 或 Pull Request。
