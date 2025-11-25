package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gptutils/config"
	"io"
	"net/http"
	"strings"
)

// HTTPClient 基于原生HTTP的客户端实现
// 用于调用通义千问API的HTTP客户端
type HTTPClient struct {
	config     *config.Config
	httpClient *http.Client
}

// NewHTTPClient 创建HTTP客户端
// cfg: 配置对象，如果为nil则使用默认配置
func NewHTTPClient(cfg *config.Config) *HTTPClient {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	return &HTTPClient{
		config:     cfg,
		httpClient: &http.Client{},
	}
}

// Chat 发送聊天请求
func (c *HTTPClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = c.config.Model
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		c.config.BaseURL+"/chat/completions",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, err
	}

	return &chatResp, nil
}

// ChatStream 流式聊天
func (c *HTTPClient) ChatStream(ctx context.Context, req ChatRequest, handler func(content string) error) error {
	if req.Model == "" {
		req.Model = c.config.Model
	}
	req.Stream = true

	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		c.config.BaseURL+"/chat/completions",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				if err := handler(chunk.Choices[0].Delta.Content); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SimpleChat 简单对话
func (c *HTTPClient) SimpleChat(ctx context.Context, message string) (string, error) {
	req := ChatRequest{
		Messages: []Message{
			{Role: "user", Content: message},
		},
	}

	resp, err := c.Chat(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", nil
}

// SimpleChatStream 简单流式对话
func (c *HTTPClient) SimpleChatStream(ctx context.Context, message string, handler func(content string) error) error {
	req := ChatRequest{
		Messages: []Message{
			{Role: "user", Content: message},
		},
	}

	return c.ChatStream(ctx, req, handler)
}
