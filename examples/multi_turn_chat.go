package main

import (
	"context"
	"fmt"
	"gptutils/client"
	"gptutils/config"
	"log"

	"github.com/openai/openai-go"
)

func main() {
	// 创建配置
	cfg := config.DefaultConfig()

	// 创建客户端
	c := client.NewClient(cfg)

	ctx := context.Background()

	// 多轮对话示例
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("你是一个友好的AI助手"),
		openai.UserMessage("我想学习Go语言"),
	}

	// 第一轮对话
	opts := client.ChatOptions{
		Messages: messages,
	}

	response1, err := c.Chat(ctx, opts)
	if err != nil {
		log.Fatalf("第一轮对话失败: %v", err)
	}

	if len(response1.Choices) > 0 {
		fmt.Println("第一轮 - AI:", response1.Choices[0].Message.Content)
		// 添加AI的回复到消息历史
		messages = append(messages, openai.AssistantMessage(response1.Choices[0].Message.Content))
	}

	// 第二轮对话
	messages = append(messages, openai.UserMessage("Go语言有哪些优势？"))
	opts.Messages = messages

	response2, err := c.Chat(ctx, opts)
	if err != nil {
		log.Fatalf("第二轮对话失败: %v", err)
	}

	if len(response2.Choices) > 0 {
		fmt.Println("\n第二轮 - AI:", response2.Choices[0].Message.Content)
		messages = append(messages, openai.AssistantMessage(response2.Choices[0].Message.Content))
	}

	// 第三轮对话
	messages = append(messages, openai.UserMessage("能给我一个简单的代码示例吗？"))
	opts.Messages = messages

	response3, err := c.Chat(ctx, opts)
	if err != nil {
		log.Fatalf("第三轮对话失败: %v", err)
	}

	if len(response3.Choices) > 0 {
		fmt.Println("\n第三轮 - AI:", response3.Choices[0].Message.Content)
	}

	// 打印Token使用情况
	fmt.Printf("\n\nToken使用情况:\n")
	fmt.Printf("输入Tokens: %d\n", response3.Usage.PromptTokens)
	fmt.Printf("输出Tokens: %d\n", response3.Usage.CompletionTokens)
	fmt.Printf("总计Tokens: %d\n", response3.Usage.TotalTokens)
}
