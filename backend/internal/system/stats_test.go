package system

import (
	"sync"
	"testing"
)

// GetStats is called concurrently in production: ws.Hub's poll ticker and
// its per-connection immediate broadcast both call it, unthrottled, so
// overlapping calls are a real occurrence, not a hypothetical. getCPU and
// getNetwork used to mutate package-level state (prevCPUStat, prevNetStats,
// prevNetTime) with no synchronization at all, which `go test -race` here
// caught immediately. Run with -race in CI, not just locally -- a plain
// `go test` won't catch a regression here.
func TestGetStatsConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = GetStats()
		}()
	}
	wg.Wait()
}
