# feature: follow 模块单元测试

## 目标
- 覆盖率从 25.2% 提升到 50%
- 测试关注核心逻辑

## 测试策略

### 方案：Mock 依赖注入

使用 `NewFollowServiceWithDeps` 注入 mock 依赖，测试关注逻辑。

## 测试用例

### 1. ListFollowing 测试

#### 场景 1.1：获取关注列表成功
- **Given** 用户 ID = 1 关注了 3 个用户
- **When** 调用 ListFollowing(1, 1, 10)
- **Then** 返回 3 个用户信息，total = 3

#### 场景 1.2：关注列表为空
- **Given** 用户 ID = 1 未关注任何人
- **When** 调用 ListFollowing(1, 1, 10)
- **Then** 返回空列表，total = 0

#### 场景 1.3：分页查询
- **Given** 用户 ID = 1 关注了 15 个用户
- **When** 调用 ListFollowing(1, 2, 10)
- **Then** 返回第 2 页数据

#### 场景 1.4：用户信息获取失败
- **Given** 关注列表中有用户信息获取失败
- **When** 调用 ListFollowing(1, 1, 10)
- **Then** 跳过失败用户，返回其他用户

### 2. ListFollowers 测试

#### 场景 2.1：获取粉丝列表成功
- **Given** 用户 ID = 1 有 3 个粉丝
- **When** 调用 ListFollowers(1, 1, 10)
- **Then** 返回 3 个粉丝信息，total = 3

#### 场景 2.2：粉丝列表为空
- **Given** 用户 ID = 1 无粉丝
- **When** 调用 ListFollowers(1, 1, 10)
- **Then** 返回空列表，total = 0

#### 场景 2.3：分页查询
- **Given** 用户 ID = 1 有 15 个粉丝
- **When** 调用 ListFollowers(1, 2, 10)
- **Then** 返回第 2 页数据

## 预计覆盖率
- 函数覆盖：ListFollowing, ListFollowers
- 预计覆盖率：50-60%
