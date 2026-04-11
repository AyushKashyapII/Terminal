package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	serve := flag.Bool("serve", false, "HTTP server for curl /terminal (plain text). Default: interactive TUI with arrow keys.")
	addr := flag.String("addr", "", "listen address when -serve (e.g. :8080); empty uses $PORT or :8080 (Fly.io sets PORT)")
	flag.Parse()

	if *serve {
		listenAddr := *addr
		if listenAddr == "" {
			if p := os.Getenv("PORT"); p != "" {
				listenAddr = ":" + p
			} else {
				listenAddr = ":8080"
			}
		}
		runServer(listenAddr)
		return
	}

	if err := runTUI(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
