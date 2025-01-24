package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"sync"

	std "github.com/zouxiaoliang/jump/std/tcp"
)

var (
	config  string
	tunnel  string
	version uint
	key     string
	help    bool
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func handleTunnel(c *Config, client net.Conn) {
	var tunnelServer *std.TcpTunnelServer
	if c.Version == 1 {
		tunnelServer = std.NewTcpTunnelServer(c.Key, c.Blacklist, c.Whitelist)
	} else {
		tunnelServer = std.NewTcpTunnelServerV2(c.Key, c.Blacklist, c.Whitelist)
	}
	tunnelServer.ForwardToTarget(client)
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
	if tunnel != "" {
		c.Tunnel = tunnel
	}
	if isFlagPassed("key") {
		c.Key = key
	}

	if isFlagPassed("version") {
		c.Version = version
	}

	if c.Version != 1 && c.Version != 2 {
		log.Fatalf("protocol version error. version: %v\n", c.Version)
		return nil
	}
	if c.Tunnel == "" {
		log.Fatalf("tunnel address error. tunnel: %v\n", c.Tunnel)
		return nil
	}

	return &c
}

func init() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&config, "config", dirname+"/.config/jump.json", "configure file path.")
	flag.StringVar(&tunnel, "tunnel", "127.0.0.1:1234", "tunnel address")
	flag.StringVar(&key, "key", "", "crypt key")
	flag.UintVar(&version, "version", 1, "protocol version")
	flag.BoolVar(&help, "h", false, "print this message")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	c := parseConfig()
	if c == nil {
		flag.Usage()
		os.Exit(2)
	}

	var listener, err = net.Listen("tcp", c.Tunnel)
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
			handleTunnel(c, client)
			wg.Done()
		}()
	}

	wg.Wait()
}
