package std

import "io"

func forwardingTcp(src io.ReadWriter, dst io.ReadWriter) {
	go io.Copy(src, dst)
	go io.Copy(dst, src)
}
