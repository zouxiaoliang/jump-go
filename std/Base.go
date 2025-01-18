package std

import "io"

type Target struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
}

func forwardingTcp(src io.ReadWriter, dst io.ReadWriter) {
	go io.Copy(src, dst)
	go io.Copy(dst, src)
}
