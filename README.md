# Gin Template

Runnable Gin service template with YAML configuration, fail-fast MySQL/Redis startup, structured logging, shared JSON responses, bcrypt password checks, JWT authentication, and a minimal user login/profile vertical slice.

## Prerequisites

- Go 1.24-compatible toolchain.
- MySQL reachable by `mysql.dsn`.
- Redis reachable by `redis.addr`.

The application initializes MySQL and Redis before starting HTTP listening. If either dependency fails, startup fails.

## Configuration

Default configuration lives in `etc/config.yaml`.

Important production settings:

- Replace `jwt.secret` with a long secret from a secure external source.
- Set `server.mode` to `release` outside development.
- Review `mysql.dsn`, Redis credentials, and log paths for the target environment.

## Run

```bash
go run .
```

The server listens on `server.addr`.

## API examples

Login:

```bash
curl -X POST http://127.0.0.1:8080/api/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret"}'
```

Get current user profile:

```bash
curl http://127.0.0.1:8080/api/v1/user/info \
  -H "Authorization: Bearer <token>"
```

Responses use:

```json
{"code":0,"message":"success","data":{}}
```

## Verification

```bash
go test ./...
go vet ./...
```
