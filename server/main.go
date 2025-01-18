package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"

	"github.com/zouxiaoliang/jump/std"
)

var (
	tunnelAddr string
	help       bool
)

func handleTunnel(client net.Conn) {
	var tunnelServer = std.TunnelServer{}
	tunnelServer.ForwardToTarget(client)
}

func init() {
	flag.StringVar(&tunnelAddr, "tunnel", "127.0.0.1:1234", "tunnel address")
	flag.BoolVar(&help, "h", false, "print this message")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	var listener, err = net.Listen("tcp", tunnelAddr)
	if err != nil {
		log.Fatalf("failed to start the proxy: %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal("Close listening failed. what: ", err)
		}
	}(listener)

	var wg sync.WaitGroup
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("Error accept client connection", err)
			break
		}
		log.Printf("new client {%v}", client.RemoteAddr().String())
		wg.Add(1)
		go func() {
			handleTunnel(client)
			wg.Done()
		}()
	}

	wg.Wait()
}
