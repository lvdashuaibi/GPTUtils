package client

import (
	"context"
	"github.com/lvdashuaibi/GPTUtils/config"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Client 通义千问客户端
type Client struct {
	client *openai.Client
	config *config.Config
}

// NewClient 创建新的客户端
func NewClient(cfg *config.Config) *Client {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
	)

	return &Client{
		client: client,
		config: cfg,
	}
}

// ChatOptions 聊天选项
type ChatOptions struct {
	Model             string                                          // 模型名称
	Messages          []openai.ChatCompletionMessageParamUnion        // 消息列表
	Stream            bool                                            // 是否流式输出
	Temperature       *float64                                        // 采样温度 [0, 2)
	TopP              *float64                                        // 核采样概率 (0, 1.0]
	MaxTokens         *int64                                          // 最大输出token数
	PresencePenalty   *float64                                        // 内容重复度惩罚 [-2.0, 2.0]
	ResponseFormat    *openai.ChatCompletionNewParamsResponseFormat   // 响应格式
	Tools             []openai.ChatCompletionToolParam                // 工具列表
	ToolChoice        openai.ChatCompletionToolChoiceOptionUnionParam // 工具选择策略
	ParallelToolCalls *bool                                           // 是否并行工具调用
	EnableSearch      *bool                                           // 是否启用联网搜索
	Seed              *int64                                          // 随机种子
	Stop              []string                                        // 停止词
	N                 *int64                                          // 生成响应数量
}

// Chat 发送聊天请求
func (c *Client) Chat(ctx context.Context, opts ChatOptions) (*openai.ChatCompletion, error) {
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
	if opts.N != nil {
		params.N = openai.F(*opts.N)
	}

	return c.client.Chat.Completions.New(ctx, params)
}

// SimpleChat 简单聊天接口
func (c *Client) SimpleChat(ctx context.Context, message string) (string, error) {
	opts := ChatOptions{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	}

	completion, err := c.Chat(ctx, opts)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) > 0 {
		return completion.Choices[0].Message.Content, nil
	}

	return "", nil
}

// GetClient 获取底层 OpenAI 客户端
func (c *Client) GetClient() *openai.Client {
	return c.client
}
