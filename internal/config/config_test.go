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

func TestMustLoadPanicsWhenRequiredFieldsAreMissing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents string
	}{
		{
			name: "server addr",
			contents: `server:
  name: test-app
  mode: test
mysql:
  dsn: "user:pass@tcp(localhost:3306)/test"
redis:
  addr: "localhost:6379"
jwt:
  secret: "test-secret"
`,
		},
		{
			name: "mysql dsn",
			contents: `server:
  addr: ":8081"
mysql:
  max_idle_conns: 2
redis:
  addr: "localhost:6379"
jwt:
  secret: "test-secret"
`,
		},
		{
			name: "redis addr",
			contents: `server:
  addr: ":8081"
mysql:
  dsn: "user:pass@tcp(localhost:3306)/test"
redis:
  db: 0
jwt:
  secret: "test-secret"
`,
		},
		{
			name: "jwt secret",
			contents: `server:
  addr: ":8081"
mysql:
  dsn: "user:pass@tcp(localhost:3306)/test"
redis:
  addr: "localhost:6379"
jwt:
  expire: 7200
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(path, []byte(tt.contents), 0o600); err != nil {
				t.Fatal(err)
			}

			defer func() {
				if recover() == nil {
					t.Fatalf("MustLoad() did not panic for missing %s", tt.name)
				}
			}()

			_ = MustLoad(path)
		})
	}
}
