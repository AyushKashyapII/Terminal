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

	if *serve {
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

	if err := runTUI(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
