package PanicHandler

import "net/http"

type PanicHandler struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real handler
func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Do stuff
	p.handler.ServeHTTP(w, r)
	// Do stuff
}

// NewPanicHandler constructs a new PanicHandler middleware handler
func NewPanicHandler(handlerToWrap http.Handler) *PanicHandler {
	return &PanicHandler{handlerToWrap}
}
