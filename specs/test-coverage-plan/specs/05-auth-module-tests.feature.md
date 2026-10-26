# feature: auth 模块单元测试

## 目标
- 覆盖率从 29.7% 提升到 50%
- 测试认证核心逻辑

## 测试策略

### 方案：Mock 依赖注入

使用 `NewAuthServiceWithDeps` 注入 mock repo，测试认证逻辑。

## 测试用例

### 1. RegisterByPhone 测试

#### 场景 1.1：手机号注册成功
- **Given** 手机号 "13800138000" 未注册，验证码 "123456"
- **When** 调用 RegisterByPhone("13800138000", "123456", "新用户")
- **Then** 返回 token 和用户信息

#### 场景 1.2：手机号已注册
- **Given** 手机号 "13800138000" 已注册
- **When** 调用 RegisterByPhone("13800138000", "123456", "新用户")
- **Then** 返回错误 ErrPhoneExists

#### 场景 1.3：验证码错误
- **Given** 验证码 "000000"
- **When** 调用 RegisterByPhone("13800138000", "000000", "新用户")
- **Then** 返回错误 ErrCodeWrong

### 2. LoginByPhone 测试

#### 场景 2.1：手机号登录成功
- **Given** 手机号 "13800138000" 已注册，验证码 "123456"
- **When** 调用 LoginByPhone("13800138000", "123456")
- **Then** 返回 token 和用户信息

#### 场景 2.2：手机号未注册
- **Given** 手机号 "13800138001" 未注册
- **When** 调用 LoginByPhone("13800138001", "123456")
- **Then** 返回错误 ErrUserNotFound

#### 场景 2.3：验证码错误
- **Given** 验证码 "000000"
- **When** 调用 LoginByPhone("13800138000", "000000")
- **Then** 返回错误 ErrCodeWrong

#### 场景 2.4：用户被禁用
- **Given** 用户状态 = 0（禁用）
- **When** 调用 LoginByPhone("13800138000", "123456")
- **Then** 返回错误 ErrUserDisabled

### 3. SendSMS 测试

#### 场景 3.1：发送验证码成功
- **Given** 手机号 "13800138000"
- **When** 调用 SendSMS("13800138000")
- **Then** 返回 nil（MVP 阶段 mock）

### 4. verifyCode 测试

#### 场景 4.1：验证码正确（debug 模式）
- **Given** debug 模式，验证码 "123456"
- **When** 调用 verifyCode("13800138000", "123456")
- **Then** 返回 nil

#### 场景 4.2：验证码错误（debug 模式）
- **Given** debug 模式，验证码 "000000"
- **When** 调用 verifyCode("13800138000", "000000")
- **Then** 返回错误 ErrCodeWrong

## 预计覆盖率
- 函数覆盖：RegisterByPhone, LoginByPhone, SendSMS, verifyCode
- 预计覆盖率：50-60%
