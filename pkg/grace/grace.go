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
func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			// Dump stacktrace to logs
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			fmt.Fprint(w, "<h1>Something went wrong!</h1>")
			// Respond with status
		}
	}()
	p.handler.ServeHTTP(w, r)
}

// NewPanicHandler constructs a new PanicHandler middleware handler
func NewPanicHandler(handlerToWrap http.Handler) *PanicHandler {
	return &PanicHandler{handlerToWrap}
}
