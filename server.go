package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func runServer(addr string) {
	mux := http.NewServeMux()
	// Plain paths only — avoids Go 1.22+ "METHOD /path{$}" patterns that panic on some toolchain builds (e.g. Fly builder).
	mux.HandleFunc("/terminal/", handleTerminalSlashRedirect)
	mux.HandleFunc("/terminal", handleTerminalGET)
	mux.HandleFunc("/", handleRootGET)

	log.Printf("listening on %s", addr)
	if runtime.GOOS == "windows" {
		fmt.Println("snapshot (PowerShell: use curl.exe, not curl): curl.exe -s http://127.0.0.1" + addr + "/terminal")
	} else {
		fmt.Println("curl snapshot: curl -s http://127.0.0.1" + addr + "/terminal")
	}
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("http server: %v", err)
	}
}

func requireGET(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func handleTerminalSlashRedirect(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	if r.URL.Path != "/terminal/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/terminal", http.StatusMovedPermanently)
}

func handleTerminalGET(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	handleTerminalPlain(w, r)
}

func handleRootGET(w http.ResponseWriter, r *http.Request) {
	if !requireGET(w, r) {
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	handleRootHint(w, r)
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
