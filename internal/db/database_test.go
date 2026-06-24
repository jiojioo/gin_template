package db

import (
	"errors"
	"testing"

	"github.com/jiojioo/gin_template/internal/config"
)

func TestInitRunsInfrastructureInStartupOrder(t *testing.T) {
	cfg := &config.Config{}
	var calls []string
	restore := replaceInitHooks(
		func(config.MySQLConfig) error {
			calls = append(calls, "mysql")
			return nil
		},
		func() error {
			calls = append(calls, "migrate")
			return nil
		},
		func() error {
			calls = append(calls, "data")
			return nil
		},
		func(config.RedisConfig) error {
			calls = append(calls, "redis")
			return nil
		},
	)
	t.Cleanup(restore)

	if err := Init(cfg); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	want := []string{"mysql", "migrate", "data", "redis"}
	if len(calls) != len(want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("calls = %#v, want %#v", calls, want)
		}
	}
}

func TestInitStopsWhenInfrastructureFails(t *testing.T) {
	failure := errors.New("mysql down")
	var calls []string
	restore := replaceInitHooks(
		func(config.MySQLConfig) error {
			calls = append(calls, "mysql")
			return failure
		},
		func() error {
			calls = append(calls, "migrate")
			return nil
		},
		func() error {
			calls = append(calls, "data")
			return nil
		},
		func(config.RedisConfig) error {
			calls = append(calls, "redis")
			return nil
		},
	)
	t.Cleanup(restore)

	if err := Init(&config.Config{}); !errors.Is(err, failure) {
		t.Fatalf("Init() error = %v, want %v", err, failure)
	}
	if len(calls) != 1 || calls[0] != "mysql" {
		t.Fatalf("calls after failure = %#v, want only mysql", calls)
	}
}
