## Context

项目需要从一个轻量 Gin 应用演进为可复用的服务端模板。当前变更建立首个可运行的基线：配置、基础设施、日志、公共安全工具和用户认证样板必须遵守同一分层边界，并能作为后续模块生成器的输入样本。

约束：保持 Gin 的轻量运行时；使用 Go 模块管理依赖；服务层不得依赖 `gin.Context`；数据库与 Redis 连接在启动期初始化；首期仅覆盖 user 模块的登录和用户信息读取。

## Goals / Non-Goals

**Goals:**

- 固化 `router -> handler -> service -> repo -> model/db` 调用方向，并通过 Service/Handler 容器完成依赖装配。
- 在启动期依次加载配置、初始化日志、MySQL、Redis、迁移/基础数据，并装配路由。
- 提供一致的 JSON 响应、JWT bearer 鉴权、bcrypt 密码校验、雪花 ID 与按小时切割的业务/Gin 日志。
- 提供登录与用户信息接口作为完整可测试的纵向样板。

**Non-Goals:**

- 不实现通用 CRUD、模块模板、CLI 生成器或数据库反向生成。
- 不实现 Swagger、RBAC/Casbin、限流、OpenTelemetry、Prometheus 或 Docker Compose。
- 不在首期抽象统一业务错误码、事务管理器或通用缓存框架。

## Decisions

### 1. 采用单向分层与容器装配

- Router 只注册路由与中间件，依赖 Handler 容器。
- Handler 负责绑定、基础校验和 HTTP 响应，只调用 Service。
- Service 承载业务规则、JWT 签发及未来事务/缓存编排，并使用标准 `context.Context`。
- Repo 封装 GORM 查询，Model 仅描述持久化结构。
- `service.NewService` 统一注入 Repo 与 Redis；`handler.NewHandler` 统一注入 Service。

选择容器装配而非在每个 Handler 内即时创建依赖，是为了使生成器只需追加固定注册点，并避免隐藏连接与跨层依赖。

### 2. 启动期集中初始化基础设施

`main` 以失败即退出的方式完成：配置 → Zap → Gin Writer → MySQL → Redis → 迁移/基础数据 → Router → HTTP 服务。MySQL 和 Redis 分别暴露受控的全局客户端，Repo/Service 在容器创建时接收它们。

选择启动期 fail-fast 而非首次请求懒初始化，避免在流量进入后才暴露配置或连接错误；测试可通过替换容器依赖绕开真实连接。

### 3. 统一安全与 HTTP 边界

- 用户密码只保存 bcrypt 哈希；登录由 Service 校验哈希并签发带 `user_id` 的 JWT。
- Auth 中间件仅解析 `Authorization: Bearer <token>`，将 `user_id` 写入 Gin 上下文；业务查询仍由 Handler 调用 Service 完成。
- `pkg/response` 以 `{code,message,data}` 输出；Handler 负责选择 HTTP 状态，Service 返回领域错误。

选择 JWT bearer token 而非服务端 session，避免首期认证能力依赖 Redis；Redis 保留为后续 token 黑名单和缓存扩展点。

### 4. 统一结构化日志并接管 Gin Writer

Zap 负责业务结构化日志；自定义按小时写入器根据配置保留最近 `keep_hours` 的日志。初始化时设置 `gin.DefaultWriter` 和 `gin.DefaultErrorWriter`，使 access/error 日志与业务日志分流且遵守相同保留策略。

选择自定义按小时切割器而非单文件日志，以满足短保留时间和按小时排障的模板需求；未来可替换为外部日志采集或成熟 rotation 库而不改变调用方。

### 5. 将 user 模块作为纵向样板

用户表包含 `uint64` 主键、用户名、密码哈希、昵称、状态和审计时间。路由只暴露 `POST /api/v1/user/login` 和受 Auth 保护的 `GET /api/v1/user/info`。用户注册/完整 CRUD 不属于当前变更，初始化数据入口只定义策略和可安全重复执行的实现边界。

这比直接实现所有 CRUD 更适合验证分层、鉴权和依赖装配；后续生成器从已稳定的模块样板抽象模板。

## Risks / Trade-offs

- [本地没有 MySQL/Redis 时服务无法启动] → 在 README 与配置中明确前置依赖；测试使用替代依赖或容器化测试环境。
- [全局 DB/Redis 客户端降低可测试性] → 容器构造函数接收 Repo/Redis 依赖，测试避免读取全局变量。
- [JWT secret 配置不当造成令牌风险] → 默认示例只用于开发，生产环境要求外部安全注入和足够长度的 secret。
- [启动期 AutoMigrate 影响生产 schema] → 首期保留该模板能力，并在部署说明中要求生产环境审查迁移；未来可切换为显式迁移工具。
- [小时切割实现处理并发或清理不当] → Writer 必须同步保护切换与清理，测试覆盖小时边界和保留时间。

## Migration Plan

1. 添加配置与公共包，保持现有 API 不被静默改写。
2. 引入数据库、日志和分层容器，创建 user 表迁移与基础数据入口。
3. 注册认证路由并补充单元及集成测试。
4. 在开发环境以 `etc/config.yaml` 启动并验证日志、登录、鉴权查询。

回滚策略：撤销该 change 的代码和迁移；若 user 表已创建，保留数据表而不由应用自动删除，避免不可逆数据丢失。

## Open Questions

- 初始管理员账号应由配置提供、环境变量注入，还是仅提供空的幂等初始化钩子？
- 开发测试是否要求内置 Docker Compose，还是仅在 README 说明 MySQL/Redis 前置条件？
- 雪花 ID 的节点 ID 应配置化还是首期固定为单机默认值？
