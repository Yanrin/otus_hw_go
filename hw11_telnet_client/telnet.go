package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return errors.WithMessage(err, "connection failed")
	}
	c.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %v\n", c.address)

	return nil
}

func (c *Client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return errors.WithMessage(err, "can't send")
	}
	fmt.Fprintln(os.Stderr, "...Sent")
	return nil
}

func (c *Client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return errors.WithMessage(err, "can't receive")
	}
	fmt.Fprintln(os.Stderr, "...Received")
	return nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "close connection error")
	}
	fmt.Fprintln(os.Stderr, "...Connection was closed")
	return nil
}
