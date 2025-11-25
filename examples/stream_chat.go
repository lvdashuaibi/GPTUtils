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

	// 创建客户端
	c := client.NewClient(cfg)

	// 流式对话
	ctx := context.Background()
	fmt.Print("AI回复: ")

	err := c.SimpleChatStream(ctx, "请用100字介绍一下人工智能", func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})

	if err != nil {
		log.Fatalf("\n流式聊天失败: %v", err)
	}

	fmt.Println()
}
