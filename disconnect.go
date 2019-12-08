package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Disconnect struct {
	Header MsgHdr
}

func (m *Disconnect) Encode() []byte {
	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m.Header)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	return b
}

func (m *Disconnect) Decode(msg []byte) error {

	return errors.New(fmt.Sprintf("Not support decode for %s", m.Desc()))
}

func (m *Disconnect) Desc() string {
	return "DISCONNECT"
}

func (m *Disconnect) RespMsg() []IPDRMsg {

	return nil
}
