package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

var ErrorFieldName = "Error"

//> public - Custom ResponseWriter
type ResponseWriter interface {
	Header() http.Header                      //> Method of http.ResponseWriter
	WriteHeader(int)                          //> Method of http.ResponseWriter
	WriteJson(v interface{}) error            //> Generate the payload
	EncodeJson(v interface{}) ([]byte, error) //> Enocode data to JSON
}

//> private
type responseWriter struct {
	http.ResponseWriter      //> Sending response to any connected HTTP clients.
	wroteHeader         bool //> Check if header is fill
}

func (rw *responseWriter) WriteHeader(code int) {
	fmt.Println("<Response:WriteHeader> code", code)
	//> Check the Extensions(MIME) of the Project - set if not defined
	if rw.Header().Get("Content-Type") == "" {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	rw.ResponseWriter.WriteHeader(code) //> Send an HTTP response with given code.
	rw.wroteHeader = true               //> Set flag
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

func (rw *responseWriter) Write(b []byte) (int, error) {
	fmt.Println("<Response:Write> : ", rw.wroteHeader)
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	fmt.Println("Byte array to string:", string(b))
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Flush() {
	fmt.Println("<Response:Flush>")
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	flusher := rw.ResponseWriter.(http.Flusher)
	flusher.Flush()
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	fmt.Println("<Response:CloseNotify>")
	notifier := rw.ResponseWriter.(http.CloseNotifier)
	return notifier.CloseNotify()
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	fmt.Println("<Response:Hijack>")
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}
