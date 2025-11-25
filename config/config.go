package config

import "os"

// Config 配置结构
type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY environment variable is not set. Please set it before running the application.")
	}

	return &Config{
		APIKey:  apiKey,
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "qwen-plus",
	}
}

// WithAPIKey 设置API Key
func (c *Config) WithAPIKey(apiKey string) *Config {
	c.APIKey = apiKey
	return c
}

// WithBaseURL 设置Base URL
func (c *Config) WithBaseURL(baseURL string) *Config {
	c.BaseURL = baseURL
	return c
}

// WithModel 设置模型名称
func (c *Config) WithModel(model string) *Config {
	c.Model = model
	return c
}
