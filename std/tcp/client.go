package std

import (
	"crypto/rc4"
	"encoding/json"
	"io"
	"log"
	"net"

	std "github.com/zouxiaoliang/jump/std"
)

type TcpTunnelClient struct {
	tunnel  string
	target  string
	key     string
	version uint8
}

func NewTcpTunnelClient(tunnel string, target string, key string) *TcpTunnelClient {
	return &TcpTunnelClient{tunnel, target, key, std.V1}
}

func NewTcpTunnelClientV2(tunnel string, target string, key string) *TcpTunnelClient {
	return &TcpTunnelClient{tunnel, target, key, std.V2}
}

func (t *TcpTunnelClient) ForwardToTunnel(incoming net.Conn) {
	var hello std.Hello
	var request std.Request
	var response std.Response
	var e error
	var outgoing io.ReadWriteCloser

	var irc4 *rc4.Cipher
	var orc4 *rc4.Cipher

	// 构建转发协议，填写转发的目标地址
	target := std.Target{Scheme: "tcp", Host: t.target}
	request.Body, e = json.Marshal(target)
	if e != nil {
		log.Println("marshal error:", e)
		goto error_out
	}
	request.Len = uint32(len(request.Body))

	// 连接到通道服务
	outgoing, e = net.Dial("tcp", t.tunnel)
	if e != nil {
		log.Println("dial error:", e)
		goto error_out
	}

	// 发送转发协议版本
	if t.version == std.V2 {
		hello.Type = std.RAW
		if len(t.key) > 0 {
			hello.Type = std.RC4
		}
		e = hello.ToStream(outgoing)
		if e != nil {
			log.Println("failed to send hello, error:", e)
			goto error_out
		}
		if hello.Type == std.RC4 {
			irc4, _ = rc4.NewCipher([]byte(t.key))
			orc4, _ = rc4.NewCipher([]byte(t.key))
			outgoing = NewRC4IO(outgoing, irc4, orc4)
		}
	}

	// 发送跳转协议
	e = request.ToStream(outgoing)
	if e != nil {
		log.Println("failed to send request, error:", e)
		goto error_out
	}

	// 接收跳转协议
	e = response.FromStream(outgoing)
	if e != nil {
		log.Println("failed to recv response, error:", e)
		goto error_out
	}
	if response.Code != 0 {

		log.Printf("connect error, response code: %d", response.Code)
		goto error_out
	}

	log.Printf("connect to target %v success.", target)

	// 开始数据流转发
	forwardingTcp(incoming, outgoing)
	return

error_out:
	incoming.Close()
}
