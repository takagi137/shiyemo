# feature: tag 模块单元测试

## 目标
- 覆盖率从 0.5% 提升到 50%
- 测试标签 CRUD 核心逻辑

## 测试策略

### 方案：接口注入 + Mock

定义 `TagRepository` 接口，注入 mock 实现，测试标签逻辑。

```go
// TagRepository 标签仓储接口
type TagRepository interface {
    GetHotTags(limit int) ([]map[string]interface{}, error)
    GetOrCreate(name string) (*Tag, error)
    UpdatePostTags(postID int64, tagIDs []int64) error
    GetPostTags(postID int64) ([]string, error)
}
```

## 测试用例

### 1. GetHotTags 测试

#### 场景 1.1：获取热门标签成功
- **Given** 数据库有 5 个标签
- **When** 调用 GetHotTags(10)
- **Then** 返回 5 个标签

#### 场景 1.2：获取热门标签（限制数量）
- **Given** 数据库有 10 个标签
- **When** 调用 GetHotTags(5)
- **Then** 返回 5 个标签

#### 场景 1.3：默认限制数量
- **Given** limit = 0
- **When** 调用 GetHotTags(0)
- **Then** 使用默认值 20

#### 场景 1.4：数据库错误
- **Given** 数据库查询返回错误
- **When** 调用 GetHotTags(10)
- **Then** 返回错误

### 2. UpdatePostTags 测试

#### 场景 2.1：更新帖子标签成功
- **Given** 帖子 ID = 1，标签 = ["Go", "后端"]
- **When** 调用 UpdatePostTags(1, ["Go", "后端"])
- **Then** 返回标签列表，数据库更新成功

#### 场景 2.2：标签不存在时自动创建
- **Given** 标签 "新标签" 不存在
- **When** 调用 UpdatePostTags(1, ["新标签"])
- **Then** 标签自动创建，关联成功

#### 场景 2.3：数据库错误
- **Given** 数据库更新返回错误
- **When** 调用 UpdatePostTags(1, ["Go"])
- **Then** 返回错误

### 3. GetPostTags 测试

#### 场景 3.1：获取帖子标签成功
- **Given** 帖子 ID = 1 有 3 个标签
- **When** 调用 GetPostTags(1)
- **Then** 返回 3 个标签名称

#### 场景 3.2：帖子无标签
- **Given** 帖子 ID = 2 无标签
- **When** 调用 GetPostTags(2)
- **Then** 返回空列表

#### 场景 3.3：数据库错误
- **Given** 数据库查询返回错误
- **When** 调用 GetPostTags(1)
- **Then** 返回错误

## 预计覆盖率
- 函数覆盖：GetHotTags, UpdatePostTags, GetPostTags
- 预计覆盖率：50-60%
