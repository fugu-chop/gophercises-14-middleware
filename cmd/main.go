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
				fmt.Fprintf(w, "<h1>panic: %s</h1><pre>%v</pre>", r, stack)
			}
		}()

		mux.ServeHTTP(w, r)
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
