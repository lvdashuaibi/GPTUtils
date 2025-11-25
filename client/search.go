package client

import (
	"context"

	"github.com/openai/openai-go"
)

// SearchOptions 联网搜索选项
type SearchOptions struct {
	EnableSearch bool   // 是否启用联网搜索
	ForcedSearch bool   // 是否强制搜索
	SearchQuery  string // 搜索查询（可选）
}

// ChatWithSearch 带联网搜索的聊天
func (c *Client) ChatWithSearch(ctx context.Context, message string, searchOpts SearchOptions) (*openai.ChatCompletion, error) {
	opts := ChatOptions{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
		EnableSearch: &searchOpts.EnableSearch,
	}

	params := openai.ChatCompletionNewParams{
		Messages: openai.F(opts.Messages),
		Model:    openai.F(c.config.Model),
	}

	// 注意: 联网搜索功能需要通过 DashScope 原生 API 实现
	// OpenAI 兼容模式可能不支持此功能

	return c.client.Chat.Completions.New(ctx, params)
}

// ChatWithSearchStream 带联网搜索的流式聊天
func (c *Client) ChatWithSearchStream(ctx context.Context, message string, searchOpts SearchOptions, handler StreamHandler) error {
	opts := ChatOptions{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
		EnableSearch: &searchOpts.EnableSearch,
	}

	params := openai.ChatCompletionNewParams{
		Messages: openai.F(opts.Messages),
		Model:    openai.F(c.config.Model),
	}

	// 注意: 联网搜索功能需要通过 DashScope 原生 API 实现
	// OpenAI 兼容模式可能不支持此功能

	// 设置流式输出选项
	params.StreamOptions = openai.F(openai.ChatCompletionStreamOptionsParam{
		IncludeUsage: openai.F(true),
	})

	stream := c.client.Chat.Completions.NewStreaming(ctx, params)

	// 处理流式响应
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				if err := handler(content); err != nil {
					return err
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return nil
}
