# 部署

统一入口参见仓库根目录 [README](../../README.md#运行docker-compose)。

将 `deploy/.env.example` 复制为 `deploy/.env`，所有实际地址与凭据只写入本地配置。

- `APP_BIND_ADDRESS`：默认 `127.0.0.1`。
- `APP_HTTP_PORT`：默认 `8090`。
- `APP_PUBLIC_URL`：用户浏览器访问的完整源，用于邮件验证和找回密码。
- `IMAGE_REMOTE_ORIGINS`：构建时使用的远程图片源，逗号分隔。
- 小内存服务器可叠加 `docker-compose.server.yml`；该文件需要 Compose 2.24.4+。

首次安装必须先执行全部数据库迁移，不能只执行旧的单模块部署脚本。
不要提交 `.env`、数据库备份、上传目录、私钥或部署打包文件。
