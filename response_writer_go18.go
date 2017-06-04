// +build go1.8

package casper

import (
	"context"
	"net/http"
)

// responseWriter transforms Link Header values into H2 Pushes
type responseWriter struct {
	// http.ResponseWriter is the wrapper http.ResponseWriter
	http.ResponseWriter
	// statusCode is used to temporarily store the http status code
	statusCode int

	ctx context.Context

	targets map[string]*http.PushOptions
}

// reset zeroes out a responseWriter
func (w *responseWriter) reset() *responseWriter {
	w.ResponseWriter = nil
	w.ctx = nil
	w.statusCode = 0
	for k := range w.targets {
		delete(w.targets, k)
	}
	return w
}

// close calls reset and returns a responseWriter to the sync.Pool
func (w *responseWriter) close() {
	w.reset()
	writerPool.Put(w)
}

// Write writes the data to the connection as part of an HTTP reply.
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(200)
	}
	return w.ResponseWriter.Write(b)
}

// WriteHeader will inspect the current response Header and generate H2 Pushes from Link Headers.
// After optionally sending Pushes WriteHeader sends an HTTP response header with status code.
func (w *responseWriter) WriteHeader(s int) {
	// Temporarily store the status code.
	if w.statusCode == 0 {
		w.statusCode = s
	}

	// If the status code is in the 200 range -> generate Pushes.
	if w.statusCode/100 == 2 && len(w.targets) > 0 {
		Push(w)
	}

	// Call WriteHeader on the wrapper http.ResponseWriter
	w.ResponseWriter.WriteHeader(w.statusCode)
}

// Flush sends any buffered data to the client.
func (w *responseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok && flusher != nil {
		flusher.Flush()
	}
}

// CloseNotify returns a channel that receives at most a
// single value (true) when the client connection has gone
// away.
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *responseWriter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := w.ResponseWriter.(http.Pusher)
	if ok && pusher != nil {
		w.targets[target] = opts
		return nil
	}
	return http.ErrNotSupported
}
