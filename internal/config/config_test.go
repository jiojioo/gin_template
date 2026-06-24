package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMustLoadDecodesAllConfigurationSections(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")
	contents := []byte(`server:
  name: test-app
  mode: test
  addr: ":8081"
mysql:
  dsn: "user:pass@tcp(localhost:3306)/test"
  max_idle_conns: 2
  max_open_conns: 4
  conn_max_lifetime: 60
redis:
  addr: "localhost:6379"
  password: "redis-pass"
  db: 3
jwt:
  secret: "test-secret"
  expire: 7200
log:
  path: "./logs"
  level: "debug"
  keep_hours: 4
  filename: "app.log"
  access_filename: "access.log"
  error_filename: "error.log"
`)
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg := MustLoad(path)
	if cfg.Server.Name != "test-app" || cfg.Server.Addr != ":8081" {
		t.Fatalf("unexpected server config: %#v", cfg.Server)
	}
	if cfg.MySQL.MaxOpenConns != 4 || cfg.Redis.DB != 3 {
		t.Fatalf("unexpected data store config: mysql=%#v redis=%#v", cfg.MySQL, cfg.Redis)
	}
	if cfg.JWT.Secret != "test-secret" || cfg.JWT.Expire != 7200 {
		t.Fatalf("unexpected jwt config: %#v", cfg.JWT)
	}
	if cfg.Log.AccessFilename != "access.log" || cfg.Log.ErrorFilename != "error.log" {
		t.Fatalf("unexpected log config: %#v", cfg.Log)
	}
}
