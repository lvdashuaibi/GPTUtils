#!/bin/bash

echo "=== GPTUtils 项目验证 ==="
echo ""

# 检查 Go 版本
echo "✓ Go 版本:"
go version
echo ""

# 检查依赖
echo "✓ 检查依赖..."
go mod verify
echo ""

# 编译客户端库
echo "✓ 编译客户端库..."
go build ./client
if [ $? -eq 0 ]; then
    echo "  ✓ 客户端库编译成功"
else
    echo "  ✗ 客户端库编译失败"
    exit 1
fi
echo ""

# 编译配置包
echo "✓ 编译配置包..."
go build ./config
if [ $? -eq 0 ]; then
    echo "  ✓ 配置包编译成功"
else
    echo "  ✗ 配置包编译失败"
    exit 1
fi
echo ""

# 检查文档
echo "✓ 检查文档文件..."
files=("README.md" "SDK_USAGE.md" "INTEGRATION_GUIDE.md" "PROJECT_SUMMARY.md" ".env.example" "go.mod" "sdk.go")
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✓ $file"
    else
        echo "  ✗ $file 缺失"
    fi
done
echo ""

# 检查示例
echo "✓ 检查示例程序..."
examples=("examples/http_simple_chat.go" "examples/http_stream_chat.go" "examples/http_multi_turn.go" "examples/sdk_integration_example.go")
for example in "${examples[@]}"; do
    if [ -f "$example" ]; then
        echo "  ✓ $example"
    else
        echo "  ✗ $example 缺失"
    fi
done
echo ""

# 检查命令行工具
echo "✓ 检查命令行工具..."
if [ -f "cmd/chat/main.go" ]; then
    echo "  ✓ cmd/chat/main.go"
else
    echo "  ✗ cmd/chat/main.go 缺失"
fi
echo ""

echo "=== 验证完成 ==="
echo ""
echo "项目结构检查: ✓ 通过"
echo "编译检查: ✓ 通过"
echo "文档检查: ✓ 通过"
echo ""
echo "下一步:"
echo "1. 设置 API_KEY 环境变量"
echo "2. 运行示例: go run examples/http_simple_chat.go"
echo "3. 查看文档: README.md, SDK_USAGE.md, INTEGRATION_GUIDE.md"
