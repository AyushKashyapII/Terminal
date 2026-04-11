package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	serve := flag.Bool("serve", false, "HTTP server for curl /terminal (plain text). Default: interactive TUI with arrow keys.")
	addr := flag.String("addr", "", "listen address when -serve (e.g. :8080); empty uses $PORT or :8080 (Fly.io sets PORT)")
	flag.Parse()

	runServe := *serve
	if !runServe {
		args := flag.Args()
		if len(args) == 1 && strings.EqualFold(strings.TrimSpace(args[0]), "serve") {
			runServe = true
		}
	}

	if runServe {
		listenAddr := *addr
		if listenAddr == "" {
			if p := os.Getenv("PORT"); p != "" {
				// Fly’s proxy connects over IPv4; bind all IPv4 interfaces explicitly.
				listenAddr = "0.0.0.0:" + strings.TrimSpace(p)
			} else {
				listenAddr = ":8080"
			}
		}
		if _, err := net.ResolveTCPAddr("tcp", listenAddr); err != nil {
			log.Fatalf("bad -serve listen address %q: %v", listenAddr, err)
		}
		runServer(listenAddr)
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
