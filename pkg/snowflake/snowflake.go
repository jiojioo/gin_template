// Package snowflake generates process-local unique identifiers.
package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	mu   sync.RWMutex
	node *snowflake.Node
)

func Init(nodeID int64) error {
	next, err := snowflake.NewNode(nodeID)
	if err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	node = next
	return nil
}

func GenID() uint64 {
	mu.RLock()
	current := node
	mu.RUnlock()
	if current == nil {
		panic("snowflake is not initialized")
	}
	return uint64(current.Generate().Int64())
}
