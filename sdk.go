// Package gptutils 提供通义千问大模型API的Go SDK
//
// 这是一个功能完整的通义千问API调用工具，支持基础对话、流式输出、多轮对话等功能。
//
// 基本使用示例:
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"gptutils/client"
//		"gptutils/config"
//	)
//
//	func main() {
//		// 创建配置
//		cfg := config.DefaultConfig()
//
//		// 创建客户端
//		c := client.NewHTTPClient(cfg)
//
//		// 简单对话
//		ctx := context.Background()
//		response, _ := c.SimpleChat(ctx, "你好")
//		fmt.Println(response)
//	}
//
// 更多示例请参考 examples/ 目录
package gptutils

import (
	"gptutils/client"
	"gptutils/config"
)

// NewClient 创建新的通义千问客户端
// 这是推荐的创建客户端的方式
func NewClient(cfg *config.Config) *client.HTTPClient {
	return client.NewHTTPClient(cfg)
}

// NewDefaultClient 使用默认配置创建客户端
// 默认配置会从 API_KEY 环境变量读取API密钥
func NewDefaultClient() *client.HTTPClient {
	return client.NewHTTPClient(config.DefaultConfig())
}

// Config 导出配置类型
type Config = config.Config

// Message 导出消息类型
type Message = client.Message

// ChatRequest 导出聊天请求类型
type ChatRequest = client.ChatRequest

// ChatResponse 导出聊天响应类型
type ChatResponse = client.ChatResponse

// StreamHandler 导出流式处理器类型
type StreamHandler = client.StreamHandler
