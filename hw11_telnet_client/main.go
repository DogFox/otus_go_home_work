package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// go run . go-telnet --timeout=5s host port
// go run . --timeout=5s host port
func main() {

	timeout := flag.Duration("timeout", 10*time.Second, "Timeout connection in seconds")
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Println("Please enter host and port. Example --timeout=5s 127.0.0.1 80")
		return
	}

	in := &bytes.Buffer{}
	port := flag.Arg(1)
	host := flag.Arg(0)
	connStr := net.JoinHostPort(host, port)

	fmt.Println(connStr, timeout)

	client := NewTelnetClient(connStr, *timeout, io.NopCloser(in), os.Stdout)
	// defer client.Close()
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()
}
