package main

import (
	"fmt"
	"net/http"
	"runtime"
)

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /terminal{$}", handleTerminalPlain)
	mux.HandleFunc("GET /terminal/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/terminal", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /{$}", handleRootHint)

	fmt.Printf("listening on http://127.0.0.1%s\n", addr)
	if runtime.GOOS == "windows" {
		fmt.Println("snapshot (PowerShell: use curl.exe, not curl): curl.exe -s http://127.0.0.1" + addr + "/terminal")
	} else {
		fmt.Println("curl snapshot: curl -s http://127.0.0.1" + addr + "/terminal")
	}
	_ = http.ListenAndServe(addr, mux)
}

func handleRootHint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	host := r.Host
	if host == "" {
		host = "your-domain"
	}
	fmt.Fprintf(w, "GET /terminal for the terminal page (try: curl -sL https://%s/terminal)\n", host)
}

func handleTerminalPlain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, plainTerminalPage())
}

func plainTerminalPage() string {
	return ansiCurlPage()
}
