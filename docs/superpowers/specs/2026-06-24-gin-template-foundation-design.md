---
comet_change: build-core-gin-template
role: technical-design
canonical_spec: openspec
---

# Gin Template Foundation Design

## Context

当前目录不含既有 Go 应用源码。本 change 将以当前目录为项目根目录，创建一个可运行的 Gin 后端模板，并把首个 user 认证模块作为后续模块生成器的纵向样板。

## Goals

- 建立 `router -> handler -> service -> repo -> model/db` 的固定单向分层。
- 提供配置、MySQL、Redis、Zap、Gin 日志、JWT、bcrypt、统一响应和雪花 ID。
- 交付用户登录与受保护的用户信息查询接口。
- 使业务依赖可通过构造函数替换，以支持不连接真实基础设施的单元测试。

## Non-Goals

- 不实现用户注册、完整 CRUD、代码生成器或数据库表反向生成。
- 不实现 Swagger、RBAC、限流、可观测性、Docker Compose、统一业务错误码或事务管理器。
- 不自动创建初始管理员账号。

## Architecture

### Startup and infrastructure

`main` 按以下顺序执行：加载 `etc/config.yaml`、初始化 Zap、初始化 Gin access/error writer、初始化 MySQL、执行迁移和 `InitData` 钩子、初始化 Redis、装配路由并启动 HTTP 服务。

MySQL 与 Redis 是启动期真实依赖；配置不合法或任一依赖不可用时，进程必须在监听 HTTP 前退出。`InitData` 必须幂等，但首期实现为空钩子，不创建管理员或默认业务数据。

MySQL `*gorm.DB` 和 Redis `*redis.Client` 在初始化后作为包级受控客户端提供给装配层。Service/Handler 容器仍使用构造函数显式注入 Repo、Redis 依赖和 JWT 配置，避免业务对象自行查找全局状态。测试通过替换 Repo 接口与最小 Redis 依赖绕开真实服务。

### Layer responsibilities

- Router：创建 Gin Engine、挂载全局中间件、注册模块路由；不得处理业务、绑定参数或访问数据。
- Handler：绑定和基础校验 HTTP 输入，读取认证上下文，调用 Service，并生成统一响应；不得访问 Repo、GORM 或 Redis。
- Service：实现业务规则、密码校验、JWT 签发和输出 DTO 组装；只接收标准 `context.Context`，不得依赖 `gin.Context` 或写 HTTP 响应。
- Repo：封装 GORM 查询，仅负责持久化操作。
- Model：定义 GORM 映射和安全 JSON 标签。

`Service` 容器负责组装全部业务 Service，`Handler` 容器负责组装模块 Handler。Router 只依赖 Handler 容器。

### Identity and authentication

用户模型包含 `uint64` ID、用户名、bcrypt 密码哈希、昵称、状态及创建/更新时间。雪花 ID 使用固定单机节点 `1`；节点 ID 配置化不属于当前 change。

`POST /api/v1/user/login` 接收必填的用户名和密码。Service 使用 Repo 查询用户、使用 bcrypt 对比密码，并生成带 `user_id` 和配置过期时间的 JWT。失败时不区分用户不存在与密码不匹配。已接受残留：该「不区分」仅作用于响应内容；bcrypt 校验仅在用户存在时执行，缺失用户时提前返回，因此仍存在基于响应耗时的用户枚举侧信道，首期不处理。

Auth 中间件仅接受 `Authorization: Bearer <token>`，验证 JWT 并写入 `gin.Context["user_id"]`。`GET /api/v1/user/info` 必须受该中间件保护，并只返回 `id`、`username`、`nickname` 和 `status`。

### Responses and error boundaries

成功响应固定为 `{code: 0, message: "success", data: ...}`。首期失败响应的 `code` 使用 HTTP 状态码；Handler 负责将错误映射为 HTTP 400、401 或 500。Service 与 Repo 仅返回错误，不产生 HTTP 响应。JWT secret 和过期时间来自配置；示例 secret 仅限开发环境。

### Logging and migration

Zap 记录业务日志。Gin access 与 error 日志使用独立文件。三类文件均通过线程安全写入器按小时轮换，并在轮换时清理超过 `keep_hours` 的文件；开发模式可额外写入控制台。

MySQL 初始化完成后执行注册模型的 `AutoMigrate`。生产 schema 变更在部署前必须审查；回滚仅撤回应用代码，不自动删除已创建的数据表。

## Testing

- 公共包：配置解析、JWT、Hash、统一响应、雪花 ID 和日志小时切换。
- 数据层：Repo 查询使用可控测试数据库或 GORM 测试策略。
- 业务层：用替身 Repo 验证有效/无效凭据与安全 DTO。
- HTTP：验证登录成功、字段缺失、凭据无效、缺失/无效 bearer token、以及用户信息不含密码字段。

## Risks and mitigations

- 本地没有 MySQL/Redis 会阻止服务启动：README 必须列出前置条件，测试不依赖真实服务。
- 包级客户端可能被滥用：只允许启动装配层访问，业务对象必须接收构造参数。
- 示例 JWT secret 被误用于生产：配置注释和 README 明确要求生产环境注入高熵 secret。
- 自定义日志轮换易出现并发或清理问题：写入器必须同步保护，测试覆盖并发写入、小时边界和保留策略。

## Spec patch

该设计不改变已确认的 OpenSpec delta specs；它将已确认的运行时、认证、测试和依赖策略具体化。
