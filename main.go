package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)


func main() {
	serve := flag.Bool("serve", false, "HTTP server for curl /terminal (plain text).")
	sshAddr := flag.String("ssh", "", "SSH listen address (e.g. :2222); empty uses $SSH_PORT or :2222")
	addr := flag.String("addr", "", "listen address when -serve (e.g. :8080); empty uses $PORT or :8080")
	flag.Parse()

	runServe := *serve
	if !runServe && len(flag.Args()) == 1 && strings.EqualFold(strings.TrimSpace(flag.Args()[0]), "serve") {
		runServe = true
	}

	// Always run SSH if we are in serve mode, or if explicitly requested
	if runServe || *sshAddr != "" {
		listenAddr := *addr
		if listenAddr == "" {
			if p := os.Getenv("PORT"); p != "" {
				listenAddr = "0.0.0.0:" + strings.TrimSpace(p)
			} else {
				listenAddr = ":8080"
			}
		}

		// Start HTTP server in a goroutine
		if runServe {
			go runServer(listenAddr)
		}

		// Run SSH server (this blocks or we handle signals there)
		runSSHServer(*sshAddr)
		return
	}


	if len(flag.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "unknown arguments:", strings.Join(flag.Args(), " "))
		fmt.Fprintln(os.Stderr, "usage: go run .           # interactive TUI")
		fmt.Fprintln(os.Stderr, "        go run . -serve   # HTTP résumé for curl (or: go run . serve)")
		os.Exit(2)
	}

	if err := runTUI(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
