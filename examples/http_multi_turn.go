package main

import (
	"context"
	"fmt"
	"github.com/lvdashuaibi/GPTUtils/client"
	"github.com/lvdashuaibi/GPTUtils/config"
	"log"
)

func main() {
	// 创建配置
	cfg := config.DefaultConfig()

	// 创建HTTP客户端
	c := client.NewHTTPClient(cfg)

	ctx := context.Background()

	// 多轮对话示例
	messages := []client.Message{
		{Role: "system", Content: "你是一个友好的AI助手"},
		{Role: "user", Content: "我想学习Go语言"},
	}

	// 第一轮对话
	req := client.ChatRequest{
		Messages: messages,
	}

	response1, err := c.Chat(ctx, req)
	if err != nil {
		log.Fatalf("第一轮对话失败: %v", err)
	}

	if len(response1.Choices) > 0 {
		content := response1.Choices[0].Message.Content
		fmt.Println("第一轮 - AI:", content)
		messages = append(messages, client.Message{Role: "assistant", Content: content})
	}

	// 第二轮对话
	messages = append(messages, client.Message{Role: "user", Content: "Go语言有哪些优势？"})
	req.Messages = messages

	response2, err := c.Chat(ctx, req)
	if err != nil {
		log.Fatalf("第二轮对话失败: %v", err)
	}

	if len(response2.Choices) > 0 {
		content := response2.Choices[0].Message.Content
		fmt.Println("\n第二轮 - AI:", content)
		messages = append(messages, client.Message{Role: "assistant", Content: content})
	}

	// 第三轮对话
	messages = append(messages, client.Message{Role: "user", Content: "能给我一个简单的代码示例吗？"})
	req.Messages = messages

	response3, err := c.Chat(ctx, req)
	if err != nil {
		log.Fatalf("第三轮对话失败: %v", err)
	}

	if len(response3.Choices) > 0 {
		content := response3.Choices[0].Message.Content
		fmt.Println("\n第三轮 - AI:", content)
	}

	// 打印Token使用情况
	fmt.Printf("\n\nToken使用情况:\n")
	fmt.Printf("输入Tokens: %d\n", response3.Usage.PromptTokens)
	fmt.Printf("输出Tokens: %d\n", response3.Usage.CompletionTokens)
	fmt.Printf("总计Tokens: %d\n", response3.Usage.TotalTokens)
}
