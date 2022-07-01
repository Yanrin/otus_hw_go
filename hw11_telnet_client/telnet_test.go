package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("stderr check", func(t *testing.T) {
		origStderr := os.Stderr
		defer func() { os.Stderr = origStderr }()

		host := "127.0.0.1"
		port := "4242"
		address := host + ":" + port

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		listener, err := net.Listen("tcp", address)
		require.NoError(t, err)
		defer func() { require.NoError(t, listener.Close()) }()

		rc, wc, _ := os.Pipe()
		os.Stderr = wc

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			timeout, err := time.ParseDuration("1s")
			require.NoError(t, err)

			client := NewTelnetClient(address, timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			require.True(t, checkFStd(rc, fmt.Sprintf("...Connected to %v", address)))

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)
			require.True(t, checkFStd(rc, "...Sent"))

			err = client.Receive()
			require.NoError(t, err)
			require.True(t, checkFStd(rc, "...Received"))

			require.NoError(t, client.Close())
			require.True(t, checkFStd(rc, "...Connection was closed"))
		}()

		go func() {
			defer wg.Done()

			conn, err := listener.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestConnect(t *testing.T) {
	t.Run("refused", func(t *testing.T) {
		host := "127.0.0.1"
		port := "111"
		address := host + ":" + port

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(address, DefaultTimeout, ioutil.NopCloser(in), out)
		err := client.Connect()

		msg := fmt.Sprintf("dial tcp %v: connect: connection refused", address)
		require.Errorf(t, errors.Unwrap(err), msg)
	})

	t.Run("opened", func(t *testing.T) {
		host := "otus.ru"
		port := "443"
		address := host + ":" + port
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(address, DefaultTimeout, ioutil.NopCloser(in), out)

		require.NoError(t, client.Connect())
	})

	t.Run("timeout", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("1ns")
		require.NoError(t, err)

		client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
		err = client.Connect()

		require.Errorf(t, errors.Unwrap(err), "i/o timeout")
	})
}

// checkFStd checks if fake stderr contains str.
func checkFStd(r io.ReadCloser, str string) bool {
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	return strings.Contains(string(buf[:n]), str)
}
