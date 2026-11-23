package app

import (
	"sync"
	"testing"
)

// TestInfoConcurrent stresses SetInfo and GetInfo from many goroutines.
// Without infoMu the Go runtime aborts with "concurrent map iteration and
// map write" or "concurrent map writes"; run with -race for full coverage.
func TestInfoConcurrent(t *testing.T) {
	const key = "test_key"
	t.Cleanup(func() {
		infoMu.Lock()
		delete(Info, key)
		infoMu.Unlock()
	})

	const workers = 32
	const iterations = 500

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				SetInfo(key, j)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				info := GetInfo()
				for range info {
				}
				// the copy is private to the caller
				info["host"] = "caller-only"
			}
		}()
	}
	wg.Wait()

	if _, ok := GetInfo()["host"]; ok {
		t.Fatal("writes to the GetInfo copy leaked into Info")
	}
}
