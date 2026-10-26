# feature: search 模块单元测试

## 目标
- 覆盖率从 0% 提升到 50%
- 测试核心搜索逻辑，mock Meilisearch 客户端

## 测试策略

### 方案：接口注入 + Mock

定义 `SearchClient` 接口，注入 mock 实现，测试搜索逻辑。

```go
// SearchClient 搜索客户端接口
type SearchClient interface {
    Search(index, query string, request *SearchRequest) (*SearchResult, error)
    AddDocuments(index string, documents []map[string]interface{}) error
}
```

## 测试用例

### 1. SearchPosts 测试

#### 场景 1.1：搜索帖子成功
- **Given** Meilisearch 返回 3 条帖子
- **When** 调用 SearchPosts("Go", 1, 10)
- **Then** 返回 3 条帖子，total = 3

#### 场景 1.2：搜索结果为空
- **Given** Meilisearch 返回空结果
- **When** 调用 SearchPosts("不存在", 1, 10)
- **Then** 返回空列表，total = 0

#### 场景 1.3：搜索分页
- **Given** Meilisearch 返回第 2 页结果
- **When** 调用 SearchPosts("Go", 2, 10)
- **Then** 返回第 2 页数据

#### 场景 1.4：搜索错误
- **Given** Meilisearch 返回错误
- **When** 调用 SearchPosts("Go", 1, 10)
- **Then** 返回错误

### 2. SearchUsers 测试

#### 场景 2.1：搜索用户成功
- **Given** Meilisearch 返回 2 个用户
- **When** 调用 SearchUsers("张三", 1, 10)
- **Then** 返回 2 个用户，total = 2

#### 场景 2.2：搜索用户为空
- **Given** Meilisearch 返回空结果
- **When** 调用 SearchUsers("不存在", 1, 10)
- **Then** 返回空列表，total = 0

### 3. GetSuggestions 测试

#### 场景 3.1：获取搜索建议成功
- **Given** Meilisearch 返回 3 条建议
- **When** 调用 GetSuggestions("Go", 5)
- **Then** 返回 3 条去重建议

#### 场景 3.2：搜索建议为空
- **Given** Meilisearch 返回空结果
- **When** 调用 GetSuggestions("不存在", 5)
- **Then** 返回空列表

### 4. SyncAllPosts 测试

#### 场景 4.1：同步帖子成功
- **Given** 数据库有 5 条帖子，Meilisearch 同步成功
- **When** 调用 SyncAllPosts()
- **Then** 返回 nil，Meilisearch 收到 5 条文档

#### 场景 4.2：数据库查询失败
- **Given** 数据库查询返回错误
- **When** 调用 SyncAllPosts()
- **Then** 返回错误

### 5. SyncAllUsers 测试

#### 场景 5.1：同步用户成功
- **Given** 数据库有 3 个用户，Meilisearch 同步成功
- **When** 调用 SyncAllUsers()
- **Then** 返回 nil，Meilisearch 收到 3 条文档

## 预计覆盖率
- 函数覆盖：SearchPosts, SearchUsers, GetSuggestions, SyncAllPosts, SyncAllUsers
- 预计覆盖率：50-60%
