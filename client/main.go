package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	std "github.com/zouxiaoliang/jump/std/tcp"
)

var (
	config     string
	tunnelAddr string
	localAddr  string
	remoteAddr string
	key        string
	v          uint
	help       bool
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
	flag.StringVar(&key, "key", "", "crypt key")
	flag.UintVar(&v, "version", 1, "protocol version")

	flag.BoolVar(&help, "h", false, "print this message")
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func parseConfig() *Config {
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
	}

	if !isFlagPassed("tunnel") {
		tunnelAddr = c.Tunnel
	}
	if isFlagPassed("local") && isFlagPassed("remote") {
		c.Forwardings = []Forwarding{{Local: localAddr, Remote: remoteAddr}}
	}

	if !isFlagPassed("version") {
		v = c.Version
	}

	if !isFlagPassed("key") {
		key = c.Key
	}

	if tunnelAddr == "" {
		fmt.Printf("tunnel address error. tunnel: %v\n", tunnelAddr)
		return nil
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

		var tunnelClient *std.TcpTunnelClient
		if v == 1 {
			tunnelClient = std.NewTcpTunnelClient(tunnel, remote, key)
		} else {
			tunnelClient = std.NewTcpTunnelClientV2(tunnel, remote, key)
		}
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
