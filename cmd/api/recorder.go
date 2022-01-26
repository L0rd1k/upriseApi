package api

import (
	"bufio"
	"net"
	"net/http"
)

type RecorderMiddleware struct{}
type recorderResponseWriter struct {
	ResponseWriter
	statusCode   int
	wroteHeader  bool
	bytesWritten int64
}

func (rm *RecorderMiddleware) MiddlewareFunc(handler HandlerFunc) HandlerFunc {
	return func(w ResponseWriter, r *Request) {
		writer := &recorderResponseWriter{w, 0, false, 0}
		handler(writer, r)
		r.environment["STATUS_CODE"] = writer.statusCode
		r.environment["BYTES_WRITTEN"] = writer.bytesWritten
	}
}

func (w *recorderResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.wroteHeader {
		return
	}
	w.statusCode = code
	w.wroteHeader = true
}

func (w *recorderResponseWriter) WriteJson(v interface{}) error {
	b, err := w.EncodeJson(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (w *recorderResponseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	flusher := w.ResponseWriter.(http.Flusher)
	flusher.Flush()
}

func (w *recorderResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}

func (w *recorderResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	writer := w.ResponseWriter.(http.ResponseWriter)
	written, err := writer.Write(b)
	w.bytesWritten += int64(written)
	return written, err
}
