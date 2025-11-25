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

	// 简单对话
	ctx := context.Background()
	response, err := c.SimpleChat(ctx, "你好，请介绍一下你自己")
	if err != nil {
		log.Fatalf("聊天失败: %v", err)
	}

	fmt.Println("AI回复:", response)
}
