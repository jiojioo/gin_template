## Why

原生 Gin 项目常将路由、HTTP 处理、业务逻辑和数据访问混杂在一起，导致配置、日志、鉴权和响应约定无法复用，也阻碍后续模块生成。现在需要建立一个可运行且边界明确的基础模板，作为团队项目和自动化代码生成的稳定起点。

## What Changes

- 建立 `router -> handler -> service -> repo -> model/db` 的固定分层和服务、处理器容器。
- 提供 YAML 配置加载，以及 MySQL、Redis、Zap 与 Gin access/error 日志初始化。
- 提供统一 JSON 响应、JWT、bcrypt 密码哈希与雪花 ID 公共能力。
- 增加用户认证样板：登录并签发 JWT，以及受保护的用户信息查询接口。
- 增加用户模型迁移与基础数据初始化入口。

## Capabilities

### New Capabilities

- `gin-template-runtime`: 可配置、可启动的 Gin 服务运行时，包含分层装配、基础设施初始化、日志和通用工具约定。
- `user-authentication`: 用户登录、JWT 鉴权和用户信息查询的标准分层实现。

### Modified Capabilities

无。当前主规格中没有既有 capability。

## Impact

- 新增或重组 `main.go`、`etc`、`internal/config`、`internal/db`、`internal/middleware`、`internal/model`、`internal/repo`、`internal/service`、`internal/handler`、`internal/router` 和 `pkg` 下的基础能力。
- 新增 HTTP API：`POST /api/v1/user/login` 与 `GET /api/v1/user/info`。
- 引入或固化 Gin、GORM/MySQL、go-redis、Zap、JWT、bcrypt、YAML 和雪花 ID 的依赖边界。
- 不包含模块/CRUD 生成器、Swagger、RBAC、可观测性、Docker Compose 或其他平台增强能力。
