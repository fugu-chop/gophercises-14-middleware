package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	wrappedMux := recoverMiddleware(mux, true)

	log.Fatal(http.ListenAndServe(":3000", wrappedMux))
}

// We never need to use Header(), so we don't need to define it
// on our interface
type responseWriterInterceptor struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

// Overwrite the interface methods in the ResponseWriter object we pass
// All of these methods are called with the appropriate arguments by
// the server functionality - we don't need to worry about manual calling
func (rwi *responseWriterInterceptor) Write(b []byte) (int, error) {
	rwi.writes = append(rwi.writes, b)
	return len(b), nil
}

func (rwi *responseWriterInterceptor) WriteHeader(statusCode int) {
	rwi.status = statusCode
}

// Write the values that have been buffered to the original
// ResponseWriter we passed in
func (rwi *responseWriterInterceptor) release() error {
	if rwi.status != 0 {
		rwi.ResponseWriter.WriteHeader(rwi.status)
	}

	for _, write := range rwi.writes {
		if _, err := rwi.ResponseWriter.Write(write); err != nil {
			return err
		}
	}
	return nil
}

func recoverMiddleware(mux http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				log.Println(r)
				stack := string(debug.Stack())
				log.Println(stack)
				if !dev {
					http.Error(w, "Something went wrong!", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %s</h1><pre>%v</pre>", r, stack)
			}
		}()

		// A ResponseWriter is responsible for showing things on the page
		// via Write() and returning an HTTP status via WriteHeader.
		// The idea here is to create a custom type that will 'capture'
		// values for a valid response - e.g. /panic-after
		// but only send them to the browser if no subsequent panic
		// i.e. valid code flow
		interceptor := &responseWriterInterceptor{ResponseWriter: w}
		// This is where the panic will occur when handling request
		// It does not actually serve the response to the request
		mux.ServeHTTP(interceptor, r)
		// This will only run on a successful request - the buffered
		// values are moved to the original ResponseWriter, w
		interceptor.release()
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
