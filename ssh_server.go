package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)


const (
	defaultSSHPort = "2222"
)

func runSSHServer(addr string) {
	if addr == "" {
		if p := os.Getenv("SSH_PORT"); p != "" {
			addr = ":" + p
		} else {
			addr = ":" + defaultSSHPort
		}
	}

	hostKeyPath := ".ssh/term_ed25519"
	if _, err := os.Stat(hostKeyPath); os.IsNotExist(err) {
		// Ensure directory exists
		_ = os.MkdirAll(".ssh", 0700)
		log.Info("Generating new SSH host key", "path", hostKeyPath)
		// We can just let wish generate one if we don't provide a path, 
		// but providing a path and having it fail if we can't write is safer for persistence.
	}

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(hostKeyPath),

		wish.WithMiddleware(
			bubbletea.Middleware(sshHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "err", err)
		return
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "address", addr)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "err", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "err", err)
	}
}
