# tasks：前端视觉重做

> 状态：**主链路完成，已部署 8090 验证**
> 前置：proposal.md、design.md、specs/*.feature.md 已评审
> 评审决策：主色=深鼠尾草绿 · 成果图标=幼苗 · 范围=主链路 5 页

## Phase 0：验证脚手架（测试优先）

- [x] 0.1 编写视觉契约检查脚本 `web/scripts/check-design.mjs`
  - 完成：脚本对未改完页面正确报红 11 项，证明检查有效

## Phase 1：修复已知缺陷

- [x] 1.1 修复 feed 筛选行重叠（spec 03）
  - 完成：去掉无效的 `[&::-webkit-scrollbar]{...}`，改用 .scroll-x-hidden

## Phase 2：关联组件

- [x] 2.1 PostCard 组件重做（复用 constants，删除重复 TYPE_CONFIG）
- [x] 2.2 LoginForm + RegisterForm 重做（裸 input → .input-field）

## Phase 3：主链路页面

- [x] 3.1 posts/[id] 帖子详情页
- [x] 3.2 posts/create 发帖页
- [x] 3.3 login + register 页面外壳
- [x] 3.4 profile 个人中心

## Phase 4：全量验证

- [x] 4.1 check-design 全绿（spec 01 令牌 + spec 02 无残留）
- [x] 4.2 vitest run：105 用例全通过（spec 04）
- [x] 4.3 构建 + 部署到 8090：主链路 6 个端点 + API 全部 HTTP 200
- [x] 4.4 主链路截图人工确认：feed / login / register / 详情 / 发帖

## Phase 5：收尾

- [x] 5.1 更新 tasks 勾选状态
- [ ] 5.2 合并回功能分支，等待推送（代理恢复后）

## 过程中修复的额外缺陷（超出原计划，均已记录提交）

- deploy/.env 被 git 追踪，同步时覆盖服务器真实密钥软链接
  → 导致 api 认证失败崩溃、nginx 502。已停止追踪并加 .gitignore。
- constants.icon 从 emoji 改为图标名后，若沿用 {typeInfo.icon} 直接渲染
  会把图标名当文本显示（详情页「emotion 情绪共鸣」）。
  → 主链路 + 连带回归（搜索/编辑页）改用 <Icon>。
- web 构建 npm install 因 registry.npmjs.org 不稳定卡死
  → Dockerfile 新增可覆盖 NPM_REGISTRY，服务器走 npmmirror。

## 未做（proposal 已列延后，需另开提案）

- 深色模式、移动端底部导航、无限滚动、全局 auth context
- 非主链路 16 页的局部旧色清理（admin/checkin/search/edit 等仍含 blue-/gray-，
  但主色与画布已随令牌自动跟随，无割裂）
