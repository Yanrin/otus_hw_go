package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type FlagCL struct { // flags from command line
	timeout time.Duration
	help    bool
}

var (
	DefaultTimeout = 10 * time.Second
	host, port     string
	fcl            FlagCL
)

func init() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage: %v [-options] host port\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVar(&fcl.help, "help", false, "output a usage message and exit")

	flag.DurationVar(&fcl.timeout, "timeout", DefaultTimeout, "timeout server connection")
}

func main() {
	flag.Parse()
	if fcl.help {
		flag.Usage()
		return
	}

	if len(flag.Args()) < 2 {
		log.Fatal("host and port should be defined")
	}
	host = flag.Arg(0)
	port = flag.Arg(1)

	tcpaddress := net.JoinHostPort(host, port)
	tc := NewTelnetClient(tcpaddress, fcl.timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer tc.Close()

	ctxC, stopC := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		defer stopC()
		if err := tc.Send(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	go func() {
		defer stopC()
		if err := tc.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	<-ctxC.Done()
}
