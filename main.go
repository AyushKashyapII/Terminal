package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from your Go server — this is what curl prints.")
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, ASCII_ART)
	})
	fmt.Println("listening on http://127.0.0.1:8080  (try: curl http://127.0.0.1:8080/  and  curl http://127.0.0.1:8080/hello)")
	http.ListenAndServe(":8080", mux)
}