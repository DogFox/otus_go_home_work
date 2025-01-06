package main

import (
	"bufio"
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
// go run . --timeout=5s localhost 4242 >./telnet.out
// ncat -l localhost 4242 >./nc.out.
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

	// fmt.Println(connStr, timeout)

	client := NewTelnetClient(connStr, *timeout, io.NopCloser(in), os.Stdout)
	defer client.Close()
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				inputStr, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				in.WriteString(inputStr)
				err = client.Send()
				if err != nil {
					fmt.Println("send error:", err)
					return
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Receive()
				if err != nil {
					fmt.Println("receive error:", err)
					return
				}
			}
		}
	}()

	<-ctx.Done()
	// fmt.Println("Контекст завершён. Программа завершает работу.")
}
