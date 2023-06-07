package espresso

import (
	"bufio"
	"net"
	"net/http"
)

type responseWriter[Data any] struct {
	http.ResponseWriter
	ctx *brewContext[Data]
}

func (w *responseWriter[Data]) write() {
	w.ctx.hasWroteResponseCode = true
}

func (w *responseWriter[Data]) WriteHeader(code int) {
	w.write()
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter[Data]) Write(b []byte) (int, error) {
	w.write()
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter[Data]) Flush() {
	flush, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	w.write()
	flush.Flush()
}

func (w *responseWriter[Data]) FlushError() error {
	flush, ok := w.ResponseWriter.(interface{ FlushError() error })
	if !ok {
		return nil
	}

	w.write()
	return flush.FlushError()
}

func (w *responseWriter[Data]) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijack, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, nil
	}

	return hijack.Hijack()
}
