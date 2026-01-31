package middleware

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
)

type ContextKey string

const contextKeyRequestID ContextKey = "requestID"

type ResponseLogger struct {
	w          http.ResponseWriter
	statusCode int
	body       []string
}

func (r *ResponseLogger) Header() http.Header {
	return r.w.Header()
}

func (r *ResponseLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.w.(http.Hijacker).Hijack()
}

func (r *ResponseLogger) Flush() {
	r.w.(http.Flusher).Flush()
}

func (r *ResponseLogger) Write(res []byte) (int, error) {
	r.body = append(r.body, string(res))
	return r.w.Write(res)
}

func (r *ResponseLogger) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.w.WriteHeader(statusCode)
}

func (r *ResponseLogger) fullResponse() string {
	if len(r.body) == 1 {
		return r.body[0]
	}
	var res string
	for i, b := range r.body {
		res += fmt.Sprintf("written %d: %s\n", i+1, b)
	}
	return res
}

func (m *middleware) RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New()
		ctx := context.WithValue(r.Context(), contextKeyRequestID, reqID.String())
		r = r.WithContext(ctx)

		rl := &ResponseLogger{w: w, body: make([]string, 0)}

		m.lg.Info(reqID.String(), "Request", fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr))

		h.ServeHTTP(rl, r)
	})
}
