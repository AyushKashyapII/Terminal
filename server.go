package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"syscall"
)

func runServer(preferredAddr string) {
	mux := http.NewServeMux()
	// Plain paths only — avoids Go 1.22+ "METHOD /path{$}" patterns that panic on some toolchain builds (e.g. Fly builder).
	mux.HandleFunc("/terminal/", handleTerminalSlashRedirect)
	mux.HandleFunc("/terminal", handleTerminalGET)
	mux.HandleFunc("/", handleRootGET)

	ln, boundAddr, err := listenTCP(preferredAddr)
	if err != nil {
		log.Fatalf("http server: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	curlBase := fmt.Sprintf("http://127.0.0.1:%d", port)

	log.Printf("listening on %s", boundAddr)
	if boundAddr != preferredAddr {
		log.Printf("(%s was busy; using %s instead)", preferredAddr, boundAddr)
	}
	if runtime.GOOS == "windows" {
		fmt.Println("resume (PowerShell: use curl.exe): curl.exe -s " + curlBase + "/")
		fmt.Println("same at /terminal: curl.exe -s " + curlBase + "/terminal")
	} else {
		fmt.Println("resume: curl -s " + curlBase + "/")
		fmt.Println("same at /terminal: curl -s " + curlBase + "/terminal")
	}

	srv := &http.Server{Handler: mux}
	if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server: %v", err)
	}
}

// listenTCP binds preferredAddr, or if that is ":8080" and in use, tries :8081 … :8099.
func listenTCP(preferredAddr string) (ln net.Listener, actualAddr string, err error) {
	tryAddrs := []string{preferredAddr}
	if preferredAddr == ":8080" {
		for p := 8081; p <= 8099; p++ {
			tryAddrs = append(tryAddrs, fmt.Sprintf(":%d", p))
		}
	}

	var firstErr error
	for _, a := range tryAddrs {
		ln, err = net.Listen("tcp", a)
		if err == nil {
			return ln, a, nil
		}
		if firstErr == nil {
			firstErr = err
		}
		if preferredAddr != ":8080" || !isAddrInUse(err) {
			break
		}
	}
	return nil, "", fmt.Errorf("%w (if you need a specific port, use e.g. -addr :3000)", firstErr)
}

func isAddrInUse(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, syscall.EADDRINUSE) {
		return true
	}
	var op *net.OpError
	if errors.As(err, &op) && op.Err != nil {
		if errors.Is(op.Err, syscall.EADDRINUSE) {
			return true
		}
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "only one usage of each socket address") ||
		strings.Contains(msg, "address already in use")
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, plainTerminalPage())
}

func handleTerminalPlain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, plainTerminalPage())
}

func plainTerminalPage() string {
	return ansiCurlPage()
}
