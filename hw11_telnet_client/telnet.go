package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	Address string
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	Conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

func (client *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", client.Address, client.Timeout)
	if err != nil {
		return err
	}

	client.Conn = conn
	return nil
}

func (client *Client) Close() error {
	err := client.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Send() error {
	_, err := io.Copy(client.Conn, client.In)
	if err != nil {
		return err
	}
	return nil

}

func (client *Client) Receive() error {
	_, err := io.Copy(client.Out, client.Conn)
	if err != nil {
		return err
	}
	return nil
}
