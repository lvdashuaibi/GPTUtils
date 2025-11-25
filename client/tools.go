package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
)

// ToolFunction 工具函数定义
type ToolFunction func(args string) (string, error)

// Tool 工具定义
type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Function    ToolFunction
}

// ToolManager 工具管理器
type ToolManager struct {
	tools map[string]*Tool
}

// NewToolManager 创建工具管理器
func NewToolManager() *ToolManager {
	return &ToolManager{
		tools: make(map[string]*Tool),
	}
}

// RegisterTool 注册工具
func (tm *ToolManager) RegisterTool(tool *Tool) {
	tm.tools[tool.Name] = tool
}

// GetTool 获取工具
func (tm *ToolManager) GetTool(name string) (*Tool, bool) {
	tool, ok := tm.tools[name]
	return tool, ok
}

// GetToolParams 获取工具参数列表
func (tm *ToolManager) GetToolParams() []openai.ChatCompletionToolParam {
	params := make([]openai.ChatCompletionToolParam, 0, len(tm.tools))
	for _, tool := range tm.tools {
		params = append(params, openai.ChatCompletionToolParam{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String(tool.Name),
				Description: openai.String(tool.Description),
				Parameters:  openai.F(openai.FunctionParameters(tool.Parameters)),
			}),
		})
	}
	return params
}

// ExecuteTool 执行工具
func (tm *ToolManager) ExecuteTool(name string, args string) (string, error) {
	tool, ok := tm.GetTool(name)
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	return tool.Function(args)
}

// ChatWithTools 带工具调用的聊天
func (c *Client) ChatWithTools(ctx context.Context, opts ChatOptions, toolManager *ToolManager, maxIterations int) (*openai.ChatCompletion, error) {
	if maxIterations <= 0 {
		maxIterations = 5 // 默认最多5轮工具调用
	}

	// 设置工具参数
	if toolManager != nil && len(toolManager.tools) > 0 {
		opts.Tools = toolManager.GetToolParams()
	}

	messages := opts.Messages
	var lastCompletion *openai.ChatCompletion

	for i := 0; i < maxIterations; i++ {
		opts.Messages = messages
		completion, err := c.Chat(ctx, opts)
		if err != nil {
			return nil, err
		}

		lastCompletion = completion

		// 检查是否需要调用工具
		if len(completion.Choices) == 0 {
			break
		}

		choice := completion.Choices[0]
		if choice.FinishReason != openai.ChatCompletionChoicesFinishReasonToolCalls {
			// 没有工具调用，返回结果
			break
		}

		// 添加助手消息
		messages = append(messages, openai.AssistantMessage(choice.Message.Content))

		// 处理工具调用
		if len(choice.Message.ToolCalls) == 0 {
			break
		}

		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Type != openai.ChatCompletionMessageToolCallTypeFunction {
				continue
			}

			functionName := toolCall.Function.Name
			functionArgs := toolCall.Function.Arguments

			// 执行工具
			result, err := toolManager.ExecuteTool(functionName, functionArgs)
			if err != nil {
				result = fmt.Sprintf("Error executing tool: %v", err)
			}

			// 添加工具响应消息
			messages = append(messages, openai.ToolMessage(toolCall.ID, result))
		}
	}

	return lastCompletion, nil
}

// CreateWeatherTool 创建天气查询工具示例
func CreateWeatherTool() *Tool {
	return &Tool{
		Name:        "get_weather",
		Description: "获取指定城市的天气信息",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "城市名称，例如：北京、上海",
				},
				"unit": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"celsius", "fahrenheit"},
					"description": "温度单位",
				},
			},
			"required": []string{"location"},
		},
		Function: func(args string) (string, error) {
			var params struct {
				Location string `json:"location"`
				Unit     string `json:"unit"`
			}
			if err := json.Unmarshal([]byte(args), &params); err != nil {
				return "", err
			}

			// 模拟天气数据
			weather := map[string]interface{}{
				"location":    params.Location,
				"temperature": 22,
				"unit":        params.Unit,
				"condition":   "晴天",
				"humidity":    65,
			}

			result, _ := json.Marshal(weather)
			return string(result), nil
		},
	}
}

// CreateCalculatorTool 创建计算器工具示例
func CreateCalculatorTool() *Tool {
	return &Tool{
		Name:        "calculator",
		Description: "执行基本的数学计算",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"expression": map[string]interface{}{
					"type":        "string",
					"description": "数学表达式，例如：2+2, 10*5",
				},
			},
			"required": []string{"expression"},
		},
		Function: func(args string) (string, error) {
			var params struct {
				Expression string `json:"expression"`
			}
			if err := json.Unmarshal([]byte(args), &params); err != nil {
				return "", err
			}

			// 这里简化处理，实际应该使用表达式解析器
			return fmt.Sprintf("计算结果: %s = [需要实现表达式计算]", params.Expression), nil
		},
	}
}
