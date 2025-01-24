package std

import (
	"crypto/rc4"
	"io"
)

type RC4IO struct {
	io io.ReadWriteCloser
	ic *rc4.Cipher
	oc *rc4.Cipher
}

func NewRC4IO(rw io.ReadWriteCloser, read_cipher *rc4.Cipher, write_cipher *rc4.Cipher) *RC4IO {
	return &RC4IO{
		io: rw,
		ic: read_cipher,
		oc: write_cipher,
	}
}

func (rw *RC4IO) Read(p []byte) (n int, err error) {
	n, err = rw.io.Read(p)
	if err != nil {
		return n, err
	}

	if n > 0 {
		rw.ic.XORKeyStream(p[:n], p[:n])
	}

	return n, err
}

func (rw *RC4IO) Write(p []byte) (n int, err error) {
	encrypted := make([]byte, len(p))
	rw.oc.XORKeyStream(encrypted, p)
	return rw.io.Write(encrypted)
}

func (rw *RC4IO) Close() error {
	return rw.io.Close()
}
