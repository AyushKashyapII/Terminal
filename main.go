package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	serve := flag.Bool("serve", false, "HTTP server for curl /terminal (plain text). Default: interactive TUI with arrow keys.")
	addr := flag.String("addr", ":8080", "listen address when -serve")
	flag.Parse()

	if *serve {
		runServer(*addr)
		return
	}

	if err := runTUI(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
