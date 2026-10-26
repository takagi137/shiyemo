# 依赖安装阶段（构建需要 devDependencies）
FROM node:20-alpine AS deps
WORKDIR /app

# npm registry 做成可覆盖参数：默认官方源，国内机器可传 https://registry.npmmirror.com
ARG NPM_REGISTRY=https://registry.npmjs.org
ENV NPM_CONFIG_REGISTRY=$NPM_REGISTRY

COPY package.json package-lock.json ./
RUN npm install

# 构建阶段
FROM node:20-alpine AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# 构建时传入环境变量
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ARG IMAGE_REMOTE_ORIGINS
ENV IMAGE_REMOTE_ORIGINS=$IMAGE_REMOTE_ORIGINS

# 限制 Node 堆大小，避免在小内存机器上构建时 OOM
ARG NODE_OPTIONS
ENV NODE_OPTIONS=$NODE_OPTIONS

RUN npm run build

# 运行阶段
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV production
ENV HOSTNAME=0.0.0.0

# 创建非 root 用户
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# 复制构建产物
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000

CMD ["node", "server.js"]
