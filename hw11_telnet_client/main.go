package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

const (
	minArgsCount   = 3
	defaultTimeout = 10
)

func init() {
	flag.DurationVar(&timeout, "timeout", defaultTimeout*time.Second, "connection timeout")
}

func main() {
	flag.Parse()

	if len(os.Args) < minArgsCount {
		log.Fatalf("Expected to have at least 3 arguments, but got %d", len(os.Args))
	}

	host, port := os.Args[2], os.Args[3]
	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if client.Close() != nil {
			log.Fatalln(err)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())
	go worker(client.Receive, cancelFunc)
	go worker(client.Send, cancelFunc)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigCh:
		cancelFunc()
		signal.Stop(sigCh)
		return

	case <-ctx.Done():
		close(sigCh)
		return
	}
}

func worker(handler func() error, cancelFunc context.CancelFunc) {
	if err := handler(); err != nil {
		cancelFunc()
	}
}
