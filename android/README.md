> 开源副本：部署地址已替换为示例；历史验证记录不代表当前部署状态。

# 失业么 Android App

## 项目简介

失业么 —— AI 时代个体的"能力重建与收入重建"社区 Android 客户端。

## 技术栈

| 层级 | 方案 |
|------|------|
| 语言 | Kotlin |
| UI 框架 | Jetpack Compose |
| 架构 | MVVM + Clean Architecture |
| 网络层 | Retrofit + OkHttp |
| 本地存储 | Room + DataStore |
| 图片加载 | Coil |
| 依赖注入 | Hilt |
| 异步处理 | Kotlin Coroutines + Flow |

## 项目结构

```
app/src/main/java/com/shiyemo/app/
├── di/                    # Hilt 依赖注入
│   └── NetworkModule.kt
├── data/                  # 数据层
│   ├── remote/            # API 接口定义
│   │   ├── AuthApi.kt
│   │   └── PostApi.kt
│   ├── local/             # 本地存储
│   │   └── prefs/
│   │       └── TokenManager.kt
│   ├── repository/        # 仓库实现（待实现）
│   └── model/             # 数据模型
│       └── dto/
│           ├── ApiResponse.kt
│           ├── UserDto.kt
│           └── PostDto.kt
├── ui/                    # UI 层
│   ├── auth/              # 登录注册
│   │   ├── LoginScreen.kt
│   │   └── RegisterScreen.kt
│   ├── feed/              # 首页信息流
│   │   └── FeedScreen.kt
│   ├── post/              # 帖子
│   │   └── CreatePostScreen.kt
│   ├── profile/           # 个人中心
│   │   └── ProfileScreen.kt
│   ├── search/            # 搜索
│   │   └── SearchScreen.kt
│   ├── notification/      # 通知
│   │   └── NotificationScreen.kt
│   ├── navigation/        # 导航
│   │   └── AppNavigation.kt
│   └── theme/             # 主题
│       ├── Color.kt
│       ├── Theme.kt
│       └── Typography.kt
├── util/                  # 工具类
│   └── AuthInterceptor.kt
├── ShiyemoApp.kt          # Application
└── MainActivity.kt        # 主 Activity
```

## 开发环境要求

- Android Studio Hedgehog (2023.1.1) 或更高版本
- JDK 17
- Android SDK 34
- Gradle 8.2

## 如何运行

1. 使用 Android Studio 打开 `android/` 目录
2. 等待 Gradle 同步完成
3. 选择模拟器或真机设备
4. 点击 Run 按钮

## 后端 API

- Base URL: `SHIYEMO_API_BASE_URL`（用户级 `~/.gradle/gradle.properties`；默认 Android 模拟器宿主地址）
- 认证方式: JWT Bearer Token
- 响应格式: `{"code": 0, "message": "ok", "data": {}}`

## 开发计划

| Phase | 内容 | 周期 |
|-------|------|------|
| Phase 1 | 脚手架 + 认证 + 导航 | 1 周 |
| Phase 2 | 首页 Feed + 帖子详情 + 发帖 + 个人中心 | 2 周 |
| Phase 3 | 社交功能（点赞/收藏/评论/关注/搜索/通知） | 1 周 |
| Phase 4 | 打磨优化 | 1 周 |

## 注意事项

- 图片上传需要相机和存储权限
- 网络请求需要 INTERNET 权限
- Token 存储在 DataStore 中，退出登录时清除
