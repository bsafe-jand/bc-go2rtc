package streams

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/AlexxIT/go2rtc/pkg/core"
)

// TestApiStreamsConcurrent stresses the streams map through the HTTP
// handlers: PUT and DELETE write it while list, single GET and .dot read it.
// Before the fix the list, .dot and DELETE paths touched the map without
// streamsMu and the runtime aborted with a concurrent map error.
func TestApiStreamsConcurrent(t *testing.T) {
	HandleFunc("race", func(url string) (core.Producer, error) { return nil, nil })

	const workers = 16
	const iterations = 200

	names := make([]string, workers)
	for i := range names {
		names[i] = fmt.Sprintf("race_%d", i)
	}
	t.Cleanup(func() {
		for _, name := range names {
			Delete(name)
		}
	})

	var wg sync.WaitGroup
	for _, name := range names {
		put := "/api/streams?" + url.Values{"name": {name}, "src": {"race:" + name}}.Encode()
		one := "/api/streams?" + url.Values{"src": {name}}.Encode()
		dot := "/api/streams.dot?" + url.Values{"src": {name}}.Encode()

		wg.Add(5)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				apiStreams(httptest.NewRecorder(), httptest.NewRequest("PUT", put, nil))
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				apiStreams(httptest.NewRecorder(), httptest.NewRequest("DELETE", one, nil))
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				w := httptest.NewRecorder()
				apiStreams(w, httptest.NewRequest("GET", "/api/streams", nil))
				if w.Code != 200 {
					t.Errorf("list status %d", w.Code)
					return
				}
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				apiStreams(httptest.NewRecorder(), httptest.NewRequest("GET", one, nil))
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				apiStreamsDOT(httptest.NewRecorder(), httptest.NewRequest("GET", dot, nil))
				apiStreamsDOT(httptest.NewRecorder(), httptest.NewRequest("GET", "/api/streams.dot", nil))
			}
		}()
	}
	wg.Wait()
}
