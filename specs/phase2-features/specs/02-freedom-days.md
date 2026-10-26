# spec：自由天数计时器

## 场景 1：设置最后工作日

### Given 前置条件
- 用户已登录

### When 操作
- 用户调用 `PUT /api/v1/user/profile`，body: `{ "last_work_date": "2026-03-15" }`

### Then 预期结果
- 响应状态码 200，code=0
- 返回用户数据中包含 `last_work_date: "2026-03-15"`
- 数据库 users.last_work_date 已更新

---

## 场景 2：查看用户详情包含自由天数

### Given 前置条件
- 用户 A 设置了 last_work_date = "2026-03-15"
- 当前日期为 2026-05-27

### When 操作
- 任何人调用 `GET /api/v1/users/:id`

### Then 预期结果
- 响应中包含 `last_work_date: "2026-03-15"`
- 前端根据该字段计算展示"自由第 73 天"

---

## 场景 3：未设置最后工作日

### Given 前置条件
- 用户未设置 last_work_date（字段为 null）

### When 操作
- 查看该用户详情

### Then 预期结果
- 响应中 `last_work_date` 为 null
- 前端不展示自由天数卡片

---

## 场景 4：清除最后工作日

### Given 前置条件
- 用户已设置 last_work_date

### When 操作
- 用户调用 `PUT /api/v1/user/profile`，body: `{ "last_work_date": null }`

### Then 预期结果
- 响应状态码 200
- last_work_date 设为 null
- 前端不再展示自由天数
