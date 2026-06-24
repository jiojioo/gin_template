package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/jiojioo/gin_template/internal/config"
)

func InitGinWriter(cfg config.LogConfig) error {
	access, err := NewRotateWriter(config.LogConfig{
		Path:      cfg.Path,
		Filename:  cfg.AccessFilename,
		KeepHours: cfg.KeepHours,
	}, nil)
	if err != nil {
		return err
	}
	errorWriter, err := NewRotateWriter(config.LogConfig{
		Path:      cfg.Path,
		Filename:  cfg.ErrorFilename,
		KeepHours: cfg.KeepHours,
	}, nil)
	if err != nil {
		return err
	}
	gin.DefaultWriter = access
	gin.DefaultErrorWriter = errorWriter
	return nil
}
