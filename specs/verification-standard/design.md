# design：验证标准

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 测试框架 | Go testing + testify | Go 标准 + 强大断言 |
| Mock 策略 | 只 Mock 外部服务 | 数据库使用真实测试库 |
| 测试数据 | 每个测试独立 | 使用事务回滚 |
| 覆盖率工具 | go test -cover | Go 内置 |

## 验证架构

```
L1 单元测试 (Unit Tests)
├── service_test.go  ← 业务逻辑测试
├── repo_test.go     ← 数据访问测试
└── *_test.go        ← 工具函数测试

L2 集成测试 (Integration Tests)
└── tests/integration/
    ├── auth_test.go
    └── post_test.go

L3 部署验证 (Deployment Verification)
└── scripts/verify-bdd.sh  ← BDD 场景自动化
```

## 测试文件结构

```
server/
├── internal/
│   ├── auth/
│   │   ├── service.go
│   │   ├── service_test.go    ← L1 单元测试
│   │   ├── repo.go
│   │   └── repo_test.go       ← L1 单元测试
│   └── post/
│       ├── service.go
│       └── service_test.go    ← L1 单元测试
├── tests/
│   └── integration/           ← L2 集成测试
│       ├── auth_test.go
│       └── post_test.go
└── testutil/                  ← 测试工具函数
    └── testutil.go
```

## 测试用例设计

### 每个函数至少 3 个测试用例

| 用例类型 | 说明 | 示例 |
|---------|------|------|
| Happy Path | 正常流程 | 注册成功 |
| Edge Case | 边界情况 | 邮箱已存在 |
| Error Case | 错误场景 | 密码格式错误 |

### 测试结构

```go
func TestAuthService_Register(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        // Arrange
        // Act
        // Assert
    })
    
    t.Run("duplicate_email", func(t *testing.T) { ... })
    t.Run("invalid_password", func(t *testing.T) { ... })
}
```

## 测试数据库配置

- 数据库名：`shiyemo_test`
- 每个测试前：清空相关表或使用事务
- 每个测试后：回滚事务

## 依赖

- Go testing 包
- testify（assert, require）
- 测试数据库

## 风险点

1. 测试数据库需要手动创建
2. 测试数据需要隔离
3. 并发测试需要注意锁
