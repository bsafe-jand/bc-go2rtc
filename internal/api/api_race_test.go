package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/AlexxIT/go2rtc/internal/app"
)

// TestAPIHandlerConcurrent reproduces issue #2443: overlapping GET /api
// requests, plus a module writing app.Info, used to crash the process with
// "concurrent map iteration and map write". Each response must also carry
// its own request's Host, not another request's.
func TestAPIHandlerConcurrent(t *testing.T) {
	// The writers re-set an existing key to its current value, so the test
	// leaves app.Info as it found it.
	const key = "version"
	value := app.GetInfo()[key]

	const workers = 32
	const iterations = 200

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			host := fmt.Sprintf("host-%d", i)
			for j := 0; j < iterations; j++ {
				r := httptest.NewRequest("GET", "/api", nil)
				r.Host = host
				w := httptest.NewRecorder()
				apiHandler(w, r)

				if w.Code != http.StatusOK {
					t.Errorf("status %d", w.Code)
					return
				}
				var info map[string]any
				if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil {
					t.Errorf("decode: %v", err)
					return
				}
				if info["host"] != host {
					t.Errorf("host %v, want %s", info["host"], host)
					return
				}
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				app.SetInfo(key, value)
			}
		}()
	}
	wg.Wait()
}
