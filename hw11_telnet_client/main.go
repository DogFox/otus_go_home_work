package main

import (
	"bytes"
	"io"
	"os"
)

func main() {

	in := &bytes.Buffer{}

	client := NewTelnetClient("192.168.1.1", 10, io.NopCloser(in), os.Stdout)
	client.Connect()
}
