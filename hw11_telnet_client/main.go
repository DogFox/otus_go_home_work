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
	"sync"
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

	fmt.Println(connStr, timeout)

	client := NewTelnetClient(connStr, *timeout, io.NopCloser(in), os.Stdout)
	defer client.Close()
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	wg := &sync.WaitGroup{}

	go func() {
		for {
			select {
			case <-ctx.Done():
				// fmt.Println("Получен сигнал, завершаем горутину.")
				return
			default:
				wg.Add(2)

				go sendWorker(wg, client, in)
				go receiveWorker(wg, client)

				wg.Wait()
			}
		}
	}()

	<-ctx.Done()
	// fmt.Println("Контекст завершён. Программа завершает работу.")
}

func sendWorker(wg *sync.WaitGroup, client TelnetClient, in *bytes.Buffer) {
	defer wg.Done()
	reader := bufio.NewReader(os.Stdin)
	for {
		inputStr, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		in.WriteString(inputStr)
		err = client.Send()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func receiveWorker(wg *sync.WaitGroup, client TelnetClient) {
	defer wg.Done()
	for {
		err := client.Receive()
		if err != nil {
			fmt.Println("receive err: ", err)
			return
		}
	}
}
