package main

import (
	"context"
	"fmt"
	"gptutils/client"
	"gptutils/config"
	"log"
)

func main() {
	// 创建配置
	cfg := config.DefaultConfig()

	// 创建客户端
	c := client.NewClient(cfg)

	// 带联网搜索的对话
	ctx := context.Background()
	searchOpts := client.SearchOptions{
		EnableSearch: true,
		ForcedSearch: false, // 让模型自主决定是否搜索
	}

	response, err := c.ChatWithSearch(ctx, "2024年最新的AI技术发展趋势是什么？", searchOpts)
	if err != nil {
		log.Fatalf("联网搜索失败: %v", err)
	}

	if len(response.Choices) > 0 {
		fmt.Println("AI回复:", response.Choices[0].Message.Content)
	}

	// 流式联网搜索
	fmt.Println("\n\n=== 流式联网搜索 ===")
	fmt.Print("AI回复: ")

	err = c.ChatWithSearchStream(ctx, "今天的新闻头条有哪些？", searchOpts, func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})

	if err != nil {
		log.Fatalf("\n流式联网搜索失败: %v", err)
	}

	fmt.Println()
}
