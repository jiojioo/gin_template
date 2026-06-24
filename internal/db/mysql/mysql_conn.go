// Package mysql owns the global GORM MySQL connection.
package mysql

import (
	"database/sql"
	"time"

	"github.com/jiojioo/gin_template/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DSN                    string
	MaxIdleConns           int
	MaxOpenConns           int
	ConnMaxLifetimeSeconds int
}

var Client *gorm.DB

func OpenConfig(cfg config.MySQLConfig) Config {
	return Config{
		DSN:                    cfg.DSN,
		MaxIdleConns:           cfg.MaxIdleConns,
		MaxOpenConns:           cfg.MaxOpenConns,
		ConnMaxLifetimeSeconds: cfg.ConnMaxLifetime,
	}
}

func Init(cfg config.MySQLConfig) error {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := configurePool(db, OpenConfig(cfg)); err != nil {
		return err
	}
	Client = db
	return nil
}

func configurePool(db *gorm.DB, cfg Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	applyPool(sqlDB, cfg)
	return nil
}

func applyPool(db *sql.DB, cfg Config) {
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetimeSeconds > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)
	}
}
