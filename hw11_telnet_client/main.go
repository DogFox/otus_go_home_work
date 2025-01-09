package main

import (
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
// go run . --timeout=5s localhost 4242 >./telnet.out
// ncat -l localhost 4242 >./nc.out.
func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout connection in seconds")
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Println("Please enter host and port. Example --timeout=5s 127.0.0.1 80")
		return
	}

	port := flag.Arg(1)
	host := flag.Arg(0)
	connStr := net.JoinHostPort(host, port)

	client := NewTelnetClient(connStr, *timeout, io.NopCloser(os.Stdin), os.Stdout)
	defer client.Close()
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	go func() {
		err = client.Send()
		if err != nil {
			fmt.Println("send error:", err)
			return
		}
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Println("receive error:", err)
			return
		}
	}()

	<-ctx.Done()
}
