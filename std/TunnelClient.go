package std

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

type TunnelClient struct {
	tunnel string
	target string
}

func NewTunnelClient(tunnel string, target string) *TunnelClient {
	return &TunnelClient{tunnel, target}
}

func (t *TunnelClient) ForwardToTunnel(client net.Conn) {
	target := Target{"tcp", t.target}
	s, e := json.Marshal(target)
	if e != nil {
		log.Println("marshal error:", e)
		client.Close()
		return
	}

	proxy, e := net.Dial("tcp", t.tunnel)
	if e != nil {
		log.Println("dial error:", e)
		client.Close()
		return
	}
	var l = make([]byte, 4)
	binary.BigEndian.PutUint32(l, uint32(len(s)))
	proxy.Write(l)
	proxy.Write([]byte(s))

	var r = make([]byte, 1)
	proxy.Read(r)

	if r[0] != 0 {
		client.Close()
		log.Printf("connect error, code: %d", r[0])
		return
	}

	log.Printf("connect to target %v success.", target)

	forwardingTcp(client, proxy)
}
