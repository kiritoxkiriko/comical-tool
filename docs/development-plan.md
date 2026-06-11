# comical-tool 开发计划

> 本文是 `comical-tool` 首个完整版本的开发计划来源。项目边界和 Agent 行为约束见 `AGENTS.md`。

## 1. 目标

将 `comical-tool` 做成一个小而完整的实用工具平台，首版包含：

- 短链接
- 图床
- 临时剪贴板
- 文件暂存
- Web 界面
- CLI
- VitePress 文档站
- Docker 自托管部署
- Cloudflare 部署支持

仓库顶层按产品和 runtime 领域区分，每个领域内部再遵循对应语言的标准结构。Go 代码不直接把 `cmd/`、`internal/`、`pkg/` 铺在仓库根目录，而是在 `server/`、`cli/` 内部分别组织。根目录可以用 `go.work` 串联 Go module。

```text
comical-tool/
  server/
    cmd/comical-tool/
    internal/
    pkg/
  cli/
    cmd/comical-cli/
    internal/
  web/
  worker/
  docs/
  deploy/
  migrations/
  scripts/
  go.work
```

## 2. 架构决策

### 2.1 Runtime 拆分

项目支持两条运行路径。

1. 自托管 runtime
   - `server/cmd/comical-tool` 运行 Go Hertz API 服务。
   - `web` 运行 Next.js Web 界面。
   - `cli/cmd/comical-cli` 提供 CLI，通过 HTTP API 调用服务。
   - 默认数据库使用 SQLite。
   - 默认缓存使用 Redis。
   - 默认对象存储使用本地文件系统。

2. Cloudflare runtime
   - `web` 通过 OpenNext 构建到 Cloudflare。
   - `worker` 提供 Cloudflare Worker API adapter。
   - D1 存储元数据。
   - R2 存储图片和文件。
   - KV 存储缓存类和短期易失数据。
   - Cron Trigger 执行清理任务。

不要尝试在 Cloudflare Worker 中运行 Hertz。Worker 应作为 adapter，复用 `server/pkg` 中的纯 Go 包；需要时可以复用编译为 WASM 的纯逻辑。

### 2.2 共享 Go 包

可复用逻辑使用标准 Go 包结构，放在 `server/pkg/` 下，不新增顶层 `core/` 或根级 `pkg/` 目录。

初始包划分：

- `server/pkg/domain`：资源模型、资源类型常量、状态辅助方法
- `server/pkg/policy`：slug、TTL、过期、口令、访问次数策略
- `server/pkg/apperror`：server、CLI、Worker adapter 共享的稳定错误码

这些包只承载不依赖具体 runtime 的逻辑：

- slug 生成和校验
- TTL 计算
- 过期判断
- 口令 hash 和校验封装
- 最大访问次数策略
- 资源状态流转
- 统一错误码

`server/pkg/` 包不得依赖 Hertz、sqlx、Redis、S3、R2、D1、KV 等框架或基础设施客户端。

### 2.3 数据归属

所有资源都归属于用户。v1 只内置 `guest` 用户，但 schema 必须支持后续接入登录系统。

资源状态使用软状态字段：

- `expires_at`
- `deleted_at`
- `revoked_at`
- `created_at`
- `updated_at`

这样可以让清理、审计和后续恢复策略保持清晰。

## 3. 仓库初始化阶段

### 3.1 创建领域化目录结构

创建以下目录：

- `server/`
- `server/cmd/comical-tool/`
- `server/internal/`
- `server/pkg/`
- `cli/`
- `cli/cmd/comical-cli/`
- `cli/internal/`
- `web/`
- `worker/`
- `docs/`
- `deploy/`
- `migrations/`
- `scripts/`

`server/`、`cli/`、`web/`、`worker/`、`docs/`、`deploy/`、`migrations/`、`scripts/` 需要简短 README，说明目录职责、本地命令和边界。Go 子包以 package doc 为主，不需要在每个 package 下放 README。

### 3.2 根目录文件

创建或更新：

- `README.md`
- `AGENTS.md`
- `.gitignore`
- `.editorconfig`
- `Makefile` 或 `justfile`
- `go.work`
- `deploy/docker-compose.yml`
- `deploy/config.example.toml`

根 README 保持简短，只放：

- 项目简介
- 快速启动
- 目录结构
- 文档入口

详细说明放到 `docs/`。

## 4. Server 阶段

### 4.1 Go server 骨架

在 `server/` 下创建 Go module。`server/cmd/comical-tool` 是 server 入口，最终二进制名为 `comical-tool`，使用：

- Hertz 作为 HTTP server
- Viper 读取 TOML 配置
- sqlx 访问数据库
- golang-migrate 或兼容迁移工具执行 migration

Go 包结构：

```text
server/
  go.mod
  cmd/
    comical-tool/
      main.go
  internal/
    config/
    http/
    middleware/
    module/
      shortlink/
      image/
      clipboard/
      filestash/
    repository/
    storage/
    cache/
    job/
  pkg/
    domain/
    policy/
    apperror/
```

共享纯 Go 包放在 `server/pkg/`：

```text
server/pkg/
  domain/
  policy/
  apperror/
```

### 4.2 配置

支持 TOML 配置，并允许环境变量覆盖。

必需配置段：

```toml
[server]
addr = "127.0.0.1:8080"
public_base_url = "http://localhost:8080"
max_body_bytes = 104857600

[database]
driver = "sqlite"
dsn = "file:comical.db?_foreign_keys=on"

[cache]
driver = "redis"
dsn = "redis://localhost:6379/0"

[storage]
driver = "local"
local_dir = "./data/objects"

[security]
admin_token = "change-me"
content_encryption_key = "change-me-32-bytes"

[modules.short_link]
default_ttl = "168h"
allow_custom_slug = true
domain_mappings = { "s.tool.sqlboy.me" = "https://tool.sqlboy.me/short" }

[modules.image_hosting]
default_ttl = "720h"
max_bytes = 10485760

[modules.clipboard]
default_ttl = "1h"
max_visits = 5

[modules.file_stash]
default_ttl = "168h"
max_bytes = 104857600
```

### 4.3 数据库 schema

先创建 SQLite migration，再补 PostgreSQL 和 MySQL migration。

数据表：

- `users`
- `short_links`
- `assets`
- `clipboard_items`
- `resource_links`
- `access_events`

初始化数据：

- `guest` 用户

索引：

- `short_links.slug`
- `short_links.expires_at`
- `assets.owner_id`
- `assets.expires_at`
- `clipboard_items.expires_at`
- `access_events.resource_type, resource_id`

### 4.4 HTTP 约定

API 路由统一使用 `/api` 前缀。

成功响应格式：

```json
{
  "data": {}
}
```

错误响应格式：

```json
{
  "error": {
    "code": "resource_expired",
    "message": "resource expired",
    "request_id": "..."
  }
}
```

中间件：

- request ID
- 日志
- panic recovery
- CORS
- guest 用户上下文
- 管理 API 的 admin token 校验
- 请求体大小限制

## 5. 功能阶段

### 5.1 短链接

能力：

- 创建随机短链接
- 创建自定义短链接
- 校验 slug 格式
- 拒绝重复的 active slug
- 设置过期时间
- revoke 短链接
- 按 slug 跳转
- 支持独立短域名，并把短域名映射到主站 path 下，例如 `s.tool.sqlboy.me/{slug}` 等价于 `tool.sqlboy.me/short/{slug}`
- 记录访问事件
- 绑定图片、剪贴板或文件资源

路由：

- `POST /api/short-links`
- `GET /short/{slug}`
- `GET /{slug}`，用于独立短域名或反向代理后的短路径
- `POST /api/short-links/{slug}/revoke`

### 5.2 图床

能力：

- 上传图片
- 校验 MIME 和大小
- 通过 storage adapter 存储对象
- 创建元数据记录
- 查看图片列表
- 删除图片
- 设置过期时间
- 创建关联短链接

路由：

- `POST /api/images`
- `GET /api/images`
- `GET /api/assets/{id}`
- `DELETE /api/images/{id}`

### 5.3 临时剪贴板

能力：

- 创建文本剪贴板条目
- 使用较短默认 TTL
- 可选口令
- 可选最大访问次数
- 读取剪贴板条目
- 过期后拒绝访问
- 达到访问次数上限后拒绝访问
- 创建关联短链接
- 删除条目

路由：

- `POST /api/clip`
- `GET /api/clip/{id}`
- `DELETE /api/clip/{id}`

### 5.4 文件暂存

能力：

- 上传文件
- 默认 7d TTL
- 可选口令
- 可选最大访问次数
- 下载文件
- 查看文件列表
- 删除文件
- 创建关联短链接

路由：

- `POST /api/files`
- `GET /api/files`
- `GET /api/assets/{id}`
- `DELETE /api/files/{id}`

## 6. 存储、缓存和清理阶段

### 6.1 Storage adapter

实现统一对象存储接口：

- 本地文件系统
- S3 兼容存储
- Worker runtime 下的 Cloudflare R2

操作：

- `Put`
- `Get`
- `Delete`
- `Head`

不要在 storage 包以上暴露具体 provider client。

### 6.2 Cache adapter

实现统一缓存和易失状态接口：

- Redis
- in-memory
- Worker runtime 下的 Cloudflare KV

缓存用途：

- 短期查询加速
- 必要时记录一次性访问状态
- 后续限流

数据库仍然是 source of truth。

### 6.3 清理任务

自托管：

- server scheduler 定期扫描过期资源
- 第一阶段先把行标记为 deleted/revoked
- 对象删除可以作为第二阶段执行

Cloudflare：

- Worker Cron Trigger 扫描 D1
- 删除过期 R2 对象
- 清理关联 KV 条目

## 7. Web 阶段

`web` 使用 Next.js、TypeScript 和 TailwindCSS。

UI 要求：

- favicon/logo 使用 `https://i.loli.net/2021/02/11/JLHnIjOvFl7PC4o.png`
- 配色使用图片中的黄/红
- 顶部导航 tab 区分模块
- v1 不做营销 landing page
- 首屏就是工具界面

页面：

- `/` 跳转到短链接工具或展示默认 tab
- `/short-links`
- `/images`
- `/clipboard`
- `/files`

共享 UI：

- 顶部导航
- 模块 tab 状态
- 过期时间选择器
- 口令输入框
- 最大访问次数输入框
- 带复制按钮的结果卡片
- 错误 toast
- loading 状态
- 资源列表表格

## 8. CLI 阶段

CLI 使用 Go + Cobra 实现。`cli/cmd/comical-cli/main.go` 只负责创建 root command、加载配置并执行命令树，最终二进制名为 `comical-cli`，具体命令放在 `cli/internal/command/` 下。

命令：

```text
comical-cli config init
comical-cli short create
comical-cli short revoke
comical-cli image upload
comical-cli image list
comical-cli image delete
comical-cli clip put
comical-cli clip get
comical-cli clip delete
comical-cli file upload
comical-cli file download
comical-cli file delete
comical-cli admin cleanup
```

CLI 读取：

- config 文件路径
- API base URL
- admin token
- 默认 TTL
- 输出格式：table 或 json

命令组织：

- root command：全局 flags、配置加载、输出格式、API client 初始化
- `short` 子命令组：短链接创建、查询、revoke
- `image` 子命令组：图片上传、列表、删除
- `clip` 子命令组：临时剪贴板创建、读取、删除
- `file` 子命令组：文件上传、下载、删除
- `admin` 子命令组：清理任务等管理命令

CLI 在 `cli/` 下单独组织 Go module。

```text
cli/
  go.mod
  cmd/
    comical-cli/
      main.go
  internal/
    config/
    client/
    command/
      root.go
      short.go
      image.go
      clip.go
      file.go
      admin.go
```

CLI 入口为 `cli/cmd/comical-cli`，命令框架使用 Cobra。CLI 应调用 HTTP API，不直接复用 `server/internal` 的 handler/service。CLI 如需共享稳定错误码或策略，可通过 `go.work` 依赖 `server/pkg/apperror` 等纯 Go 包，但不能依赖 `server/internal`。

## 9. Cloudflare 阶段

### 9.1 Worker adapter

`worker` 实现 Cloudflare 特化 API routing。

Bindings：

- D1 存关系型元数据
- R2 存图片和文件对象
- KV 存缓存和易失数据
- Cron Trigger 做清理

Worker 应尽量和 server API 保持一致。差异必须记录在 `docs/cloudflare.md`。

### 9.2 从 Go 包编译 WASM

选择部分 `server/pkg/` 函数编译为 WASM。

首批范围：

- slug 校验
- TTL 计算
- 资源状态判断
- 访问次数上限判断

不要把数据库、存储或 HTTP framework 逻辑放进 WASM。

### 9.3 Web 上 Cloudflare

使用 OpenNext Cloudflare 构建 `web`。

部署配置包含：

- Worker entry
- static assets binding
- R2 bucket binding
- D1 binding
- KV namespace binding
- OpenNext 所需 compatibility flags

## 10. 文档阶段

`docs` 使用 VitePress。

必需文档：

- 项目概览
- 快速开始
- 目录结构
- 配置说明
- API 说明
- CLI 说明
- 本地开发
- Docker 部署
- Cloudflare 部署
- 数据库迁移
- 存储后端配置
- 故障排查

文档应偏可执行，优先提供命令和预期结果。

## 11. 部署阶段

### 11.1 Docker

提供：

- server Dockerfile
- web Dockerfile
- 本地依赖 docker-compose
- SQLite 和本地对象存储 volume
- 可选 PostgreSQL/MySQL profile
- Redis service

本地开发默认：

- SQLite
- Redis
- 本地对象存储

### 11.2 生产自托管

文档需要说明：

- SQLite 单机模式
- PostgreSQL/MySQL 生产模式
- S3 存储模式
- 反向代理示例
- DB 和对象存储备份建议

## 12. 测试计划

Server：

- `server/pkg/` 逻辑单测
- SQLite repository 测试
- PostgreSQL/MySQL migration smoke
- 所有模块 API 集成测试
- storage adapter contract test

Web：

- 共享控件组件测试
- 各模块页面表单和列表测试
- build 验证

CLI：

- Cobra command 解析和 flags 测试
- 配置读取测试
- API client 错误处理测试

Worker：

- routing 测试
- D1/R2/KV binding smoke test
- WASM import 测试
- Cron 清理测试

部署：

- `docker compose config`
- server container build
- web container build
- OpenNext Cloudflare build
- 可行时执行 Wrangler dry-run

## 13. 实施顺序

1. 仓库骨架和根文档
2. `server/` Go module、`server/cmd`、`server/internal`、`server/pkg` 骨架
3. `cli/` Go module、`cli/cmd`、`cli/internal` 骨架
4. `server/pkg` models、errors、TTL、slug、password、visit policy
5. server 配置和 Hertz 骨架
6. SQLite migration 和 repository 层
7. 短链接模块
8. 图床模块
9. 剪贴板模块
10. 文件暂存模块
11. storage adapters
12. cache adapters 和 cleanup jobs
13. Next.js web app
14. CLI
15. VitePress docs
16. Docker 部署
17. Cloudflare Worker adapter
18. Cloudflare web 部署
19. 最终集成测试和 README 完善

每个阶段完成后，仓库都应保持可运行、可验证。

## 14. 验收标准

首个完整版本满足以下条件：

- 四个模块都可以从 Web UI 使用
- 四个模块都可以从 CLI 使用
- `server/cmd/comical-tool` 可以使用 SQLite 在本地运行
- Redis-backed cache 可以在本地 Docker 中工作
- 本地对象存储可用
- SQLite、PostgreSQL、MySQL 都有 migration
- Docker 部署可以启动系统
- docs site 可以构建
- Cloudflare Worker build 存在，并包含 D1/R2/KV bindings
- Cloudflare 部署路径有文档
- 项目 README 链接到 docs，并说明快速启动方式
