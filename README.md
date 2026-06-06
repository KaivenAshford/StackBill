# StackBill

面向 PC 浏览器的数字资产与订阅管理平台。

## 功能

- **订阅管理** — 管理软件订阅、AI 工具订阅、周期性支出，支持按分类/状态筛选
- **数字资产管理** — 管理域名、服务器、Docker 服务、SSL 证书、API Key 等数字资产
- **分类管理** — 自定义分类，系统为新用户初始化 8 个默认分类
- **仪表盘** — 本月/年度支出统计、分类支出占比图表、即将续费与到期列表
- **到期提醒** — 订阅续费提醒、资产到期提醒、异常资产警告
- **多用户数据隔离** — 所有数据严格按用户隔离
- **PC 端界面** — 侧边栏布局，优先优化桌面浏览器使用体验
- **国际化** — 支持 zh-CN / en-US
- **暗黑模式** — 支持浅色/深色主题切换

## 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go / Gin / GORM / PostgreSQL / JWT |
| 前端 | Vue 3 / Vite / TypeScript / Pinia / Naive UI / ECharts |
| 部署 | Docker Compose |

## 项目结构

```
StackBill/
├── backend/                 # Go 后端
│   ├── cmd/server/          # 入口
│   ├── internal/
│   │   ├── api/             # 控制器
│   │   ├── config/          # 配置
│   │   ├── dto/             # 数据传输对象
│   │   ├── middleware/       # 中间件（JWT、CORS）
│   │   ├── model/           # 数据模型
│   │   ├── repository/      # 数据访问层
│   │   ├── router/          # 路由
│   │   └── service/         # 业务逻辑
│   ├── docs/                # Swagger 生成的 API 文档
│   ├── pkg/
│   │   ├── database/        # 数据库连接与迁移
│   │   └── response/        # 统一响应
├── frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── api/             # API 请求
│   │   ├── composables/     # 组合式函数
│   │   ├── i18n/            # 国际化配置
│   │   ├── layouts/         # 布局组件
│   │   ├── locales/         # 语言文件（zh-CN / en-US）
│   │   ├── router/          # 路由
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── types/           # TypeScript 类型
│   │   ├── utils/           # 工具函数
│   │   └── views/           # 页面视图
│   └── nginx.conf           # Nginx 配置
├── docker-compose.yml
├── .env.example
├── config.example.yaml
└── CLAUDE.md
```

## 快速开始

### Docker Compose 部署（推荐）

```bash
cp .env.example .env
# 编辑 .env 设置 JWT_SECRET 和数据库密码
docker compose up -d
```

访问 http://localhost 即可使用。

服务包含三个容器：

- **frontend** — Nginx 托管前端静态文件，反向代理 API 到后端
- **backend** — Go API 服务
- **postgres** — PostgreSQL 17 数据库

### 本地开发

**前提条件：** Go 1.24+、Node.js 20+、PostgreSQL 17

**后端：**

```bash
cd backend
cp ../config.example.yaml config.yaml
# 编辑 config.yaml 设置数据库连接
go run ./cmd/server
```

后端默认启动在 `http://localhost:8080`，首次启动会通过 GORM AutoMigrate 自动创建数据库表。

**前端：**

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器默认启动在 `http://localhost:3000`，已配置 API 代理到后端。

## API 概览

所有接口前缀：`/api/v1`

启动服务后访问 **http://localhost/swagger/index.html** 查看完整的 Swagger API 文档。

### 统一返回格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 分页格式

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "page_size": 20
}
```

### 接口列表

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /auth/register | 用户注册 |
| POST | /auth/login | 用户登录 |
| GET | /auth/me | 获取当前用户 |
| PUT | /users/profile | 修改资料 |
| PUT | /users/password | 修改密码 |
| GET | /dashboard | 仪表盘统计 |
| GET/POST/PUT/DELETE | /categories | 分类 CRUD |
| GET/POST/PUT/DELETE | /subscriptions | 订阅 CRUD |
| GET/POST/PUT/DELETE | /assets | 资产 CRUD |
| GET | /health | 健康检查 |
| GET | /reminders | 提醒列表 |
| PUT | /reminders/:id/read | 标记已读 |
| PUT | /reminders/read-all | 全部已读 |
| DELETE | /reminders/:id | 忽略提醒 |

所有写操作需要 JWT 鉴权（`Authorization: Bearer <token>`）。

## 配置

### 环境变量（.env）

| 变量 | 说明 | 默认值 |
|---|---|---|
| JWT_SECRET | JWT 签名密钥 | — |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 5432 |
| DB_USER | 数据库用户 | stackbill |
| DB_PASSWORD | 数据库密码 | — |
| DB_NAME | 数据库名称 | stackbill |
| SERVER_PORT | 后端端口 | 8080 |

### 配置文件（config.yaml）

参考 `config.example.yaml`，支持服务器、数据库、JWT、日志等配置。环境变量会覆盖配置文件中的对应项。

## 数据库迁移

项目提供两种迁移方式：

**自动迁移（开发/默认）：** 后端启动时通过 GORM AutoMigrate 自动创建和更新数据库表结构，无需手动操作。

**手动迁移（生产推荐）：** 使用 `backend/migrations/` 目录下的 SQL 文件，配合 [golang-migrate](https://github.com/golang-migrate/migrate) 等工具执行版本化迁移：

```bash
# 示例：使用 golang-migrate CLI
migrate -path backend/migrations -database "postgres://stackbill:stackbill_password@localhost:5432/stackbill?sslmode=disable" up
```

迁移文件命名规则：`{序号}_{描述}.up.sql`（执行） / `{序号}_{描述}.down.sql`（回滚）。

## License

MIT
