package main

import (
	"context"
	"fmt"
	"gptutils"
	"log"
)

// 这个示例演示如何在其他项目中集成 GPTUtils SDK

func main() {
	// 示例1: 基础对话
	fmt.Println("=== 示例1: 基础对话 ===")
	basicChat()

	fmt.Println("\n=== 示例2: 流式对话 ===")
	streamChat()

	fmt.Println("\n=== 示例3: 多轮对话 ===")
	multiTurnChat()

	fmt.Println("\n=== 示例4: 自定义参数 ===")
	customParamsChat()
}

// 基础对话示例
func basicChat() {
	client := gptutils.NewDefaultClient()
	ctx := context.Background()

	response, err := client.SimpleChat(ctx, "Go语言的主要特点是什么？")
	if err != nil {
		log.Printf("错误: %v", err)
		return
	}

	fmt.Println("回复:", response)
}

// 流式对话示例
func streamChat() {
	client := gptutils.NewDefaultClient()
	ctx := context.Background()

	fmt.Print("AI: ")
	err := client.SimpleChatStream(ctx, "请用50字介绍一下Docker", func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})

	if err != nil {
		log.Printf("错误: %v", err)
		return
	}

	fmt.Println()
}

// 多轮对话示例
func multiTurnChat() {
	client := gptutils.NewDefaultClient()
	ctx := context.Background()

	// 初始化消息
	messages := []gptutils.Message{
		{Role: "system", Content: "你是一个编程助手"},
		{Role: "user", Content: "如何在Go中处理错误？"},
	}

	// 第一轮
	req := gptutils.ChatRequest{
		Messages: messages,
	}

	resp1, err := client.Chat(ctx, req)
	if err != nil {
		log.Printf("错误: %v", err)
		return
	}

	if len(resp1.Choices) > 0 {
		content := resp1.Choices[0].Message.Content
		fmt.Println("第一轮回复:", content[:100]+"...")
		messages = append(messages, gptutils.Message{
			Role:    "assistant",
			Content: content,
		})
	}

	// 第二轮
	messages = append(messages, gptutils.Message{
		Role:    "user",
		Content: "能给我一个具体的代码示例吗？",
	})

	req.Messages = messages
	resp2, err := client.Chat(ctx, req)
	if err != nil {
		log.Printf("错误: %v", err)
		return
	}

	if len(resp2.Choices) > 0 {
		fmt.Println("第二轮回复:", resp2.Choices[0].Message.Content[:100]+"...")
	}
}

// 自定义参数示例
func customParamsChat() {
	client := gptutils.NewDefaultClient()
	ctx := context.Background()

	// 设置参数
	temperature := 0.8
	maxTokens := 500

	req := gptutils.ChatRequest{
		Messages: []gptutils.Message{
			{Role: "user", Content: "写一个关于春天的短诗"},
		},
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	}

	resp, err := client.Chat(ctx, req)
	if err != nil {
		log.Printf("错误: %v", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Println("回复:", resp.Choices[0].Message.Content)
		fmt.Printf("\nToken使用: 输入=%d, 输出=%d, 总计=%d\n",
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens,
			resp.Usage.TotalTokens)
	}
}
