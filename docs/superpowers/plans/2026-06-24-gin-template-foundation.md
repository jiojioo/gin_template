# Gin Template Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Change:** `build-core-gin-template`

**Design Doc:** `docs/superpowers/specs/2026-06-24-gin-template-foundation-design.md`

**Goal:** Create a runnable Gin template with real startup dependencies and a testable user login/JWT profile vertical slice.

**Architecture:** Initialize configuration, logs, MySQL, migrations and Redis before route assembly. Router injects a Handler container; handlers call services; services use repository interfaces and standard contexts; repositories use GORM only.

**Tech Stack:** Go, Gin, GORM/MySQL, go-redis/v9, Zap, golang-jwt/jwt/v5, bcrypt, YAML.

## Global Constraints

- MySQL and Redis connection failures MUST prevent HTTP listening.
- Handler MUST NOT access GORM, Redis, Repo or Model directly.
- Service MUST NOT depend on `gin.Context` or write HTTP responses.
- User passwords MUST use bcrypt and MUST NOT be serialized.
- The implementation MUST NOT add CRUD generation, Swagger, RBAC, Docker Compose or an initial administrator.

---

### Task 1: Bootstrap module, configuration and public contracts

**Files:**
- Create: `go.mod`, `main.go`, `etc/config.yaml`, `internal/config/config.go`
- Create: `pkg/response/response.go`, `pkg/hash/hash.go`, `pkg/jwt/jwt.go`, `pkg/snowflake/snowflake.go`
- Test: `internal/config/config_test.go`, `pkg/hash/hash_test.go`, `pkg/jwt/jwt_test.go`, `pkg/response/response_test.go`

**Interfaces:**
- Produces: `config.MustLoad(path string) *config.Config`, `hash.Make(string) (string, error)`, `hash.Check(string, string) bool`, `jwt.GenerateToken(uint64) (string, error)`, `jwt.ParseToken(string) (*Claims, error)`, `response.Success(*gin.Context, any)`, `response.Fail(*gin.Context, int, string)`.

- [x] **Step 1: Write failing configuration and utility tests**

```go
func TestPasswordRoundTrip(t *testing.T) {
  encoded, err := hash.Make("secret")
  require.NoError(t, err)
  assert.True(t, hash.Check("secret", encoded))
}
```

- [x] **Step 2: Run tests to verify failure**

Run: `go test ./internal/config ./pkg/hash ./pkg/jwt ./pkg/response`

Expected: FAIL because packages do not exist.

- [x] **Step 3: Implement module, typed config and utilities**

```go
type Config struct { Server ServerConfig `yaml:"server"`; MySQL MySQLConfig `yaml:"mysql"`; Redis RedisConfig `yaml:"redis"`; JWT JWTConfig `yaml:"jwt"`; Log LogConfig `yaml:"log"` }
type Body struct { Code int `json:"code"`; Message string `json:"message"`; Data any `json:"data,omitempty"` }
```

- [x] **Step 4: Run tests to verify pass**

Run: `go test ./internal/config ./pkg/hash ./pkg/jwt ./pkg/response`

Expected: PASS.

- [x] **Step 5: Commit**

Run: `git add go.mod main.go etc internal/config pkg && git commit -m "feat: add template configuration and shared contracts"`

### Task 2: Add logging and startup infrastructure

**Files:**
- Create: `pkg/logger/logger.go`, `pkg/logger/rotate_writer.go`, `pkg/logger/gin_writer.go`
- Create: `internal/db/database.go`, `internal/db/mysql/mysql_conn.go`, `internal/db/mysql/auto_migrate.go`, `internal/db/mysql/init_data.go`, `internal/db/redis/redis_conn.go`
- Test: `pkg/logger/rotate_writer_test.go`, `internal/db/mysql/mysql_conn_test.go`

**Interfaces:**
- Produces: `logger.Init(config.LogConfig) error`, `logger.InitGinWriter(config.LogConfig) error`, `mysql.Init(config.MySQLConfig) error`, `redis.Init(config.RedisConfig) error`, `db.Init(*config.Config) error`.

- [x] **Step 1: Write failing rotation and startup-order tests**

```go
func TestRotateWriterCreatesHourlyFile(t *testing.T) { /* inject clock; assert hourly file name */ }
```

- [x] **Step 2: Run tests to verify failure**

Run: `go test ./pkg/logger ./internal/db/...`

Expected: FAIL because logger and database packages do not exist.

- [x] **Step 3: Implement fail-fast clients and writers**

```go
func Init(cfg *config.Config) error {
  if err := mysql.Init(cfg.MySQL); err != nil { return err }
  if err := mysql.AutoMigrate(); err != nil { return err }
  if err := mysql.InitData(); err != nil { return err }
  return redis.Init(cfg.Redis)
}
```

- [x] **Step 4: Run tests to verify pass**

Run: `go test ./pkg/logger ./internal/db/...`

Expected: PASS.

- [x] **Step 5: Commit**

Run: `git add pkg/logger internal/db && git commit -m "feat: add logging and startup infrastructure"`

### Task 3: Implement model, repository and service boundaries

**Files:**
- Create: `internal/model/base_model.go`, `internal/model/user_model.go`, `internal/repo/user_repo.go`, `internal/repo/repo.go`, `internal/service/service.go`, `internal/service/user_service.go`
- Test: `internal/service/user_service_test.go`

**Interfaces:**
- Produces: `repo.UserRepository`, `service.NewUserService(repo.UserRepository, *redis.Client, jwt.Config)`, `(*UserService).Login(context.Context, *LoginReq) (*LoginResp, error)`, `(*UserService).GetUserInfo(context.Context, uint64) (*GetUserInfoResp, error)`.

- [x] **Step 1: Write failing Service tests with fake repository**

```go
type fakeUsers struct { user *model.User; err error }
func (f fakeUsers) FindByUsername(context.Context, string) (*model.User, error) { return f.user, f.err }
```

- [x] **Step 2: Run tests to verify failure**

Run: `go test ./internal/service`

Expected: FAIL because service and repository interfaces do not exist.

- [x] **Step 3: Implement GORM repository, DTOs and container**

```go
type UserRepository interface { FindByID(context.Context, uint64) (*model.User, error); FindByUsername(context.Context, string) (*model.User, error) }
```

- [x] **Step 4: Run tests to verify pass**

Run: `go test ./internal/service ./internal/repo`

Expected: PASS.

- [x] **Step 5: Commit**

Run: `git add internal/model internal/repo internal/service && git commit -m "feat: add user data and service layers"`

### Task 4: Implement middleware, handlers and routes

**Files:**
- Create: `internal/middleware/auth.go`, `internal/middleware/cors.go`, `internal/handler/handler.go`, `internal/handler/v1/user/user_handler.go`, `internal/handler/v1/user/login.go`, `internal/handler/v1/user/get_user_info.go`, `internal/router/router.go`, `internal/router/v1/user_router.go`
- Test: `internal/router/router_test.go`

**Interfaces:**
- Consumes: `service.UserService`; produces `router.InitRouter(*config.Config, *handler.Handler) *gin.Engine` and routes `/api/v1/user/login`, `/api/v1/user/info`.

- [x] **Step 1: Write failing HTTP tests**

```go
req := httptest.NewRequest(http.MethodPost, "/api/v1/user/login", strings.NewReader(`{"username":"u","password":"p"}`))
req.Header.Set("Content-Type", "application/json")
```

- [x] **Step 2: Run tests to verify failure**

Run: `go test ./internal/router`

Expected: FAIL because routes are not registered.

- [x] **Step 3: Implement HTTP boundary and JWT middleware**

```go
func Auth() gin.HandlerFunc { return func(c *gin.Context) { /* parse Bearer JWT; set user_id; abort 401 on failure */ } }
```

- [x] **Step 4: Run tests to verify pass**

Run: `go test ./internal/handler/... ./internal/middleware ./internal/router`

Expected: PASS.

- [x] **Step 5: Commit**

Run: `git add internal/middleware internal/handler internal/router && git commit -m "feat: add user authentication HTTP flow"`

### Task 5: Wire entrypoint, verify and document operation

**Files:**
- Modify: `main.go`, `README.md`
- Test: `main_test.go`, `internal/router/router_test.go`

**Interfaces:**
- Consumes: `config.MustLoad`, logger initialization, `db.Init`, `router.InitRouter`.

- [ ] **Step 1: Write failing startup composition test**

```go
func TestApplicationBuildsRouterAfterDependencies(t *testing.T) { /* inject initialized test dependencies; assert routes */ }
```

- [ ] **Step 2: Run tests to verify failure**

Run: `go test ./...`

Expected: FAIL until entrypoint composition and remaining dependencies are complete.

- [ ] **Step 3: Implement startup composition and operational documentation**

```go
cfg := config.MustLoad("etc/config.yaml")
if err := db.Init(cfg); err != nil { panic(err) }
router.InitRouter(cfg).Run(cfg.Server.Addr)
```

- [ ] **Step 4: Run final verification**

Run: `gofmt -w .; go test ./...; go vet ./...`

Expected: all commands exit 0.

- [ ] **Step 5: Commit**

Run: `git add main.go README.md && git commit -m "docs: document gin template startup"`

## Coverage review

- Startup/configuration, layers, infrastructure, migration, logs and shared contracts: Tasks 1-2 and 5.
- Login, JWT authentication and safe user profile response: Tasks 3-4.
- Required unit and HTTP verification: Tasks 1-5.
