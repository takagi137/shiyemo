# spec：L2 集成测试框架建立

## 场景描述

建立集成测试基础设施，验证核心 API 端到端功能。

### Given 前置条件
- 测试数据库可连接（使用独立测试 DB）
- 后端服务可通过 httptest 启动

### When 操作
- 创建 `tests/integration/` 目录
- 实现测试 DB 初始化、数据清理、HTTP 请求工具
- 编写 auth 模块集成测试：注册 → 登录 → 获取 Profile
- 编写 post 模块集成测试：创建 → 查询 → 更新 → 删除

### Then 预期结果
- `go test ./tests/integration/...` 全部通过
- auth 场景：注册返回 200 + token → 登录返回 200 + token → Profile 返回用户信息
- post 场景：创建返回 201 → 列表包含新帖 → 更新成功 → 删除后列表不包含
