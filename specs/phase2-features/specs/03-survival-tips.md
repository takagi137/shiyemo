# spec：生存小贴士

## 场景 1：获取随机贴士

### Given 前置条件
- 数据库中有 20+ 条贴士（is_active = true）

### When 操作
- 调用 `GET /api/v1/tips/random`（无需登录）

### Then 预期结果
- 响应状态码 200，code=0
- 返回一条贴士：`{ "id": N, "content": "...", "category": "finance" }`
- 同一用户同一天多次调用返回同一条（基于 cookie/日期种子）
- 不同用户或不同天返回不同条

---

## 场景 2：刷新获取另一条

### Given 前置条件
- 用户已看到当日贴士

### When 操作
- 调用 `GET /api/v1/tips/random?refresh=true`

### Then 预期结果
- 返回与当前不同的一条贴士

---

## 场景 3：管理员添加贴士

### Given 前置条件
- 用户为管理员（is_admin = true）

### When 操作
- 调用 `POST /api/v1/tips`，body: `{ "content": "新贴士", "category": "life" }`

### Then 预期结果
- 响应状态码 200，code=0
- 数据库新增一条贴士记录
- 返回 `{ "id": N }`

---

## 场景 4：管理员查看贴士列表

### Given 前置条件
- 用户为管理员

### When 操作
- 调用 `GET /api/v1/tips?page=1&size=20`

### Then 预期结果
- 返回分页贴士列表
- 格式：`{ "list": [...], "page": 1, "size": 20, "has_more": false }`

---

## 场景 5：贴士库为空时

### Given 前置条件
- 数据库无活跃贴士

### When 操作
- 调用 `GET /api/v1/tips/random`

### Then 预期结果
- 返回默认兜底贴士："今天也是自由的一天。"
