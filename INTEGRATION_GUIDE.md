# GPTUtils SDK 集成指南

本指南说明如何将 GPTUtils SDK 集成到你的 Go 项目中。

## 📋 目录

1. [快速集成](#快速集成)
2. [项目结构](#项目结构)
3. [使用示例](#使用示例)
4. [常见集成场景](#常见集成场景)
5. [故障排除](#故障排除)

## 快速集成

### 步骤 1: 添加依赖

在你的项目 `go.mod` 中添加：

```
require github.com/lvdashuaibi/GPTUtils v0.1.0
```

或使用 go get：

```bash
go get github.com/lvdashuaibi/GPTUtils@latest
```

### 步骤 2: 设置环境变量

```bash
export API_KEY="your-api-key-here"
```

### 步骤 3: 在代码中使用

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

## 项目结构

```
GPTUtils/
├── client/                 # 核心客户端实现
│   ├── http_client.go     # HTTP 客户端（推荐使用）
│   ├── types.go           # 数据类型定义
│   ├── client.go          # OpenAI SDK 客户端（备选）
│   ├── stream.go          # 流式输出支持
│   ├── tools.go           # 工具调用（规划中）
│   └── search.go          # 联网搜索（规划中）
├── config/                # 配置管理
│   └── config.go
├── cmd/                   # 命令行工具
│   └── chat/main.go
├── examples/              # 示例程序
│   ├── http_simple_chat.go
│   ├── http_stream_chat.go
│   ├── http_multi_turn.go
│   └── sdk_integration_example.go
├── sdk.go                 # SDK 公共接口
├── SDK_USAGE.md           # SDK 详细文档
├── README.md              # 项目文档
└── INTEGRATION_GUIDE.md   # 本文件
```

## 使用示例

### 示例 1: 基础对话

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

    response, err := client.SimpleChat(ctx, "你好，请介绍一下你自己")
    if err != nil {
        panic(err)
    }

    fmt.Println("AI:", response)
}
```

### 示例 2: 流式输出

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
    err := client.SimpleChatStream(ctx, "请介绍一下人工智能", func(chunk string) error {
        fmt.Print(chunk)
        return nil
    })

    if err != nil {
        panic(err)
    }
    fmt.Println()
}
```

### 示例 3: 多轮对话

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

    messages := []gptutils.Message{
        {Role: "system", Content: "你是一个编程助手"},
        {Role: "user", Content: "如何学习Go语言？"},
    }

    // 第一轮
    req := gptutils.ChatRequest{Messages: messages}
    resp1, _ := client.Chat(ctx, req)

    if len(resp1.Choices) > 0 {
        content := resp1.Choices[0].Message.Content
        fmt.Println("AI:", content)
        messages = append(messages, gptutils.Message{
            Role: "assistant",
            Content: content,
        })
    }

    // 第二轮
    messages = append(messages, gptutils.Message{
        Role: "user",
        Content: "推荐一些学习资源",
    })

    req.Messages = messages
    resp2, _ := client.Chat(ctx, req)
    if len(resp2.Choices) > 0 {
        fmt.Println("AI:", resp2.Choices[0].Message.Content)
    }
}
```

### 示例 4: 自定义配置

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

    client := gptutils.NewClient(cfg)
    ctx := context.Background()

    // 设置参数
    temperature := 0.9
    maxTokens := 1000

    req := gptutils.ChatRequest{
        Messages: []gptutils.Message{
            {Role: "user", Content: "写一首诗"},
        },
        Temperature: &temperature,
        MaxTokens: &maxTokens,
    }

    resp, _ := client.Chat(ctx, req)
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## 常见集成场景

### 场景 1: 集成到 Web 服务

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "gptutils"
    "net/http"
)

var client *gptutils.HTTPClient

func init() {
    client = gptutils.NewDefaultClient()
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Message string `json:"message"`
    }

    json.NewDecoder(r.Body).Decode(&req)

    response, err := client.SimpleChat(context.Background(), req.Message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "response": response,
    })
}

func main() {
    http.HandleFunc("/chat", chatHandler)
    http.ListenAndServe(":8080", nil)
}
```

### 场景 2: 集成到 CLI 工具

```go
package main

import (
    "bufio"
    "context"
    "fmt"
    "gptutils"
    "os"
    "strings"
)

func main() {
    client := gptutils.NewDefaultClient()
    ctx := context.Background()

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("你: ")
        if !scanner.Scan() {
            break
        }

        input := strings.TrimSpace(scanner.Text())
        if input == "exit" {
            break
        }

        response, _ := client.SimpleChat(ctx, input)
        fmt.Println("AI:", response)
    }
}
```

### 场景 3: 集成到数据处理管道

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
)

func processData(data []string) []string {
    client := gptutils.NewDefaultClient()
    ctx := context.Background()

    results := make([]string, len(data))

    for i, item := range data {
        prompt := fmt.Sprintf("分析以下内容：%s", item)
        response, _ := client.SimpleChat(ctx, prompt)
        results[i] = response
    }

    return results
}

func main() {
    data := []string{"数据1", "数据2", "数据3"}
    results := processData(data)

    for _, result := range results {
        fmt.Println(result)
    }
}
```

### 场景 4: 集成到日志分析系统

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
)

func analyzeLog(logContent string) string {
    client := gptutils.NewDefaultClient()
    ctx := context.Background()

    prompt := fmt.Sprintf(`
分析以下日志内容，识别问题和建议：

%s

请提供：
1. 问题诊断
2. 根本原因
3. 解决方案
`, logContent)

    response, _ := client.SimpleChat(ctx, prompt)
    return response
}

func main() {
    logContent := `
2025-11-25 10:30:45 ERROR [main] Connection timeout
2025-11-25 10:30:46 ERROR [main] Retry failed
2025-11-25 10:30:47 ERROR [main] Service unavailable
`

    analysis := analyzeLog(logContent)
    fmt.Println("日志分析结果：")
    fmt.Println(analysis)
}
```

## 常见集成场景

### 场景 5: 集成到测试框架

```go
package main

import (
    "context"
    "fmt"
    "gptutils"
    "testing"
)

func TestAIResponse(t *testing.T) {
    client := gptutils.NewDefaultClient()
    ctx := context.Background()

    response, err := client.SimpleChat(ctx, "1+1等于多少？")
    if err != nil {
        t.Fatalf("API 调用失败: %v", err)
    }

    if response == "" {
        t.Error("收到空响应")
    }

    fmt.Println("测试通过，AI 回复:", response)
}
```

## 故障排除

### 问题 1: API_KEY 未设置

**错误信息:**
```
panic: API_KEY environment variable is not set
```

**解决方案:**
```bash
export API_KEY="your-api-key-here"
```

### 问题 2: 网络连接超时

**错误信息:**
```
context deadline exceeded
```

**解决方案:**
```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

response, err := client.SimpleChat(ctx, "你好")
```

### 问题 3: 模型不存在

**错误信息:**
```
model not found
```

**解决方案:**
检查模型名称是否正确，支持的模型有：
- qwen-plus（推荐）
- qwen-turbo
- qwen-max
- qwen-long

### 问题 4: Token 超限

**错误信息:**
```
max_tokens exceeded
```

**解决方案:**
```go
maxTokens := 500
req := gptutils.ChatRequest{
    Messages: messages,
    MaxTokens: &maxTokens,
}
```

## 最佳实践

1. **错误处理**: 始终检查返回的错误
   ```go
   response, err := client.SimpleChat(ctx, "你好")
   if err != nil {
       log.Printf("错误: %v", err)
       return
   }
   ```

2. **超时设置**: 为长时间操作设置超时
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

3. **资源管理**: 复用客户端实例
   ```go
   // 好的做法
   client := gptutils.NewDefaultClient()
   // 在整个应用生命周期中复用 client

   // 不好的做法
   for i := 0; i < 1000; i++ {
       client := gptutils.NewDefaultClient()  // 重复创建
   }
   ```

4. **流式输出**: 对长响应使用流式输出
   ```go
   err := client.SimpleChatStream(ctx, "长问题", func(chunk string) error {
       fmt.Print(chunk)
       return nil
   })
   ```

5. **参数调优**: 根据场景调整参数
   ```go
   // 创意写作：高温度
   temperature := 1.5

   // 精确回答：低温度
   temperature := 0.3
   ```

## 获取帮助

- 📖 [SDK 详细文档](./SDK_USAGE.md)
- 📚 [项目 README](./README.md)
- 🔗 [官方 API 文档](https://help.aliyun.com/zh/model-studio/developer-reference/qwen-api)
- 🐛 [GitHub Issues](https://github.com/lvdashuaibi/GPTUtils/issues)
- 💬 [GitHub Discussions](https://github.com/lvdashuaibi/GPTUtils/discussions)

## 许可证

MIT License
