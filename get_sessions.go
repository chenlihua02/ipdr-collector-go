package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type GetSessions struct {
	Header    MsgHdr
	RequestId uint16
}

func (m *GetSessions) Encode() []byte {

	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	//slice for msgLen
	msgLen := b[4:8]
	binary.Write(bytesBuffer, endian, uint32(len(b)))
	copy(msgLen, bytesBuffer.Bytes())

	return b
}

func (m *GetSessions) Decode(msg []byte) error {

	return errors.New(fmt.Sprintf("Not support decode for %s", m.Desc()))
}

func (m *GetSessions) Desc() string {
	return "GET_SESSIONS"
}

func (m *GetSessions) RespMsg() []IPDRMsg {

	return nil
}
