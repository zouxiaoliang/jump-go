package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/zouxiaoliang/jump/std"
)

var (
	config string

	tunnelAddr string

	localAddr  string
	remoteAddr string

	help bool
)

func init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&tunnelAddr, "tunnel", "", "tunnel address, example: 127.0.0.1:1234")
	flag.StringVar(&localAddr, "local", "", "local address, example: 127.0.0.1:1234")
	flag.StringVar(&remoteAddr, "remote", "", "remote address, example: 127.0.0.1:1234")
	flag.StringVar(&config, "config", dirname+"/.config/jump.json", "configure file path.")

	flag.BoolVar(&help, "h", false, "print this message")
}

func parseConfig() *Config {
	if config == "" && (tunnelAddr == "" || localAddr == "" || remoteAddr == "") {
		fmt.Printf("tunnel: %v, local: %v, remote: %v, config: %v\n", tunnelAddr, localAddr, remoteAddr, config)
		return nil
	}

	var c Config
	if config != "" {
		content, err := os.ReadFile(config)
		if err != nil {
			log.Fatalf("Error when opening file, what: %v\n", err)
			return nil
		}

		err = json.Unmarshal(content, &c)
		if err != nil {
			log.Fatalf("Error when unmarshal json, what: %v\n", err)
			return nil
		}
	} else {
		c.Tunnel = tunnelAddr
		c.Forwardings = append(c.Forwardings, Forwardings{localAddr, remoteAddr})
	}

	return &c
}

func forward(tunnel string, local string, remote string) {
	log.Printf("tcp proxy server is running on tcp://%v\n", local)
	// 开始转发流程
	var listener, err = net.Listen("tcp", local)
	if err != nil {
		log.Fatalf("failed to start the proxy: %v\n", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("Close listening failed. what: %v\n", err)
		}
	}(listener)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("Error accept client connection", err)
			break
		}
		log.Printf("new client {%v}", client.RemoteAddr().String())

		var tunnelClient = std.NewTunnelClient(tunnel, remote)
		tunnelClient.ForwardToTunnel(client)
	}
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(1)
	}
	c := parseConfig()
	if c == nil {
		flag.Usage()
		os.Exit(2)
	}

	var wg sync.WaitGroup

	for _, f := range c.Forwardings {
		// Add your forwarding logic here
		wg.Add(1)
		go func() {
			forward(c.Tunnel, f.Local, f.Remote)
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("All forwarding tasks completed")
}
