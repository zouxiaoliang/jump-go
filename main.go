package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	role   string
	local  string
	tunnel string
	target string
)

type Target struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
}

func forwardingTcp(src io.ReadWriter, dst io.ReadWriter) {
	go io.Copy(src, dst)
	go io.Copy(dst, src)
}

func handleClient(client net.Conn) {
	var tunnelClient = NewTunnelClient(tunnel, target)
	tunnelClient.forwardToTunnel(client)
}

func handleTunnel(client net.Conn) {
	var tunnelServer = TunnelServer{}
	tunnelServer.forwardToTarget(client)
}

func init() {
	flag.StringVar(&role, "role", "local", "choise in {local, tunnel}, default is local")
	flag.StringVar(&local, "local", "127.0.0.1:1234", "local listen address")
	flag.StringVar(&tunnel, "tunnel", "127.0.0.1:1234", "tunnel address")
	flag.StringVar(&target, "target", "127.0.0.1:2234", "destination address")
}

func main() {
	flag.Parse()
	var addr = local
	if role == "local" {
		addr = local
	} else if role == "tunnel" {
		addr = tunnel
	} else {
		log.Fatal("role must be in {local, tunnel}")
	}
	log.Printf("tcp proxy server is running on tcp://%v, role: %v", local, role)

	var listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to start the %v proxy: %v", role, err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal("Close listening failed. what: ", err)
		}
	}(listener)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("Error accept client connection", err)
			continue
		}
		log.Printf("new client {%v}", client.RemoteAddr().String())
		if role == "local" {
			handleClient(client)
		} else {
			handleTunnel(client)
		}
	}
}
