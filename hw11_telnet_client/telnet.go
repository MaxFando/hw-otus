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

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func (c *telnetClient) Connect() (err error) {
	c.connection, err = net.DialTimeout("tcp", c.address, c.timeout)
	return err
}

func (c *telnetClient) Close() error {
	return c.connection.Close()
}

func (c *telnetClient) Send() error {
	_, err := io.Copy(c.connection, c.in)
	return err
}

func (c *telnetClient) Receive() error {
	_, err := io.Copy(c.out, c.connection)
	return err
}
