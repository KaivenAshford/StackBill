# StackBill

多端同步数字资产与订阅管理平台。

## 功能

- 软件订阅管理（AI 工具、开发工具、云服务等）
- 数字资产管理（域名、服务器、SSL 证书、API Key 等）
- 分类管理
- 到期提醒
- 仪表盘统计
- 多用户数据隔离
- PC 端 + 移动端响应式界面
- 国际化（zh-CN / en-US）

## 技术栈

**后端:** Go / Gin / GORM / PostgreSQL / JWT

**前端:** Vue 3 / Vite / TypeScript / Pinia / Naive UI / ECharts

**部署:** Docker Compose

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
│   │   ├── service/         # 业务逻辑
│   │   └── task/            # 定时任务
│   ├── pkg/
│   │   ├── database/        # 数据库连接
│   │   └── response/        # 统一响应
│   ├── migrations/          # 数据库迁移
│   └── docs/                # API 文档
├── frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── api/             # API 请求
│   │   ├── components/      # 组件
│   │   ├── i18n/            # 国际化
│   │   ├── layouts/         # 布局
│   │   ├── locales/         # 语言文件
│   │   ├── router/          # 路由
│   │   ├── stores/          # Pinia 状态
│   │   ├── types/           # TypeScript 类型
│   │   ├── utils/           # 工具函数
│   │   └── views/           # 页面
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
# 编辑 .env 设置 JWT_SECRET
docker compose up -d
```

访问 http://localhost

### 本地开发

**后端:**

```bash
cd backend
cp config.example.yaml config.yaml
# 编辑 config.yaml 设置数据库连接
# 需要本地运行 PostgreSQL
go run ./cmd/server
```

**前端:**

```bash
cd frontend
npm install
npm run dev
```

访问 http://localhost:3000

## API

所有接口前缀: `/api/v1`

统一返回格式:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 后续开发

1. 完成用户注册登录 API
2. 完成订阅 CRUD
3. 完成资产 CRUD
4. 完成分类管理
5. 完成提醒模块
6. 完成仪表盘统计
7. 移动端适配优化
