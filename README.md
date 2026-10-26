# 失业么 · Shiyemo

面向失业、待转型与重新出发的人，围绕情绪共鸣、能力重建和行动记录搭建的社区项目。

包括帖子、评论、点赞收藏、关注、搜索、管理后台，以及静思空间、日记和共读等模块；同时保留 Android 客户端源码。

这是原项目主分支的脱敏开源快照，使用全新的 Git 历史。私人交接记录、部署归档、运行数据及实际环境配置不随仓库发布。`specs/` 中的进度与验证结果是历史开发记录，不代表当前部署验收。

## 可以在这里做什么

有时候需要的不是一句“赶紧振作”，而是一个能说说近况、记下一点进展，也能安静待一会儿的地方。

| 想做的事 | 对应功能 |
| --- | --- |
| 看看其他人怎么度过这段时间 | 社区动态：五类帖子、最新 / 最热、评论、点赞与收藏 |
| 把今天的感受记下来 | 每日打卡：五种心情、短文字、最多三张图片，可选公开或私密 |
| 分享经验，或者发起一个小行动 | 发帖：富文本、标签、匿名选项，区分情绪、认知、案例、行动和成果 |
| 找到处境相近的人 | 同温层：按城市、行业、离职时间等资料查看伙伴与统计 |
| 暂时放下信息流 | 静思空间：呼吸练习、环境音、今日阅读与私密日记 |
| 找回有用的内容，持续交流 | 搜索、收藏、关注、通知与个人资料 |
| 维护自己的社区 | 管理后台：内容审核、举报、评论、反馈、图片与站点设置 |

## 界面预览

以下是开源 Web 前端的实际页面截图，接口内容使用虚构示例数据，人数和互动数不代表真实运营情况。截图用于介绍界面，不代表后端、音频播放或生产部署验收。

### 社区动态 · 先看看大家最近过得怎么样

心情打卡和长帖分区展示；顶部可以按内容类型筛选、切换最新与最热。点击帖子继续阅读与交流。

![社区动态：心情打卡与不同类型的示例帖子](docs/images/feed.png)

### 静思空间 · 给自己一点安静的时间

呼吸练习、环境音、短篇阅读放在同一个页面；登录后可以书写和回看私密日记。

![静思空间：呼吸练习、环境音、阅读和日记入口](docs/images/sanctuary.png)

### 每日打卡 · 不需要写很多，也能留下一点记录

选择心情，写几句话，决定是否公开到社区。下面展示的是填写中的示例，尚未提交保存。

![每日打卡：心情选择、文字、图片和公开开关](docs/images/checkin.png)

继续查看 [完整功能与界面导览](docs/FEATURES.md)，了解同温层、手机端页面、发帖流程和各模块入口。

## 技术栈

- Web：Next.js 16.2、React 19.2、TypeScript、Tailwind CSS。
- API：Go 1.24、Gin、GORM。
- 服务：PostgreSQL 15、Redis 7、Meilisearch；本地文件存储或 MinIO。
- Android：Kotlin、Jetpack Compose、Retrofit、Hilt。

## 运行（Docker Compose）

需要 Docker Compose v2、Git，以及 Bash（Windows 可使用 Git Bash / WSL）。以下命令从仓库根目录运行。

```bash
git clone https://github.com/takagi137/shiyemo.git
cd shiyemo
cp deploy/.env.example deploy/.env
```

编辑 `deploy/.env`，为 `DB_PASSWORD`、`JWT_SECRET`、`MEILI_KEY`、`MINIO_SECRET_KEY` 分别填入独立随机值，例如逐次运行 `openssl rand -hex 32`。这些变量未填时 Compose 会拒绝启动。实际 `.env` 已被 Git 忽略。

```bash
# 先启动数据库，确认 ready 后再应用迁移
docker compose --project-directory deploy -f deploy/docker-compose.yml up -d postgres
docker compose --project-directory deploy -f deploy/docker-compose.yml exec postgres pg_isready -U shiyemo -d shiyemo
bash scripts/migrate.sh
# 迁移成功后启动完整应用
docker compose --project-directory deploy -f deploy/docker-compose.yml up -d --build
```

默认访问 `http://localhost:8090`，API 为 `http://localhost:8090/api/v1`。默认只绑定本机；外网部署请配置自己的 HTTPS 反向代理和 `APP_PUBLIC_URL`。这是部署说明，本次开源整理未对全新 Docker 环境做上线验收。

## 配置入口

| 场景 | 配置方式 |
| --- | --- |
| Compose 服务密码、访问端口、邮件链接 | `deploy/.env`，参考 `.env.example` |
| 前端独立开发 | 复制 `web/.env.example` 为 `web/.env.local` |
| 邮箱验证 / 密码重置链接 | 后端环境变量 `APP_PUBLIC_URL` |
| 远程图片来源 | `IMAGE_REMOTE_ORIGINS`，逗号分隔的精确 HTTP(S) 源，修改后重新构建 Web |
| Android API | 用户目录 `~/.gradle/gradle.properties` 设置 `SHIYEMO_API_BASE_URL=https://app.example.com/api/v1/`，末尾保留 `/` |
| SMTP | 环境变量或管理后台设置，凭据只留在私有环境 |

Android 默认地址 `10.0.2.2` 是 Android 模拟器访问宿主机的标准地址；真机需要设置自己的服务地址。Android 上传解析仍有原项目 TODO，不代表已完成真机验收。

独立运行后端时需要自行提供数据库、缓存、搜索等连接环境变量；Go 程序不会自动读取 `deploy/.env`。代码中兼容旧测试的固定默认值只适用于本地开发，不能当作生产凭据。Compose 配置已强制要求外部注入主要密码和密钥。

## 目录与验证

- `server/`：API、测试与数据库迁移；已有迁移文件保持原样。
- `web/`：页面、组件和前端测试。
- `android/`：安卓客户端源码。
- `deploy/`：容器与环境变量模板。
- `docs/`、`specs/`：产品、架构及历史需求资料。
- `scripts/`：测试与迁移工具。旧的 `deploy-sanctuary.sh` 是模块维护脚本，首次安装使用上面的全量迁移流程。

常用检查：`cd server && go test ./internal/verify ./internal/common/config`；`cd web && npm ci && npm run lint && npm run build`。集成测试需要独立测试数据库，不能指向生产数据。

## 许可证

项目代码采用 [MIT](LICENSE)。依赖库沿用各自许可证。种子数据引用了 Unsplash、BigSoundBank、OpenGameArt 等外部图片与音频链接，MIT 不替代这些素材各自的授权条款；使用时请核对原来源。

## CI 模板

自动检查配置保留在 [docs/ci.yml.example](docs/ci.yml.example)。当前发布使用的 GitHub CLI 授权没有 workflow scope，因此未启用 Actions；后续获得对应权限后可复制为 .github/workflows/ci.yml。

