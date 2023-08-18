package espresso

import (
	"net/http"

	"golang.org/x/exp/slog"
)

type responseWriter struct {
	http.ResponseWriter
	logger *slog.Logger

	wroteHeader  bool
	responseCode int
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		w.logger.Error("Already wrote header")
		return
	}

	w.wroteHeader = true
	w.responseCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	w.ensureWriteHeader()

	return w.ResponseWriter.Write(p)
}

func (w *responseWriter) ensureWriteHeader() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
}
