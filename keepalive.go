package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type KeepAlive struct {
	Header MsgHdr
}

func (m *KeepAlive) Encode() []byte {
	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m.Header)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	return b
}

func (m *KeepAlive) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])
	err := binary.Read(bytesBuffer, endian, &m.Header)

	return err
}

func (m *KeepAlive) Desc() string {
	return fmt.Sprintf("KEEP_ALIVE")
}

func (m *KeepAlive) RespMsg() []IPDRMsg {
	return nil
}

func NewKeepAliveMsg() IPDRMsg {
	var h MsgHdr = MsgHdr{
		Version: 2,
		MsgId:   KEEP_ALIVE,
		SessId:  0,
		MsgFlag: 0,
		MsgLen:  8,
	}

	m := &KeepAlive{
		Header: h,
	}

	return m
}
