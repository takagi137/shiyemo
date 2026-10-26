# feature: interaction 模块单元测试

## 目标
- 覆盖率从 25.9% 提升到 50%
- 测试互动核心逻辑

## 测试策略

### 方案：Mock 依赖注入

使用 `NewInteractionServiceWithDeps` 注入 mock 依赖，测试互动逻辑。

## 测试用例

### 1. ListBookmarks 测试

#### 场景 1.1：获取收藏列表成功
- **Given** 用户 ID = 1 收藏了 3 个帖子
- **When** 调用 ListBookmarks(1, 1, 10)
- **Then** 返回 3 个帖子信息，total = 3

#### 场景 1.2：收藏列表为空
- **Given** 用户 ID = 1 无收藏
- **When** 调用 ListBookmarks(1, 1, 10)
- **Then** 返回空列表，total = 0

#### 场景 1.3：分页查询
- **Given** 用户 ID = 1 收藏了 15 个帖子
- **When** 调用 ListBookmarks(1, 2, 10)
- **Then** 返回第 2 页数据

#### 场景 1.4：帖子信息获取失败
- **Given** 收藏列表中有帖子信息获取失败
- **When** 调用 ListBookmarks(1, 1, 10)
- **Then** 跳过失败帖子，返回其他帖子

## 预计覆盖率
- 函数覆盖：ListBookmarks
- 预计覆盖率：50-60%
