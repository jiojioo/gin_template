## 1. 项目骨架与配置

- [x] 1.1 建立目标目录结构，并整理 `main.go`、`etc/config.yaml` 与项目依赖。
- [x] 1.2 实现配置结构与 YAML 加载；为必需配置提供启动期校验和错误信息。
- [x] 1.3 实现统一响应、随机工具、雪花 ID、bcrypt 密码哈希和 JWT 公共包。

## 2. 日志与基础设施

- [x] 2.1 实现 Zap 业务日志初始化、级别解析和同步关闭。
- [x] 2.2 实现线程安全的按小时日志写入与保留期清理，并配置 Gin access/error writer。
- [x] 2.3 实现 MySQL 初始化、连接池配置、模型迁移和可幂等的基础数据入口。
- [x] 2.4 实现 Redis 初始化、连接检测与客户端暴露。
- [x] 2.5 实现统一 db 初始化流程，并确保任一基础设施失败时服务不开始监听。

## 3. 分层装配与 HTTP 基线

- [x] 3.1 定义基础模型和用户模型，包含安全的 JSON 字段与 GORM 表映射。
- [x] 3.2 实现 `UserRepo` 的按 ID、按用户名和写入数据访问方法，并只接受标准 context。
- [x] 3.3 实现 Service 容器和 Handler 容器，明确依赖注入边界。
- [x] 3.4 实现 CORS、JWT Auth 和 Recovery/日志中间件装配。
- [x] 3.5 实现总路由与 user 路由注册，确保 Router 不承载业务逻辑。

## 4. 用户认证样板

- [x] 4.1 实现 `UserService.Login`：查询用户、校验 bcrypt 密码并签发带 user ID 的 JWT。
- [x] 4.2 实现 `UserService.GetUserInfo`：查询用户并仅组装安全的个人资料字段。
- [x] 4.3 实现登录 Handler 的 JSON 绑定、校验和统一成功/失败响应。
- [x] 4.4 实现用户信息 Handler，并通过 Auth 中间件读取 `user_id`。
- [x] 4.5 注册 `POST /api/v1/user/login` 与受保护的 `GET /api/v1/user/info`。

## 5. 验证与文档

- [x] 5.1 为配置、JWT、Hash、响应和小时日志切割编写单元测试。
- [x] 5.2 为 Repo/Service 及登录、鉴权用户信息链路编写可重复执行的测试。
- [x] 5.3 执行格式化、静态检查和测试；修复所有回归。
- [x] 5.4 更新 README，说明配置、MySQL/Redis 前置条件、启动方式、API 示例与安全配置要求。

## 6. 验证修复（verify-fail 回退）

- [x] 6.1 修复登录必填字段校验：为 `LoginReq` 添加 `binding:"required"`，确保省略 username/password 返回 HTTP 400（验证报告 WARNING 1）。
- [x] 6.2 修复错误响应信息泄露：repo 将 `gorm.ErrRecordNotFound` 翻译为领域错误，service 不区分用户不存在与密码错误，handler 返回通用消息且不暴露持久层细节（验证报告 WARNING 2）。
