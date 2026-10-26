# spec：前端 DRY 重构

## 场景描述

消除前端重复代码，抽取公共模块。

### Given 前置条件
- 当前 7 个文件各自定义 POST_TYPES/TYPE_LABELS
- 4 个文件各自定义 timeAgo 函数
- 5 个页面各自实现分页 UI

### When 操作
- 创建 `lib/constants.ts` 统一定义 POST_TYPES 和 TYPE_LABELS
- 创建 `lib/utils.ts` 统一定义 timeAgo
- 创建 `components/Pagination.tsx` 复用分页 UI
- 各页面改为 import 公共模块

### Then 预期结果
- `npm run build` 通过，零错误
- `grep -r "POST_TYPES\|TYPE_LABELS" web/app/*/page.tsx` 只有 import 引用，无本地定义
- `grep -r "function timeAgo" web/` 只在 lib/utils.ts 中定义一次
- 各页面视觉效果与重构前一致
