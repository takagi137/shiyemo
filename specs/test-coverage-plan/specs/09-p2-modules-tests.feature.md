# feature: P2 模块单元测试

## 目标
- 将 P2 模块覆盖率提升到 50%

## 模块列表

### 1. comment 模块（38.8% → 50%）

#### 测试用例

##### 场景 1.1：空评论列表
- **Given** 帖子 ID = 1 无评论
- **When** 调用 List(1, 1, 10)
- **Then** 返回空列表，total = 0

##### 场景 1.2：帖子不存在时创建评论
- **Given** 帖子 ID = 999 不存在
- **When** 调用 Create(1, 999, CreateRequest{Content: "评论"})
- **Then** 返回错误 ErrPostNotFound

##### 场景 1.3：大量评论分页
- **Given** 帖子有 100 条评论
- **When** 调用 List(1, 5, 20)
- **Then** 返回第 5 页数据

### 2. post 模块（39.6% → 50%）

#### 测试用例

##### 场景 2.1：FollowingFeed 成功
- **Given** 用户关注了 3 个用户，他们有 5 条帖子
- **When** 调用 FollowingFeed(userID, 1, 10)
- **Then** 返回 5 条帖子

##### 场景 2.2：FollowingFeed 无关注
- **Given** 用户未关注任何人
- **When** 调用 FollowingFeed(userID, 1, 10)
- **Then** 返回空列表

##### 场景 2.3：缓存命中
- **Given** 帖子详情已缓存
- **When** 调用 GetDetail(postID)
- **Then** 从缓存返回，不查询数据库

### 3. moderation 模块（29.8% → 50%）

#### 测试用例

##### 场景 3.1：CheckTitle 成功
- **Given** 标题 "正常标题"
- **When** 调用 CheckTitle("正常标题")
- **Then** 返回 Passed = true

##### 场景 3.2：CheckTitle 命中敏感词
- **Given** 标题包含敏感词
- **When** 调用 CheckTitle("色情内容")
- **Then** 返回 Passed = false

##### 场景 3.3：AddWord 添加敏感词
- **Given** 敏感词库为空
- **When** 调用 AddWord("新敏感词")
- **Then** 敏感词库包含 "新敏感词"

##### 场景 3.4：normalize 文本预处理
- **Given** 文本 "Hello World"
- **When** 调用 normalize("Hello World")
- **Then** 返回 "helloworld"

##### 场景 3.5：contains 变体绕过检测
- **Given** 文本 "色*情*内*容"
- **When** 调用 contains("色*情*内*容", "色情")
- **Then** 返回 true

### 4. upload 模块（31.5% → 50%）

#### 测试用例

##### 场景 4.1：不支持的文件格式
- **Given** 文件名 "test.txt"
- **When** 调用 Upload(1, file, "test.txt", "image")
- **Then** 返回错误 "不支持的文件格式"

##### 场景 4.2：文件大小超限
- **Given** 文件大小 11MB
- **When** 调用 Upload(1, file, "test.jpg", "image")
- **Then** 返回错误 "文件大小超过限制"

##### 场景 4.3：存储后端错误
- **Given** 存储后端返回错误
- **When** 调用 Upload(1, file, "test.jpg", "image")
- **Then** 返回错误

### 5. user 模块（39.5% → 50%）

#### 测试用例

##### 场景 5.1：UploadAvatar 成功
- **Given** 用户 ID = 1，文件 "avatar.jpg"
- **When** 调用 UploadAvatar(1, file, "avatar.jpg")
- **Then** 返回头像 URL

##### 场景 5.2：UploadAvatar 格式错误
- **Given** 文件名 "avatar.txt"
- **When** 调用 UploadAvatar(1, file, "avatar.txt")
- **Then** 返回错误 ErrParamInvalid

##### 场景 5.3：缓存命中
- **Given** 用户详情已缓存
- **When** 调用 GetProfile(userID)
- **Then** 从缓存返回，不查询数据库

### 6. websocket 模块（32.2% → 50%）

#### 测试用例

##### 场景 6.1：SendToUser 成功
- **Given** 用户 ID = 1 在线
- **When** 调用 SendToUser(1, "notification", content)
- **Then** 用户收到消息

##### 场景 6.2：Broadcast 成功
- **Given** 3 个用户在线
- **When** 调用 Broadcast("announcement", content)
- **Then** 所有用户收到消息

##### 场景 6.3：GetOnlineCount 成功
- **Given** 3 个用户在线
- **When** 调用 GetOnlineCount()
- **Then** 返回 3

##### 场景 6.4：IsUserOnline 成功
- **Given** 用户 ID = 1 在线
- **When** 调用 IsUserOnline(1)
- **Then** 返回 true

##### 场景 6.5：IsUserOnline 用户不在线
- **Given** 用户 ID = 2 不在线
- **When** 调用 IsUserOnline(2)
- **Then** 返回 false

## 预计覆盖率
- 每个模块预计覆盖率：50-60%
