package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultDuration = 10 * time.Second

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", defaultDuration, "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "Usage: telnet [--timeout=<duration>] <host> <port>")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "...Failed to connect:", err)
		return
	}
	defer client.Close()
	fmt.Fprintln(os.Stderr, "...Connected to ", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "...Send error:", err)
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		}

		fmt.Fprintln(os.Stderr, "...EOF")
		cancel()
	}()

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "...Receive error:", err)
		}

		fmt.Fprintln(os.Stderr, "...Exit receive")
		cancel()
	}()

	<-ctx.Done()
}
