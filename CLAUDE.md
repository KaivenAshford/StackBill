# CLAUDE.md

## 项目名称

StackBill

## 项目定位

StackBill 是一个面向技术用户和自托管用户的 PC 浏览器数字资产与订阅管理平台。

它不只是普通的订阅记账工具，而是用于统一管理：

* 软件订阅
* AI 工具订阅
* 域名
* 服务器
* Docker 服务
* SSL 证书
* API Key
* 云服务账单
* 周期性支出
* 服务到期提醒
* 数字资产状态

当前阶段优先支持 PC 端浏览器访问。移动端体验暂不纳入当前计划，等用户明确提及时再单独规划。

---

## 核心目标

第一版目标是完成一个可自托管部署的基础可用版本。

必须实现：

* 用户注册登录
* 多用户数据隔离
* 订阅管理
* 数字资产管理
* PC 端页面
* 首页数据统计
* 到期提醒
* Docker Compose 部署
* 国际化基础结构

暂时不要实现过度复杂的功能，例如邮件账单自动解析、远程 Docker 控制、真实支付系统、复杂团队协作。

---

## 技术栈要求

### 后端

使用：

* Go
* Gin
* GORM
* PostgreSQL
* JWT
* RESTful API
* config.yaml
* .env
* Docker

后端需要采用清晰的分层结构：

```text
cmd/
internal/
  api/
  middleware/
  service/
  repository/
  model/
  dto/
  config/
  router/
  task/
pkg/
migrations/
docs/
```

### 前端

使用：

* Vue 3
* Vite
* TypeScript
* Pinia
* Vue Router
* Axios
* Naive UI 或 Element Plus
* ECharts

前端需要支持：

* PC 浏览器访问
* 桌面端布局

暂时不做移动端浏览器体验、PWA 或原生手机 App。移动端相关工作等用户明确提及时再单独设计和排期。

### 部署

必须提供：

```text
docker-compose.yml
.env.example
config.example.yaml
README.md
```

项目需要可以通过 Docker Compose 一键启动。

---

## 数据同步设计

PC 端前端通过同一个后端服务和数据库读写数据。

所有数据必须通过后端 API 读写，前端不允许单独存储核心业务数据。

### 同步原则

* 后端数据库是唯一真实数据源
* 前端只做缓存，不做长期数据源
* 所有写操作都必须经过鉴权
* 每条用户数据都必须包含 user_id
* 所有查询都必须按当前登录用户过滤
* 不允许出现跨用户数据访问

PC 端页面重点优化：

* 表格管理
* 筛选搜索
* 批量查看
* 数据统计图表
* 资产详情页

---

## 第一版功能范围

### 1. 用户模块

必须实现：

* 用户注册
* 用户登录
* JWT 鉴权
* 获取当前用户信息
* 修改基础资料
* 修改密码

暂时不做：

* 第三方登录
* 邮箱验证码
* 找回密码
* 双因素认证

---

### 2. 订阅管理模块

订阅字段至少包括：

```text
id
user_id
name
description
category_id
amount
currency
billing_cycle
billing_interval
next_payment_date
start_date
payment_method
auto_renew
status
website_url
remark
created_at
updated_at
```

billing_cycle 支持：

```text
weekly
monthly
quarterly
yearly
custom
one_time
```

status 支持：

```text
active
paused
cancelled
expired
```

必须实现：

* 新增订阅
* 编辑订阅
* 删除订阅
* 查询订阅列表
* 查询订阅详情
* 按分类筛选
* 按状态筛选
* 按即将续费筛选
* 自动计算下一次付款时间
* 统计本月预计支出
* 统计今年预计支出

---

### 3. 分类模块

分类字段至少包括：

```text
id
user_id
name
type
color
icon
sort_order
created_at
updated_at
```

type 支持：

```text
subscription
asset
```

系统需要为新用户初始化默认分类：

* AI 工具
* 开发工具
* 云服务
* 域名
* 服务器
* 娱乐
* 办公
* 其他

用户可以自定义分类。

---

### 4. 数字资产模块

资产字段至少包括：

```text
id
user_id
name
asset_type
provider
identifier
url
expire_date
cost_amount
cost_currency
billing_cycle
status
description
remark
created_at
updated_at
```

asset_type 支持：

```text
domain
server
docker_service
ssl_certificate
api_key
repository
other
```

status 支持：

```text
active
inactive
expired
warning
```

必须实现：

* 新增资产
* 编辑资产
* 删除资产
* 查询资产列表
* 查询资产详情
* 按资产类型筛选
* 按到期时间筛选
* 30 天内到期提醒
* 资产与订阅的关联关系

---

### 5. 仪表盘模块

首页需要返回以下统计数据：

```text
本月预计支出
今年预计支出
订阅总数
资产总数
7 天内即将续费数量
30 天内即将到期资产数量
异常资产数量
最近新增订阅
最近新增资产
即将续费列表
即将到期列表
分类支出占比
```

PC 端使用卡片 + 图表 + 表格。

---

### 6. 提醒模块

第一版先只做站内提醒，不接外部通知。

提醒类型包括：

```text
subscription_renewal
asset_expiration
service_warning
```

提醒字段至少包括：

```text
id
user_id
target_type
target_id
remind_type
remind_date
title
content
is_read
created_at
updated_at
```

需要支持：

* 生成提醒
* 查询提醒列表
* 标记已读
* 删除提醒

---

## API 设计要求

接口风格使用 RESTful API。

统一前缀：

```text
/api/v1
```

接口返回格式统一：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

分页格式统一：

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "page_size": 20
}
```

错误码需要统一管理。

禁止在业务代码中随意返回不一致的数据结构。

---

## 数据库要求

使用 PostgreSQL。

必须提供 migration 文件。

所有表必须包含：

```text
id
created_at
updated_at
```

涉及用户数据的表必须包含：

```text
user_id
```

所有查询必须严格根据 user_id 过滤。

---

## 安全要求

必须注意：

* 密码必须加密存储
* JWT 密钥从环境变量读取
* 不允许把密钥写死在代码中
* API Key 类资产默认不明文展示
* 删除操作需要软删除或二次确认
* 后端必须做参数校验
* 前端不能依赖隐藏按钮来保证权限

---

## 国际化要求

第一版至少预留国际化结构。

默认支持：

```text
zh-CN
en-US
```

前端所有菜单、按钮、提示文案都不要写死在页面里。

后端错误信息可以先返回统一错误码，前端负责翻译。

---

## 前端页面规划

### PC 端页面

需要实现：

```text
登录页
注册页
仪表盘
订阅列表
订阅详情
新增/编辑订阅
资产列表
资产详情
新增/编辑资产
分类管理
提醒中心
用户设置
```

## UI 风格要求

整体风格：

* 简洁
* 清晰
* 偏技术感
* 不要过度花哨
* 数据展示优先

PC 端适合使用侧边栏布局。

---

## 开发优先级

### 第一阶段

完成基础工程结构：

* 后端项目初始化
* 前端项目初始化
* Docker Compose
* 数据库连接
* 配置读取
* 日志系统
* 基础路由

### 第二阶段

完成用户系统：

* 注册
* 登录
* JWT
* 当前用户信息
* 鉴权中间件

### 第三阶段

完成核心业务：

* 订阅管理
* 分类管理
* 资产管理
* 提醒管理

### 第四阶段

完成仪表盘：

* 本月支出
* 年度支出
* 即将续费
* 即将到期
* 分类统计

### 第五阶段

完善部署和文档：

* docker-compose.yml
* .env.example
* config.example.yaml
* README.md
* API 文档
* 初始化数据说明

### 后续阶段

移动端浏览器体验、PWA、原生手机 App 暂不纳入当前计划。等用户明确提出移动端需求时，再重新设计移动端信息架构、导航、表单和首页体验。

---

## 暂不实现的功能

以下功能不要在第一版实现：

* 邮件账单解析
* 支付系统
* 银行卡绑定
* 真实扣费
* 远程 Docker 管理
* 服务器 SSH 管理
* 团队空间
* 第三方登录
* 手机原生 App
* 小程序
* 浏览器插件

可以预留扩展点，但不要提前复杂化。

---

## 代码质量要求

开发时必须遵守：

* 代码结构清晰
* 命名语义明确
* 不写无意义注释
* 不写重复代码
* 不把业务逻辑写在控制器里
* 参数校验集中处理
* 错误处理统一
* 接口返回统一
* 前后端类型尽量一致

---

## Git 要求

默认分支使用：

```text
master
```

提交信息尽量清晰，例如：

```text
feat: add subscription module
fix: fix jwt middleware
refactor: clean dashboard service
docs: update deployment guide
```

---

## 你作为 AI 编程助手的工作方式

当我提出需求时，你需要：

1. 先理解需求
2. 拆分任务
3. 判断是否影响已有结构
4. 给出实现方案
5. 再修改代码
6. 修改后说明改动点
7. 如果有风险，需要明确指出

不要在没有确认项目结构的情况下大面积重构。

不要随意更换技术栈。

不要引入不必要的复杂依赖。

优先保证项目可以运行。

---

## 当前第一步任务

请先完成项目初始化。

需要生成：

```text
backend/
frontend/
docker-compose.yml
.env.example
config.example.yaml
README.md
CLAUDE.md
```

后端使用 Go 初始化。

前端使用 Vue 3 + Vite + TypeScript 初始化。

Docker Compose 中至少包含：

* backend
* frontend
* postgres

完成后请说明：

* 项目目录结构
* 启动方式
* 后续开发建议
