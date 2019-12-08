package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Connect struct {
	Header       MsgHdr
	InitAddr     uint32
	InitPort     uint16
	Capabilities uint32
	KaInterval   uint32
	VendorId     UTF8String
}

func (m *Connect) Encode() []byte {

	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m.Header)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	binary.Write(bytesBuffer, endian, m.InitAddr)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	binary.Write(bytesBuffer, endian, m.InitPort)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	binary.Write(bytesBuffer, endian, m.Capabilities)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	binary.Write(bytesBuffer, endian, m.KaInterval)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	b = append(b, m.VendorId.Encode()...)

	//slice for msgLen
	msgLen := b[4:8]
	binary.Write(bytesBuffer, endian, uint32(len(b)))
	copy(msgLen, bytesBuffer.Bytes())

	return b
}

func (m *Connect) Decode(msg []byte) error {

	return errors.New(fmt.Sprintf("Not support decode for %s", m.Desc()))
}

func (m *Connect) Desc() string {
	return "CONNECT"
}

func (m *Connect) RespMsg() []IPDRMsg {

	return nil
}
