# feature: admin 模块单元测试

## 目标
- 覆盖率从 22.6% 提升到 50%
- 测试管理员核心逻辑

## 测试策略

### 方案：真实数据库测试

使用已有 TestMain，测试管理员功能。

## 测试用例

### 1. ListPosts 测试

#### 场景 1.1：获取帖子列表成功
- **Given** 数据库有 5 条帖子
- **When** 调用 ListPosts(1, 10, "active")
- **Then** 返回 5 条帖子

#### 场景 1.2：按状态筛选
- **Given** 数据库有 3 条活跃帖子，2 条已删除
- **When** 调用 ListPosts(1, 10, "deleted")
- **Then** 返回 2 条已删除帖子

#### 场景 1.3：分页查询
- **Given** 数据库有 15 条帖子
- **When** 调用 ListPosts(2, 10, "active")
- **Then** 返回第 2 页数据

#### 场景 1.4：默认参数
- **Given** page = 0, pageSize = 0
- **When** 调用 ListPosts(0, 0, "")
- **Then** 使用默认值 page=1, pageSize=20, status="active"

### 2. DeletePost 测试

#### 场景 2.1：删除帖子成功
- **Given** 帖子 ID = 1 存在
- **When** 调用 DeletePost(1)
- **Then** 返回 nil，帖子软删除

#### 场景 2.2：帖子不存在
- **Given** 帖子 ID = 999 不存在
- **When** 调用 DeletePost(999)
- **Then** 返回错误

### 3. ListUsers 测试

#### 场景 3.1：获取用户列表成功
- **Given** 数据库有 5 个用户
- **When** 调用 ListUsers(1, 10)
- **Then** 返回 5 个用户

#### 场景 3.2：分页查询
- **Given** 数据库有 15 个用户
- **When** 调用 ListUsers(2, 10)
- **Then** 返回第 2 页数据

### 4. UpdateUserStatus 测试

#### 场景 4.1：更新用户状态成功
- **Given** 用户 ID = 1 存在
- **When** 调用 UpdateUserStatus(1, 0)
- **Then** 返回 nil，用户状态更新为 0

#### 场景 4.2：用户不存在
- **Given** 用户 ID = 999 不存在
- **When** 调用 UpdateUserStatus(999, 0)
- **Then** 返回错误 "用户不存在"

### 5. GetStats 测试

#### 场景 5.1：获取统计数据成功
- **Given** 数据库有用户、帖子、评论
- **When** 调用 GetStats()
- **Then** 返回统计数据

### 6. SetPostTop 测试

#### 场景 6.1：设置置顶成功
- **Given** 帖子 ID = 1 存在
- **When** 调用 SetPostTop(1, true)
- **Then** 返回 nil，帖子置顶

### 7. SetPostFeatured 测试

#### 场景 7.1：设置加精成功
- **Given** 帖子 ID = 1 存在
- **When** 调用 SetPostFeatured(1, true)
- **Then** 返回 nil，帖子加精

## 预计覆盖率
- 函数覆盖：ListPosts, DeletePost, ListUsers, UpdateUserStatus, GetStats, SetPostTop, SetPostFeatured
- 预计覆盖率：50-60%
