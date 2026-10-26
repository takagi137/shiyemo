> 开源副本：部署地址已替换为示例；历史验证记录不代表当前部署状态。

# tasks：静思空间

## Phase 1：数据层

- [x] 任务 1.1：创建 journal 表 + readings 表
  - 验收：模型定义完成，GORM AutoMigrate 可创建
  - Spec：04-journal.md, 05-reading.md

- [x] 任务 1.2：实现 journal repo（CRUD）
  - 验收：单元测试通过（6 用例）
  - Spec：04-journal.md

- [x] 任务 1.3：实现 reading repo（CRUD + 随机 + 今日）
  - 验收：单元测试通过（4 用例）
  - Spec：05-reading.md

## Phase 2：业务层

- [x] 任务 2.1：实现 journal service
  - 验收：单元测试通过（6 用例）
  - Spec：04-journal.md

- [x] 任务 2.2：实现 reading service
  - 验收：单元测试通过（4 用例）
  - Spec：05-reading.md

## Phase 3：接口层

- [x] 任务 3.1：实现 journal handlers（POST/GET/DELETE）
  - 验收：go build 通过
  - Spec：04-journal.md

- [x] 任务 3.2：实现 reading handlers（today/random）
  - 验收：go build 通过
  - Spec：05-reading.md

- [x] 任务 3.3：注册路由到 main.go
  - 验收：服务器启动无报错
  - Spec：全部

## Phase 4：前端

- [x] 任务 4.1：创建 BreathingCircle 组件
  - 验收：呼吸动画正常运行（4-4-6节奏）
  - Spec：02-breathing.md

- [x] 任务 4.2：创建 AmbientPlayer 组件
  - 验收：音效选择/播放/暂停 UI 正常（音频文件待补充）
  - Spec：03-ambient-music.md

- [x] 任务 4.3：创建 JournalEntry 组件
  - 验收：日记保存/查看/删除正常
  - Spec：04-journal.md

- [x] 任务 4.4：创建 ReadingCard 组件
  - 验收：今日阅读展示正常
  - Spec：05-reading.md

- [x] 任务 4.5：创建静思空间主页面 /sanctuary
  - 验收：npm run build 通过，四个模块正常展示
  - Spec：01-enter-sanctuary.md

## Phase 5：验证

- [x] 任务 5.1：创建 BDD 场景文件
  - 验收：5 个 feature.md 文件创建完成
  - 文件：docs/05-场景验证/sanctuary/01-05

- [x] 任务 5.2：部署到验证环境
  - 验收：/sanctuary 页面正常访问（HTTP 200）
  - 说明：原 app.example.com 已不可达，改部署到共享生产机 app.example.com（:8090）

- [x] 任务 5.3：BDD 场景验证
  - 验收：BDD feature 已与实现对齐；后端 service 单测全绿（tips/reading/journal/dailyimage/ambientsound）
  - 备注：04-breathing 更新为 4-6 共振式；05-ambient-music 去进度条 + 补失败态

- [x] 任务 5.4：移动端适配验证
  - 验收：手机端已实测（首页/静思空间正常显示、可交互）

## Phase 6：上线后修复与规范化

- [x] 任务 6.1：修首页「换一句/换一张」失效（金句表未初始化 + 图片小时缓存）
  - Spec：迁移 022_seed_survival_tips.sql；dailyimage RandomImageFresh
- [x] 任务 6.2：呼吸练习改共振式呼吸 4-6（去满肺屏息，消除窒息感）
- [x] 任务 6.3：环境音补 6 个 CC0 音频，随后重构为后端维护列表 + 外链
  - Spec：ambientsound 模块 + 迁移 023_create_ambient_sounds.sql，对齐 daily_images
- [x] 任务 6.4：补测试与 BDD 对齐（ambientsound / dailyimage service 单测，覆盖 ≥80%）
