package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type FlowStop struct {
	Header     MsgHdr
	ReasonCode uint16
	ReasonInfo UTF8String
}

func (m *FlowStop) Encode() []byte {
	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m.Header)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	binary.Write(bytesBuffer, endian, m.ReasonCode)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	b = append(b, m.ReasonInfo.Encode()...)

	//slice for msgLen
	msgLen := b[4:8]
	binary.Write(bytesBuffer, endian, uint32(len(b)))
	copy(msgLen, bytesBuffer.Bytes())

	return b
}

func (m *FlowStop) Decode(msg []byte) error {

	return errors.New(fmt.Sprintf("Not support decode for %s", m.Desc()))
}

func (m *FlowStop) Desc() string {
	return fmt.Sprintf("FLOW_STOP - id: %d", m.Header.SessId)
}

func (m *FlowStop) RespMsg() []IPDRMsg {

	return nil
}
