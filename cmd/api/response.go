package api

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
)

type ResponseWriter interface {
	Header() http.Header
	WriteJson(v interface{}) error
	EncodeJson(v interface{}) ([]byte, error)
	WriteHeader(int)
}

type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

var ErrorFieldName = "Error"

func Error(rw ResponseWriter, error string, code int) {
	rw.WriteHeader(code)
	err := rw.WriteJson(map[string]string{ErrorFieldName: error})
	if err != nil {
		panic(err)
	}
}

func NotFound(rw ResponseWriter, r *Request) {
	Error(rw, "Resource not found", http.StatusNotFound)
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.Header().Get("Content-Type") == "" {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) EncodeJson(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (rw *responseWriter) WriteJson(v interface{}) error {
	b, err := rw.EncodeJson(v)
	if err != nil {
		return err
	}
	_, err = rw.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Flush() {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	flusher := rw.ResponseWriter.(http.Flusher)
	flusher.Flush()
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	notifier := rw.ResponseWriter.(http.CloseNotifier)
	return notifier.CloseNotify()
}

// In case if we want don't want to use the built-in server's implementation of the HTTP protocol.
// This might be because you want to switch protocols (to WebSocket for example) or the built-in server is getting in your way.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}
