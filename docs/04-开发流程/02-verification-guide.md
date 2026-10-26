> 开源副本：部署地址已替换为示例；历史验证记录不代表当前部署状态。

# 验证优先开发指南

> 本文档定义"验证优先"原则的具体执行规范。
> **验证优先 = 先写验证，再写代码；验证不过，代码不入库。**

---

## 🎯 核心理念

```
传统模式：写代码 → 测试 → 发现问题 → 修复
验证优先：写验证 → 写代码 → 验证通过 → 完成
```

**好处**：
1. 需求前置澄清（写验证时发现需求不清晰）
2. 设计前置验证（写验证时发现架构问题）
3. 质量内建（不是事后补测试）

---

## 📋 验证层次模型

### L1 - 单元测试（Unit Tests）

**目标**：验证每个函数/方法的行为正确性

**覆盖范围**：
- service 层的所有公开方法
- repo 层的所有公开方法
- 工具函数（jwt, password, validator 等）

**测试框架**：
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

**测试文件命名**：
- 源文件：`service.go`
- 测试文件：`service_test.go`

**测试函数命名**：
```go
func TestAuthService_Register(t *testing.T) { ... }
func TestAuthService_Register_DuplicateEmail(t *testing.T) { ... }
```

**每个函数至少 3 个测试用例**：
| 用例类型 | 说明 | 示例 |
|---------|------|------|
| Happy Path | 正常流程 | 正常注册成功 |
| Edge Case | 边界情况 | 邮箱已存在 |
| Error Case | 错误场景 | 密码格式错误 |

**测试结构**：
```go
func TestAuthService_Register(t *testing.T) {
    // Arrange - 准备测试数据
    req := RegisterRequest{
        Email:    "test@example.com",
        Password: "a123456",
        Nickname: "测试用户",
    }
    
    // Act - 执行被测函数
    resp, err := service.Register(req)
    
    // Assert - 验证结果
    require.NoError(t, err)
    assert.Equal(t, req.Email, resp.Email)
}
```

**Mock 策略**：
- ✅ 可以 Mock：外部服务（SMS、邮件）
- ❌ 不可 Mock：数据库（使用测试数据库）
- ❌ 不可 Mock：核心业务逻辑

---

### L2 - 集成测试（Integration Tests）

**目标**：验证 API 端点的完整行为

**覆盖范围**：
- 每个 API 端点的正常/异常场景
- 认证/授权行为
- 数据库交互

**测试文件位置**：
```
tests/
├── integration/
│   ├── auth_test.go
│   ├── post_test.go
│   └── ...
```

**测试工具**：
```go
import (
    "net/http"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
)
```

**测试结构**：
```go
func TestRegisterByEmail(t *testing.T) {
    // Setup - 初始化测试环境
    router := setupTestRouter()
    db := setupTestDB()
    
    t.Run("success", func(t *testing.T) {
        // 准备请求
        body := `{"email":"test@example.com","password":"a123456","nickname":"Test"}`
        req := httptest.NewRequest("POST", "/api/v1/auth/register/email", 
            strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        
        // 执行请求
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        // 验证响应
        assert.Equal(t, 200, w.Code)
        // 验证数据库
        var user User
        db.Where("email = ?", "test@example.com").First(&user)
        assert.Equal(t, "Test", user.Nickname)
    })
    
    t.Run("duplicate_email", func(t *testing.T) { ... })
    t.Run("missing_fields", func(t *testing.T) { ... })
}
```

**测试数据库**：
- 使用独立的测试数据库（`shiyemo_test`）
- 每个测试前清空相关表
- 使用事务回滚

---

### L3 - 部署验证（Deployment Verification）

**目标**：验证部署后系统整体行为正确

**验证内容**：
1. 健康检查
2. BDD 场景执行
3. 关键路径手动验证

**验证脚本**：
```bash
#!/bin/bash
# scripts/verify-bdd.sh

BASE_URL="http://app.example.com:8080"

echo "=== 验证用户注册场景 ==="
# Given: 用户未注册
# When: 提交注册请求
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register/email" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com","password":"a123456","nickname":"NewUser"}')
# Then: 注册成功
echo "$RESPONSE" | jq '.code' | grep -q 0 && echo "✅ 注册成功" || echo "❌ 注册失败"

echo "=== 验证登录场景 ==="
# Given: 用户已注册
# When: 提交登录请求
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login/email" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com","password":"a123456"}')
# Then: 登录成功，返回 token
echo "$RESPONSE" | jq '.data.token' | grep -q "." && echo "✅ 登录成功" || echo "❌ 登录失败"
```

**BDD 场景文件**：
- 位置：`specs/*/specs/*.feature.md`
- 格式：Given/When/Then
- 必须与代码同步更新

---

## 📝 验证检查清单

### 开发前

- [ ] 阅读并理解 BDD 场景文件
- [ ] 确认测试环境可用
- [ ] 准备测试数据

### 开发中

- [ ] 每个函数写完立即写测试
- [ ] 测试通过后再继续下一个函数
- [ ] 使用 `go test -v` 查看详细输出

### 开发后

- [ ] 运行所有单元测试：`go test ./...`
- [ ] 检查测试覆盖率：`go test -cover ./...`
- [ ] 运行集成测试
- [ ] 部署到验证环境
- [ ] 执行 BDD 场景验证
- [ ] 记录验证结果

---

## 🚫 常见验证错误

| 错误模式 | 正确做法 |
|---------|---------|
| Mock 数据库 | 使用真实测试数据库 |
| 只写 happy path | 至少 3 个测试用例 |
| 测试后不清理 | 使用事务回滚或清理脚本 |
| 测试依赖顺序 | 每个测试独立，可单独运行 |
| 只测 HTTP 状态码 | 验证响应体和数据库状态 |
| 部署后不验证 | 必须运行 BDD 场景 |

---

## 📊 验证覆盖率报告

每次提交前，必须生成覆盖率报告：

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html
```

**覆盖率要求**：
- 整体覆盖率：≥ 70%
- service 层：≥ 80%
- 新增代码：≥ 90%

---

## 🔄 验证与 OpenSpec 的关系

```
OpenSpec 流程：
  proposal.md → design.md → specs/*.feature.md → tasks.md

验证流程：
  specs/*.feature.md → L1 单元测试 → L2 集成测试 → L3 部署验证
```

**关键规则**：
1. BDD 场景文件是验证的**唯一标准**
2. 验证结果必须与 BDD 场景一一对应
3. 验证不通过，不得标记任务完成

---

## 📚 参考文档

| 文档 | 路径 |
|------|------|
| 开发约束 | `docs/02-开发规范/02-ai-development-guide.md` |
| OpenSpec 流程 | `docs/04-开发流程/00-openspec-guide.md` |
| BDD 指南 | `docs/04-开发流程/01-bdd-guide.md` |
| 已有场景 | `specs/*/specs/*.feature.md` |

---

*本文档是验证优先原则的核心指南，所有开发活动必须遵循。*