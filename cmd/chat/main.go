package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"gptutils/client"
	"gptutils/config"
	"log"
	"os"
	"strings"
)

func main() {
	// 命令行参数
	stream := flag.Bool("stream", true, "使用流式输出(默认开启)")
	model := flag.String("model", "qwen-plus", "模型名称")
	temperature := flag.Float64("temperature", 0.7, "采样温度(0-2)")
	flag.Parse()

	// 创建配置
	cfg := config.DefaultConfig()
	cfg.Model = *model

	// 创建客户端
	c := client.NewHTTPClient(cfg)

	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	// 消息历史
	messages := []client.Message{
		{Role: "system", Content: "你是一个友好、专业的AI助手"},
	}

	fmt.Println("=== 通义千问对话工具 ===")
	fmt.Printf("模型: %s\n", *model)
	fmt.Printf("流式输出: %v\n", *stream)
	fmt.Printf("温度: %.1f\n", *temperature)
	fmt.Println("\n命令:")
	fmt.Println("  exit/quit - 退出程序")
	fmt.Println("  clear     - 清空对话历史")
	fmt.Println("  history   - 查看对话历史")
	fmt.Println("========================\n")

	for {
		fmt.Print("你: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// 处理特殊命令
		switch strings.ToLower(input) {
		case "exit", "quit":
			fmt.Println("再见！")
			return
		case "clear":
			messages = []client.Message{
				{Role: "system", Content: "你是一个友好、专业的AI助手"},
			}
			fmt.Println("对话历史已清空")
			continue
		case "history":
			fmt.Println("\n=== 对话历史 ===")
			for i, msg := range messages {
				if msg.Role == "system" {
					continue
				}
				fmt.Printf("%d. [%s]: %s\n", i, msg.Role, msg.Content)
			}
			fmt.Println("================\n")
			continue
		}

		// 添加用户消息
		messages = append(messages, client.Message{Role: "user", Content: input})

		fmt.Print("AI: ")

		if *stream {
			// 流式输出
			req := client.ChatRequest{
				Messages:    messages,
				Temperature: temperature,
			}

			var fullResponse strings.Builder
			err := c.ChatStream(ctx, req, func(chunk string) error {
				fmt.Print(chunk)
				fullResponse.WriteString(chunk)
				return nil
			})

			if err != nil {
				log.Printf("\n错误: %v\n", err)
				// 移除最后添加的用户消息
				messages = messages[:len(messages)-1]
				continue
			}

			fmt.Println()
			// 添加AI回复到历史
			messages = append(messages, client.Message{Role: "assistant", Content: fullResponse.String()})

		} else {
			// 普通对话
			req := client.ChatRequest{
				Messages:    messages,
				Temperature: temperature,
			}

			response, err := c.Chat(ctx, req)
			if err != nil {
				log.Printf("\n错误: %v\n", err)
				messages = messages[:len(messages)-1]
				continue
			}

			if len(response.Choices) > 0 {
				content := response.Choices[0].Message.Content
				fmt.Println(content)
				messages = append(messages, client.Message{Role: "assistant", Content: content})

				// 显示Token使用情况
				fmt.Printf("\n[Token使用: 输入=%d, 输出=%d, 总计=%d]\n",
					response.Usage.PromptTokens,
					response.Usage.CompletionTokens,
					response.Usage.TotalTokens)
			}
		}

		fmt.Println()
	}
}
