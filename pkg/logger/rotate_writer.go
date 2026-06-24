// Package logger initializes application and Gin log writers.
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jiojioo/gin_template/internal/config"
)

type Clock func() time.Time

type RotateWriter struct {
	mu          sync.Mutex
	cfg         config.LogConfig
	now         Clock
	currentHour string
	file        *os.File
}

func NewRotateWriter(cfg config.LogConfig, now Clock) (*RotateWriter, error) {
	if now == nil {
		now = time.Now
	}
	if strings.TrimSpace(cfg.Path) == "" {
		return nil, fmt.Errorf("log.path is required")
	}
	if strings.TrimSpace(cfg.Filename) == "" {
		return nil, fmt.Errorf("log.filename is required")
	}
	if err := os.MkdirAll(cfg.Path, 0o755); err != nil {
		return nil, err
	}
	return &RotateWriter{cfg: cfg, now: now}, nil
}

func (w *RotateWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	hour := w.now().Format("2006010215")
	if w.file == nil || w.currentHour != hour {
		if err := w.rotate(hour); err != nil {
			return 0, err
		}
	}
	return w.file.Write(p)
}

func (w *RotateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

func (w *RotateWriter) rotate(hour string) error {
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return err
		}
	}
	name := fmt.Sprintf("%s.%s", w.cfg.Filename, hour)
	file, err := os.OpenFile(filepath.Join(w.cfg.Path, name), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	w.file = file
	w.currentHour = hour
	w.cleanupLocked()
	return nil
}

func (w *RotateWriter) cleanupLocked() {
	if w.cfg.KeepHours <= 0 {
		return
	}
	cutoff := w.now().Add(-time.Duration(w.cfg.KeepHours) * time.Hour)
	entries, err := os.ReadDir(w.cfg.Path)
	if err != nil {
		return
	}
	prefix := w.cfg.Filename + "."
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, prefix) {
			continue
		}
		suffix := strings.TrimPrefix(name, prefix)
		t, err := time.Parse("2006010215", suffix)
		if err == nil && t.Before(cutoff) {
			_ = os.Remove(filepath.Join(w.cfg.Path, name))
		}
	}
}
