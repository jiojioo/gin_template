package logger

import (
	"fmt"
	"strings"

	"github.com/jiojioo/gin_template/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

func Init(cfg config.LogConfig) error {
	writer, err := NewRotateWriter(config.LogConfig{
		Path:      cfg.Path,
		Filename:  cfg.Filename,
		KeepHours: cfg.KeepHours,
	}, nil)
	if err != nil {
		return err
	}

	level, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(writer), level)
	Log = zap.New(core, zap.AddCaller())
	return nil
}

func Sync() error {
	if Log == nil {
		return nil
	}
	return Log.Sync()
}

func parseLevel(raw string) (zapcore.Level, error) {
	if strings.TrimSpace(raw) == "" {
		return zapcore.InfoLevel, nil
	}
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(raw))); err != nil {
		return zapcore.InfoLevel, fmt.Errorf("parse log level %q: %w", raw, err)
	}
	return level, nil
}
