# tasks：二次开发计划 - 任务清单

## Phase 1: 用户主页数据

### 数据层
- [x] 无改动（复用现有查询）

### 业务层
- [x] user/service.go 添加 GetUserStats 方法

### 接口层
- [x] user/handler.go 扩展 GetProfile 接口

### 验证
- [x] L1: 单元测试通过
- [x] L3: BDD 场景通过

---

## Phase 2: 通知分类 + 批量操作

### 数据层
- [x] notification/model.go 添加 type 字段枚举
- [x] 执行数据库迁移

### 业务层
- [x] notification/repo.go 添加批量操作方法
- [x] notification/service.go 添加分类查询、批量操作

### 接口层
- [x] notification/handler.go 添加批量 API

### 验证
- [x] L1: 单元测试通过
- [x] L3: BDD 场景通过

---

## Phase 3: 评论回复通知

### 数据层
- [x] 无改动

### 业务层
- [x] comment/service.go 添加回复通知逻辑

### 接口层
- [x] 无改动（通知自动触发）

### 验证
- [x] L1: 单元测试通过
- [x] L3: BDD 场景通过

---

## Phase 4: 搜索增强

### 数据层
- [x] 无改动

### 业务层
- [x] search/service.go 添加用户搜索方法
- [x] search/service.go 添加标签搜索方法
- [x] search/service.go 添加搜索建议方法

### 接口层
- [x] search/handler.go 添加搜索建议接口

### 验证
- [x] L1: 单元测试通过
- [x] L3: BDD 场景通过

---

## Phase 5: 帖子编辑增强

### 数据层
- [x] 执行数据库迁移（content 改为 JSONB）

### 业务层
- [x] post/service.go 支持富文本存储
- [x] upload/service.go 实现图片上传

### 接口层
- [x] upload/handler.go 添加图片上传接口

### 验证
- [x] L1: 单元测试通过
- [x] L3: BDD 场景通过

---

## Phase 6: 内容审核队列

### 数据层
- [x] 执行数据库迁移（添加审核字段）

### 业务层
- [x] admin/service.go 添加审核队列查询
- [x] admin/service.go 添加审核操作（approve/reject）

### 接口层
- [x] admin/handler.go 添加审核 API

### 验证
- [x] L1: 单元测试通过
- [x] L3: BDD 场景通过

> 状态同步：2026-05-29，代码已全部实现
