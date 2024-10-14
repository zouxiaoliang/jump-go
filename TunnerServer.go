package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

type TunnelServer struct {
}

func (server *TunnelServer) forwardToTarget(client net.Conn) {
	r := make([]byte, 1)
	r[0] = 1

	l := make([]byte, 4)
	_, err := client.Read(l)
	if err != nil {
		log.Printf("read failed: %v", err)
		client.Write(r)
		client.Close()
		return
	}
	length := binary.BigEndian.Uint32(l)
	b := make([]byte, length)
	_, err = client.Read(b)
	if err != nil {
		log.Printf("read failed: %v", err)
		client.Write(r)
		client.Close()
		return
	}

	targetInfo := Target{}
	d := json.NewDecoder(bytes.NewReader(b))
	err = d.Decode(&targetInfo)
	if err != nil {
		log.Printf("json Decode failed: %v", err)
		client.Write(r)
		client.Close()
		return
	}

	target, err := net.Dial(targetInfo.Scheme, targetInfo.Host)
	if err != nil {
		log.Printf("connect to target {%v} failed: %v", targetInfo, err)
		client.Write(r)
		client.Close()
		return
	}
	r[0] = 0
	client.Write(r)

	log.Printf("connect to target %v success.", targetInfo)

	forwardingTcp(client, target)
}
