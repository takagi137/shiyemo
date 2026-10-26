# design：前端视觉重做

> 状态：**待评审**
> 前置：proposal.md 已评审通过（2026-07-21）
> 评审决策：主色=深鼠尾草绿 · 成果图标=幼苗 · 范围=主链路 5 页

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 令牌载体 | `@theme` CSS 变量 | Tailwind v4 原生机制，`bg-accent` 等工具类自动可用 |
| 主色 | 深鼠尾草绿 `oklch(48% 0.075 165)` | 已评审确定 |
| 颜色空间 | oklch | 感知均匀，便于让 5 类标签保持等明度 |
| 图标方案 | 自建 Icon 组件（lucide 路径） | 不引第三方运行时依赖，体积可控 |
| 成果图标 | 幼苗 | 已评审确定 |
| 组件类实现 | `@layer components` 原生 CSS | 已在 b18ee0a 确定，不用 @apply |
| 铺开范围 | 主链路 5 页 | 已评审确定 |

## 主链路 5 页定义

| 页面 | 路由 | 关联组件 |
|------|------|---------|
| Feed | `/feed` | Icon、constants |
| 帖子详情 | `/posts/[id]` | PostCard、Avatar、评论区 |
| 发帖 | `/posts/create` | RichTextEditor、POST_CREATE_TYPES |
| 登录/注册 | `/login`、`/register` | LoginForm、RegisterForm |
| 个人中心 | `/profile` | Avatar、FreedomDayCard |

> 其余 16 个路由不在本次范围，但因共享 globals.css 令牌，主色/画布会自动跟随，
> 不会出现「半绿半蓝」的割裂——只是它们内部的 `blue-*` 局部色暂不清理。

## 模块影响

| 模块 | 改动类型 | 说明 | 状态 |
|------|---------|------|------|
| globals.css | 修改 | 语义令牌层 | 样本已完成，需评审 |
| components/Icon.tsx | 新增 | 线性图标 | 样本已完成 |
| lib/constants.ts | 修改 | 类型配色→令牌 | 样本已完成 |
| app/feed/page.tsx | 修改 | 重做布局 | 样本已完成，**含已知 bug** |
| components/Navbar.tsx | 修改 | 配色跟随 | 样本已完成 |
| app/layout.tsx | 修改 | body 底色 | 样本已完成 |
| components/PostCard.tsx | 修改 | 旧色 14 处→令牌、emoji→Icon | 待做 |
| app/posts/[id]/page.tsx | 修改 | 旧色 53 处 | 待做 |
| app/posts/create/page.tsx | 修改 | 旧色 26 处 | 待做 |
| components/LoginForm.tsx | 修改 | 旧色 16 处 | 待做 |
| components/RegisterForm.tsx | 修改 | 旧色 36 处 | 待做 |
| app/profile/page.tsx | 修改 | 旧色 28 处 | 待做 |

## 必须先修的已知缺陷

feed 页筛选行 `[&::-webkit-scrollbar]{display:none}` 是无效语法（花括号形式不被
Tailwind v4 支持），导致溢出容器未生效、筛选按钮与排序控件重叠。

修法：改用标准的任意变体 `[&::-webkit-scrollbar]:hidden` 或直接用内联 style
`scrollbarWidth:'none'`；同时给排序控件容器加 `shrink-0`，给筛选区加正确的
`min-w-0 overflow-x-auto`。这是本次第一个要修的项。

## 令牌命名（已在样本确立）

```
canvas / surface / surface-sunk   画布与卡片表面
ink / ink-soft / ink-faint        文字三级
line / line-soft                  描边两级
accent / accent-hover / accent-soft / accent-line   主色四态
{emotion|insight|case|action|achievement}-{soft|line|ink}   五类标签
```

## 可访问性验证方法

每个语义令牌对搭配（ink on surface、accent on surface、各标签 ink on soft）
用 oklch 亮度差估算对比度，正文需 ≥4.5:1、大字号 ≥3:1。
样本中标签统一在 L≈50%(ink) 与 L≈96.5%(soft)，对比度约 7:1，满足 AA。

## 风险点

| 风险 | 应对 |
|------|------|
| 改动页面多，可能漏改局部旧色导致色彩不一致 | 每页改完 grep 复查 `blue-/gray-/rose-` 等残留 |
| PostCard 与 feed 页 PostRow 结构重复 | 本次不强行合并组件，仅统一视觉，避免扩大改动面 |
| 富文本编辑器 RichTextEditor 内联样式 | 仅改容器与按钮，不动 tiptap 内部渲染 |
| 既有前端测试断言了颜色类名 | 改前先看 __tests__，必要时同步更新断言 |

## 不引入的东西

- 不加深色模式（proposal 已列延后）
- 不改信息架构（导航结构、分页/无限滚动）
- 不引第三方 UI 库或图标库运行时依赖
