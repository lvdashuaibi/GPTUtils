package main

import (
	"context"
	"fmt"
	"github.com/lvdashuaibi/GPTUtils/client"
	"github.com/lvdashuaibi/GPTUtils/config"
	"log"

	"github.com/openai/openai-go"
)

func main() {
	// 创建配置
	cfg := config.DefaultConfig()

	// 创建客户端
	c := client.NewClient(cfg)

	// 创建工具管理器
	toolManager := client.NewToolManager()

	// 注册天气查询工具
	toolManager.RegisterTool(client.CreateWeatherTool())

	// 注册计算器工具
	toolManager.RegisterTool(client.CreateCalculatorTool())

	// 带工具调用的对话
	ctx := context.Background()
	opts := client.ChatOptions{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("北京今天天气怎么样？"),
		},
	}

	completion, err := c.ChatWithTools(ctx, opts, toolManager, 5)
	if err != nil {
		log.Fatalf("工具调用失败: %v", err)
	}

	if len(completion.Choices) > 0 {
		fmt.Println("AI回复:", completion.Choices[0].Message.Content)
	}
}
