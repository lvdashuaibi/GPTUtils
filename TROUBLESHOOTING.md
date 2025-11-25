# GPTUtils 故障排除指南

## 包无法被 pkg.go.dev 索引

### 问题描述

在尝试访问 `github.com/lvdashuaibi/GPTUtils` 时收到错误：
```
Oops! We couldn't find "github.com/lvdashuaibi/GPTUtils".
```

### 原因分析

1. **Module 路径不正确**: `go.mod` 中的 module 名称必须与 GitHub 仓库路径一致
2. **缺少版本标签**: 需要创建 Git tag 来标记版本
3. **索引延迟**: pkg.go.dev 需要时间来索引新包

### 解决方案

#### 步骤 1: 验证 go.mod 文件

确保 `go.mod` 的第一行是：
```
module github.com/lvdashuaibi/GPTUtils
```

❌ 错误的方式：
```
module gptutils
```

✅ 正确的方式：
```
module github.com/lvdashuaibi/GPTUtils
```

#### 步骤 2: 更新所有导入语句

所有 Go 文件中的导入语句应该使用完整路径：

❌ 错误的方式：
```go
import "gptutils/client"
import "gptutils/config"
```

✅ 正确的方式：
```go
import "github.com/lvdashuaibi/GPTUtils/client"
import "github.com/lvdashuaibi/GPTUtils/config"
```

#### 步骤 3: 创建版本标签

```bash
# 创建标签
git tag v0.1.0

# 推送标签到 GitHub
git push origin v0.1.0

# 推送所有标签
git push origin --tags
```

#### 步骤 4: 等待索引

pkg.go.dev 通常在 24 小时内索引新包。你可以：

1. 访问 https://pkg.go.dev/github.com/lvdashuaibi/GPTUtils 检查状态
2. 或者手动触发索引：访问 https://pkg.go.dev/github.com/lvdashuaibi/GPTUtils?tab=doc

### 验证修复

#### 方式 1: 在本地测试

```bash
# 创建测试项目
mkdir test-gptutils
cd test-gptutils
go mod init test-gptutils

# 添加依赖
go get github.com/lvdashuaibi/GPTUtils@v0.1.0

# 创建测试文件
cat > main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "github.com/lvdashuaibi/GPTUtils"
)

func main() {
    client := gptutils.NewDefaultClient()
    response, _ := client.SimpleChat(context.Background(), "你好")
    fmt.Println(response)
}
EOF

# 运行测试
go run main.go
```

#### 方式 2: 检查 pkg.go.dev

访问以下 URL 检查包是否已被索引：
- https://pkg.go.dev/github.com/lvdashuaibi/GPTUtils
- https://pkg.go.dev/github.com/lvdashuaibi/GPTUtils/client
- https://pkg.go.dev/github.com/lvdashuaibi/GPTUtils/config

### 常见错误

#### 错误 1: "module declares its identity as X but was required as Y"

**原因**: go.mod 中的 module 名称与导入路径不匹配

**解决方案**:
```bash
# 检查 go.mod
cat go.mod | head -1

# 应该是
module github.com/lvdashuaibi/GPTUtils
```

#### 错误 2: "cannot find module providing package"

**原因**: 导入路径不正确

**解决方案**:
```bash
# 错误的导入
import "gptutils/client"

# 正确的导入
import "github.com/lvdashuaibi/GPTUtils/client"
```

#### 错误 3: "go get: no matching versions for query"

**原因**: 版本标签不存在或格式不正确

**解决方案**:
```bash
# 查看已有的标签
git tag -l

# 创建正确格式的标签
git tag v0.1.0
git push origin v0.1.0
```

### 完整修复步骤

如果上述方法都不起作用，按照以下步骤完整修复：

```bash
# 1. 确保 go.mod 正确
cat go.mod | head -1
# 应该输出: module github.com/lvdashuaibi/GPTUtils

# 2. 更新所有导入语句
find . -name "*.go" -type f -exec sed -i '' 's/"gptutils\//"github.com\/lvdashuaibi\/GPTUtils\//g' {} \;

# 3. 验证编译
go build ./client
go build ./config

# 4. 提交更改
git add -A
git commit -m "Fix: Update module path for pkg.go.dev compatibility"

# 5. 创建版本标签
git tag v0.1.0

# 6. 推送到 GitHub
git push origin main
git push origin v0.1.0

# 7. 验证
go get github.com/lvdashuaibi/GPTUtils@v0.1.0
```

### 后续版本发布

对于后续版本，遵循以下流程：

```bash
# 1. 修改 VERSION 文件
echo "0.2.0" > VERSION

# 2. 更新 go.mod 中的版本注释（可选）
# 3. 提交更改
git add -A
git commit -m "Release v0.2.0"

# 4. 创建标签
git tag v0.2.0

# 5. 推送
git push origin main
git push origin v0.2.0
```

### 获取帮助

如果问题仍未解决，请：

1. 检查 GitHub 仓库设置是否正确
2. 查看 [Go Modules 官方文档](https://golang.org/doc/modules/managing-dependencies)
3. 在 [GitHub Issues](https://github.com/lvdashuaibi/GPTUtils/issues) 中报告问题

### 相关资源

- [pkg.go.dev 文档](https://pkg.go.dev/about)
- [Go Modules 指南](https://golang.org/doc/modules/managing-dependencies)
- [GitHub 发布指南](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository)
