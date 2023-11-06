package espresso

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	hasWritten bool
}

func (w *responseWriter) Write(p []byte) (int, error) {
	w.hasWritten = true
	return w.ResponseWriter.Write(p)
}

func (w *responseWriter) WriteHeader(code int) {
	w.hasWritten = true
	w.ResponseWriter.WriteHeader(code)
}
