# feature: notification 模块单元测试

## 目标
- 覆盖率从 23.4% 提升到 50%
- 测试通知核心逻辑

## 测试策略

### 方案：真实数据库测试

使用已有 TestMain，测试通知功能。

## 测试用例

### 1. ListByType 测试

#### 场景 1.1：按类型筛选成功
- **Given** 用户有 3 条点赞通知，2 条评论通知
- **When** 调用 ListByType(userID, "like", false, 1, 10)
- **Then** 返回 3 条点赞通知

#### 场景 1.2：筛选未读通知
- **Given** 用户有 5 条通知，3 条未读
- **When** 调用 ListByType(userID, "", true, 1, 10)
- **Then** 返回 3 条未读通知

#### 场景 1.3：类型 + 未读组合筛选
- **Given** 用户有 3 条点赞通知，2 条未读
- **When** 调用 ListByType(userID, "like", true, 1, 10)
- **Then** 返回 2 条未读点赞通知

### 2. BatchMarkAsRead 测试

#### 场景 2.1：批量标记已读成功
- **Given** 用户有 5 条未读通知
- **When** 调用 BatchMarkAsRead([1,2,3], userID)
- **Then** 返回 nil，未读数量变为 2

#### 场景 2.2：空 ID 列表
- **Given** 空 ID 列表
- **When** 调用 BatchMarkAsRead([], userID)
- **Then** 返回错误 "通知ID列表不能为空"

### 3. BatchDelete 测试

#### 场景 3.1：批量删除成功
- **Given** 用户有 5 条通知
- **When** 调用 BatchDelete([1,2], userID)
- **Then** 返回 nil，剩余 3 条通知

#### 场景 3.2：空 ID 列表
- **Given** 空 ID 列表
- **When** 调用 BatchDelete([], userID)
- **Then** 返回错误 "通知ID列表不能为空"

#### 场景 3.3：删除他人通知
- **Given** 通知 ID = 1 属于用户 A
- **When** 调用 BatchDelete([1], userB_ID)
- **Then** 返回错误（无权删除）

### 4. GetDetail 测试

#### 场景 4.1：获取通知详情成功
- **Given** 通知 ID = 1 存在
- **When** 调用 GetDetail(1)
- **Then** 返回通知详情

#### 场景 4.2：通知不存在
- **Given** 通知 ID = 999 不存在
- **When** 调用 GetDetail(999)
- **Then** 返回错误

### 5. GetUnreadCount 测试

#### 场景 5.1：获取未读数量成功
- **Given** 用户有 3 条未读通知
- **When** 调用 GetUnreadCount(userID)
- **Then** 返回 3

## 预计覆盖率
- 函数覆盖：ListByType, BatchMarkAsRead, BatchDelete, GetDetail, GetUnreadCount
- 预计覆盖率：50-60%
