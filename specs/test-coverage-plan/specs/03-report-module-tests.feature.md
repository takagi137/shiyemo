# feature: report 模块单元测试

## 目标
- 覆盖率从 1.5% 提升到 50%
- 测试举报核心逻辑

## 测试策略

### 方案：接口注入 + Mock

定义 `ReportRepository` 接口，注入 mock 实现，测试举报逻辑。

```go
// ReportRepository 举报仓储接口
type ReportRepository interface {
    ExistsByUserAndTarget(userID int64, targetType string, targetID int64) (bool, error)
    Create(report *Report) error
    List(page, size int) ([]Report, int64, error)
    UpdateStatus(reportID int64, status int) error
}
```

## 测试用例

### 1. Create 测试

#### 场景 1.1：创建举报成功
- **Given** 用户 ID = 1，举报帖子 ID = 100
- **When** 调用 Create(1, CreateRequest{TargetType: "post", TargetID: 100, Reason: "违规"})
- **Then** 返回 nil，举报创建成功

#### 场景 1.2：重复举报
- **Given** 用户已举报过帖子 ID = 100
- **When** 调用 Create(1, CreateRequest{TargetType: "post", TargetID: 100, Reason: "违规"})
- **Then** 返回错误 "你已经举报过该内容"

#### 场景 1.3：数据库查询错误
- **Given** ExistsByUserAndTarget 返回错误
- **When** 调用 Create(1, CreateRequest{...})
- **Then** 返回数据库错误

#### 场景 1.4：数据库创建错误
- **Given** Create 返回错误
- **When** 调用 Create(1, CreateRequest{...})
- **Then** 返回数据库错误

### 2. List 测试

#### 场景 2.1：获取举报列表成功
- **Given** 数据库有 5 条举报
- **When** 调用 List(1, 10)
- **Then** 返回 5 条举报，total = 5

#### 场景 2.2：举报列表为空
- **Given** 数据库无举报
- **When** 调用 List(1, 10)
- **Then** 返回空列表，total = 0

#### 场景 2.3：分页查询
- **Given** 数据库有 15 条举报
- **When** 调用 List(2, 10)
- **Then** 返回第 2 页数据

### 3. Handle 测试

#### 场景 3.1：处理举报成功
- **Given** 举报 ID = 1 存在
- **When** 调用 Handle(1, 2)
- **Then** 返回 nil，状态更新为 2

#### 场景 3.2：举报不存在
- **Given** 举报 ID = 999 不存在
- **When** 调用 Handle(999, 2)
- **Then** 返回错误

## 预计覆盖率
- 函数覆盖：Create, List, Handle
- 预计覆盖率：50-60%
