package grace

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type PanicHandler struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real handler
func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			// Dump stacktrace to logs
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong!"))
			// Respond with status
		}
	}()
	p.handler.ServeHTTP(w, req)
}

// NewPanicHandler constructs a new PanicHandler middleware handler
func NewPanicHandler(handlerToWrap http.Handler) *PanicHandler {
	return &PanicHandler{handlerToWrap}
}
