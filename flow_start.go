package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type FlowStart struct {
	Header MsgHdr
}

func (m *FlowStart) Encode() []byte {
	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m.Header)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	return b
}

func (m *FlowStart) Decode(msg []byte) error {

	return errors.New(fmt.Sprintf("Not support decode for %s", m.Desc()))
}

func (m *FlowStart) Desc() string {
	return fmt.Sprintf("FLOW_START - id: %d", m.Header.SessId)
}

func (m *FlowStart) RespMsg() []IPDRMsg {

	return nil
}
