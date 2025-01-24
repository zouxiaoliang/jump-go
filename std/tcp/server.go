package std

import (
	"bytes"
	"crypto/rc4"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"

	std "github.com/zouxiaoliang/jump/std"
)

type TcpTunnelServer struct {
	key          string
	blacklist    sync.Map
	blacklist_on bool
	whitelist    sync.Map
	whitelist_on bool
	version      uint8
}

func NewTcpTunnelServer(key string, blacklist []string, whitelist []string) *TcpTunnelServer {
	server := &TcpTunnelServer{key: key}

	server.blacklist_on = false
	server.whitelist_on = false

	for _, item := range blacklist {
		server.blacklist.Store(item, struct{}{})
		server.blacklist_on = true
	}
	for _, item := range whitelist {
		server.whitelist.Store(item, struct{}{})
		server.whitelist_on = true
	}
	server.version = std.V1
	return server
}

func NewTcpTunnelServerV2(key string, blacklist []string, whitelist []string) *TcpTunnelServer {
	server := &TcpTunnelServer{key: key}

	server.blacklist_on = false
	server.whitelist_on = false

	for _, item := range blacklist {
		server.blacklist.Store(item, struct{}{})
		server.blacklist_on = true
	}
	for _, item := range whitelist {
		server.whitelist.Store(item, struct{}{})
		server.whitelist_on = true
	}
	server.version = std.V2
	return server
}

func (server *TcpTunnelServer) isBlacklisted(target std.Target) bool {
	_, ok := server.blacklist.Load(target.Host)
	return ok
}

func (server *TcpTunnelServer) isWhitelisted(target std.Target) bool {
	_, ok := server.whitelist.Load(target.Host)
	return ok
}

func (server *TcpTunnelServer) ForwardToTarget(incoming io.ReadWriteCloser) {
	var hello std.Hello
	var request std.Request
	var response = &std.Response{Code: 1}
	var outgoing net.Conn
	var targetInfo std.Target
	var d *json.Decoder
	var err error
	var irc4 *rc4.Cipher
	var orc4 *rc4.Cipher

	if server.version == std.V2 {
		err = hello.FromStream(incoming)
		if err != nil {
			log.Printf("read hello failed: %v", err)
			goto error_out
		}

		if hello.Type == std.RC4 {
			irc4, err = rc4.NewCipher([]byte(server.key))
			if err != nil {
				log.Printf("new rc4 cipher failed: %v", err)
				goto error_out
			}
			orc4, err = rc4.NewCipher([]byte(server.key))
			if err != nil {
				log.Printf("new rc4 cipher failed: %v", err)
				goto error_out
			}
			incoming = NewRC4IO(incoming, irc4, orc4)
		}
	}
	err = request.FromStream(incoming)
	if err != nil {
		log.Printf("read request failed: %v", err)
		goto error_out
	}

	d = json.NewDecoder(bytes.NewReader(request.Body))
	err = d.Decode(&targetInfo)
	if err != nil {
		log.Printf("json Decode failed: %v, len: %v, body: %v", err, request.Len, request.Body)
		goto error_out
	}

	if server.blacklist_on && server.isBlacklisted(targetInfo) {
		log.Printf("target %v is in blacklist", targetInfo.Host)
		goto error_out
	}
	if server.whitelist_on && !server.isWhitelisted(targetInfo) {
		log.Printf("target %v is not in whitelist", targetInfo.Host)
		goto error_out
	}

	outgoing, err = net.Dial(targetInfo.Scheme, targetInfo.Host)
	if err != nil {
		log.Printf("connect to target {%v} failed: %v", targetInfo, err)
		goto error_out
	}
	response.Code = 0
	response.ToStream(incoming)

	log.Printf("connect to target %v success.", targetInfo)

	forwardingTcp(incoming, outgoing)
	return

error_out:
	response.ToStream(incoming)
	incoming.Close()
}
