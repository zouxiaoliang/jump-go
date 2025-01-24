package std

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Target struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
}

const (
	V1 uint8 = 1
	V2 uint8 = 2
)

const (
	RAW uint8 = 0
	RC4 uint8 = 1
	KCP uint8 = 2
)

type Hello struct {
	Type uint8 // RAW,RC4,KCP
}

func (hello *Hello) ToBytes() (buf []byte, err error) {
	b := bytes.Buffer{}

	err = binary.Write(&b, binary.BigEndian, hello.Type)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (hello *Hello) FromBytes(buf []byte) (err error) {
	if len(buf) == 0 {
		return errors.New("buffer not enough")
	}

	hello.Type = uint8(buf[0])
	return nil
}

func (hello *Hello) ToStream(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, hello.Type)
}

func (hello *Hello) FromStream(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, &hello.Type)
}

type Request struct {
	Len  uint32 // body length
	Body []byte // body data
}

func (request *Request) ToBytes() (buf []byte, err error) {
	b := bytes.Buffer{}
	err = binary.Write(&b, binary.BigEndian, request.Len)
	if err != nil {
		return nil, err
	}
	_, err = b.Write(request.Body)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (request *Request) FromBytes(buf []byte) (err error) {
	if len(buf) < 4 {
		return errors.New("buffer not enough")
	}

	request.Len = uint32(binary.BigEndian.Uint32(buf[:4]))
	if int(request.Len) > len(buf)-4 {
		return errors.New("buffer not enough")
	}
	request.Body = buf[4 : 4+request.Len]

	return nil
}

func (request *Request) ToStream(writer io.Writer) (err error) {
	err = binary.Write(writer, binary.BigEndian, request.Len)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.BigEndian, request.Body)
	if err != nil {
		return err
	}
	return nil
}

func (request *Request) FromStream(reader io.Reader) (err error) {
	err = binary.Read(reader, binary.BigEndian, &request.Len)
	if err != nil {
		return err
	}
	request.Body = make([]byte, request.Len)
	err = binary.Read(reader, binary.BigEndian, request.Body)
	if err != nil {
		return err
	}

	return nil
}

type Response struct {
	Code uint8
}

func (response Response) ToBytes() (buf []byte, err error) {
	b := bytes.Buffer{}

	err = binary.Write(&b, binary.BigEndian, response.Code)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (response *Response) FromBytes(buf []byte) (err error) {
	if len(buf) == 0 {
		return errors.New("buffer not enough")
	}

	response.Code = uint8(buf[0])
	return nil
}

func (response *Response) ToStream(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, response.Code)
}

func (response *Response) FromStream(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, &response.Code)
}
