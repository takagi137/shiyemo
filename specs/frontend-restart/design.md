# design：前端重新初始化

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 前端框架 | Next.js 14 | 项目技术栈要求，SSR/SEO友好 |
| 语言 | TypeScript 5.x | 类型安全，AI编码效率高 |
| 样式方案 | TailwindCSS 3.x | 实用优先，快速开发 |
| 包管理 | pnpm | 快速，磁盘空间效率高 |
| 项目结构 | App Router | Next.js 14 推荐方式 |

## 模块影响

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| web/ | 新建 | 重新初始化前端项目 |
| web/app/login/ | 新建 | 登录页面 |
| web/lib/ | 新建 | API客户端、工具函数 |
| web/components/ | 新建 | 基础组件 |

## 接口设计

### 登录页面路由

```
GET /login
```

### 登录API调用

```
POST /api/v1/auth/login/email
请求：{ email: string, password: string }
响应：{ code: number, message: string, data: { token: string, user: object } }
```

## 数据变更

无数据库变更，仅前端项目重新初始化。

## 风险点

1. 现有前端代码可能需要备份或重命名
2. 需要确保与现有后端API兼容
3. 登录页面设计需要与现有风格保持一致

## 实现阶段

### Phase 1：项目初始化
- 创建 Next.js 14 项目
- 配置 TypeScript + TailwindCSS
- 设置基础目录结构

### Phase 2：登录页面
- 实现登录表单
- 对接后端登录API
- 错误处理和加载状态

### Phase 3：基础组件
- 创建通用按钮、输入框组件
- 实现响应式布局
- 配置全局样式