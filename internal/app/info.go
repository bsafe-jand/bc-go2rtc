package app

import "sync"

// infoMu guards Info. Modules write Info during Init (for example
// rtsp.Init) after the API listener may already be serving GET /api, so
// every reader and writer must take the lock.
// https://github.com/AlexxIT/go2rtc/issues/2443
var infoMu sync.RWMutex

// SetInfo sets one key in Info. It is safe for concurrent use.
func SetInfo(key string, value any) {
	infoMu.Lock()
	Info[key] = value
	infoMu.Unlock()
}

// GetInfo returns a shallow copy of Info. It is safe for concurrent use,
// and the caller may modify the returned map.
func GetInfo() map[string]any {
	infoMu.RLock()
	info := make(map[string]any, len(Info)+1)
	for k, v := range Info {
		info[k] = v
	}
	infoMu.RUnlock()
	return info
}
