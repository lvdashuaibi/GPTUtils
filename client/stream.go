package client

import (
	"context"
	"io"

	"github.com/openai/openai-go"
)

// ChatStream 流式聊天
func (c *Client) ChatStream(ctx context.Context, opts ChatOptions, handler StreamHandler) error {
	if opts.Model == "" {
		opts.Model = c.config.Model
	}

	params := openai.ChatCompletionNewParams{
		Messages: openai.F(opts.Messages),
		Model:    openai.F(opts.Model),
	}

	// 设置可选参数
	if opts.Temperature != nil {
		params.Temperature = openai.F(*opts.Temperature)
	}
	if opts.TopP != nil {
		params.TopP = openai.F(*opts.TopP)
	}
	if opts.MaxTokens != nil {
		params.MaxTokens = openai.F(*opts.MaxTokens)
	}
	if opts.PresencePenalty != nil {
		params.PresencePenalty = openai.F(*opts.PresencePenalty)
	}
	if len(opts.Tools) > 0 {
		params.Tools = openai.F(opts.Tools)
	}
	if opts.ToolChoice != nil {
		params.ToolChoice = openai.F(opts.ToolChoice)
	}
	if opts.ParallelToolCalls != nil {
		params.ParallelToolCalls = openai.F(*opts.ParallelToolCalls)
	}
	if opts.Seed != nil {
		params.Seed = openai.F(*opts.Seed)
	}
	if len(opts.Stop) > 0 {
		params.Stop = openai.F[openai.ChatCompletionNewParamsStopUnion](
			openai.ChatCompletionNewParamsStopArray(opts.Stop),
		)
	}

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

// SimpleChatStream 简单流式聊天
func (c *Client) SimpleChatStream(ctx context.Context, message string, handler StreamHandler) error {
	opts := ChatOptions{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	}

	return c.ChatStream(ctx, opts, handler)
}

// ChatStreamToWriter 流式输出到 Writer
func (c *Client) ChatStreamToWriter(ctx context.Context, opts ChatOptions, writer io.Writer) error {
	return c.ChatStream(ctx, opts, func(chunk string) error {
		_, err := writer.Write([]byte(chunk))
		return err
	})
}
