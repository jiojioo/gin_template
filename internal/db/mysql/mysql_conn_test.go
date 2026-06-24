package mysql

import (
	"testing"

	"github.com/jiojioo/gin_template/internal/config"
)

func TestOpenConfigAppliesConnectionPoolSettings(t *testing.T) {
	cfg := config.MySQLConfig{
		DSN:             "user:pass@tcp(localhost:3306)/test",
		MaxIdleConns:    2,
		MaxOpenConns:    4,
		ConnMaxLifetime: 60,
	}

	openCfg := OpenConfig(cfg)

	if openCfg.DSN != cfg.DSN {
		t.Fatalf("DSN = %q, want %q", openCfg.DSN, cfg.DSN)
	}
	if openCfg.MaxIdleConns != 2 || openCfg.MaxOpenConns != 4 || openCfg.ConnMaxLifetimeSeconds != 60 {
		t.Fatalf("unexpected pool config: %#v", openCfg)
	}
}
