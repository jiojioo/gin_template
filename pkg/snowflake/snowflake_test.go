package snowflake_test

import (
	"testing"

	"github.com/jiojioo/gin_template/pkg/snowflake"
)

func TestGenIDProducesDistinctPositiveIDs(t *testing.T) {
	if err := snowflake.Init(1); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	first := snowflake.GenID()
	second := snowflake.GenID()
	if first == 0 || second == 0 {
		t.Fatalf("generated IDs must be positive: %d, %d", first, second)
	}
	if first == second {
		t.Fatalf("generated IDs must be unique: %d", first)
	}
}
