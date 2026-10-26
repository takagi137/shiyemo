# tasks：单元测试覆盖率提升

> 覆盖率标准：只统计 service 层核心业务逻辑，详见 `docs/04-开发流程/03-test-coverage-standard.md`

## 任务清单

### Phase 1: P0 模块（接口化 + mock）✅ 已完成

#### 任务 1.1: search 模块接口化 ✅
- **状态**: ✅ 完成
- **产出**: `search_client.go` 接口 + `mock_test.go` + `service_test.go`（7 个测试）
- **service 覆盖率**: 67.0%

#### 任务 1.2: tag 模块接口化 ✅
- **状态**: ✅ 完成
- **产出**: `tag_repo.go` 接口 + mock + `service_test.go`（11 个测试）
- **service 覆盖率**: 33.9%（core service 100%，skill_service 0%）
- **额外修复**: service 返回类型 `*bizerr.BizError` → `error`，与其他模块统一

#### 任务 1.3: report 模块接口化 ✅
- **状态**: ✅ 完成
- **产出**: `report_repo.go` 接口 + mock + `service_test.go`（12 个测试）
- **service 覆盖率**: 84.6%

### Phase 2: P1 模块补充 ✅ 已完成

#### 任务 2.1: auth 模块补充 ✅
- **状态**: ✅ 完成
- **产出**: 新增 LoginByPhone/user_disabled 测试
- **service 覆盖率**: 83.9%

#### 任务 2.2: 其他 P1 模块 ✅
- **状态**: ✅ 完成
- **产出**: 修复 notification ListByType 参数、user 缺少 fmt import
- **follow**: 81.8% | **interaction**: 83.3% | **comment**: 84.6%

### Phase 3: 待补充模块 ⏳

#### 任务 3.1: tag/skill_service.go 测试
- **状态**: ✅ 完成
- **优先级**: P1
- **完成日期**: 2026-05-21
- **产出**: skill_service.go 函数覆盖率 83-100%（AddUserSkill 83.3%, SearchBySkill 88.9%, 其余 100%）

#### 任务 3.2: moderation/admin_service.go 测试
- **状态**: ✅ 完成
- **优先级**: P2
- **完成日期**: 2026-05-21
- **产出**: admin_service.go 所有函数 100% 覆盖（AddWord/DeleteWord/ListWords）

#### 任务 3.3: admin/notification 远程验证
- **状态**: ⏳ 待开始
- **优先级**: P2
- **预计工时**: 1h
- **依赖**: 远程 DB 环境

---

## 当前覆盖率总览（service 层口径）

| 模块 | service 覆盖率 | 状态 |
|------|---------------|------|
| upload | 88.5% | ✅ |
| export | 87.5% | ✅ |
| comment | 84.6% | ✅ |
| report | 84.6% | ✅ |
| auth | 83.9% | ✅ |
| interaction | 83.3% | ✅ |
| follow | 81.8% | ✅ |
| websocket | 70.7% | ✅ |
| search | 67.0% | ✅ |
| moderation | 62.8% | ⚠️ admin_service 0% |
| tag | 33.9% | ⚠️ skill_service 0% |
| admin | N/A | ⏳ 需 DB |
| notification | N/A | ⏳ 需 DB |
