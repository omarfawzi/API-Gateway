package errors

import (
	"context"
	"gateway/utils"
	"net/http"
)

const (
	serverHandlerName = "errors-server-handler"
)

type ServerHandler struct {
}

func NewServerHandler() *ServerHandler {
	return &ServerHandler{}
}

type responseCapture struct {
	http.ResponseWriter
	header      http.Header
	body        []byte
	statusCode  int
	wroteHeader bool
}

func (rw *responseCapture) Header() http.Header {
	return rw.header
}

func (rw *responseCapture) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.statusCode = statusCode
		rw.wroteHeader = true
	}
}

func (rw *responseCapture) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return len(b), nil
}

func (r *ServerHandler) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(serverHandlerName, r.middleware)
}

//nolint:gochecknoglobals // Reason: map is constant-like and safe as a global
var errorStatusCodes = map[int]bool{
	http.StatusInternalServerError: true,
	http.StatusMethodNotAllowed:    true,
}

func (r *ServerHandler) middleware(
	_ context.Context,
	_ map[string]interface{},
	next http.Handler,
) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rw := &responseCapture{
			ResponseWriter: w,
			header:         http.Header{},
		}

		next.ServeHTTP(rw, req)

		if isError, ok := errorStatusCodes[rw.statusCode]; ok && isError {
			utils.WriteJSONError(w, rw.statusCode)
			return
		}

		for k, v := range rw.header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(rw.statusCode)
		_, err := w.Write(rw.body)
		if err != nil {
			return
		}

	}), nil
}
