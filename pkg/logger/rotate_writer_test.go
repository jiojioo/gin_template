package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jiojioo/gin_template/internal/config"
)

func TestRotateWriterCreatesHourlyFile(t *testing.T) {
	dir := t.TempDir()
	writer, err := NewRotateWriter(config.LogConfig{
		Path:      dir,
		Filename:  "app.log",
		KeepHours: 4,
	}, func() time.Time {
		return time.Date(2026, 6, 24, 15, 30, 0, 0, time.UTC)
	})
	if err != nil {
		t.Fatalf("NewRotateWriter() error = %v", err)
	}
	t.Cleanup(func() {
		if err := writer.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})

	if _, err := writer.Write([]byte("hello\n")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("created files = %d, want 1", len(entries))
	}
	name := entries[0].Name()
	if !strings.Contains(name, "app.log.2026062415") {
		t.Fatalf("created file %q, want hourly app log", name)
	}
	contents, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(contents) != "hello\n" {
		t.Fatalf("file contents = %q, want hello newline", string(contents))
	}
}
