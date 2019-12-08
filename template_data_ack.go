package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type TemplateDataAck struct {
	Header MsgHdr
}

func (m *TemplateDataAck) Encode() []byte {
	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m.Header)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	return b
}

func (m *TemplateDataAck) Decode(msg []byte) error {

	return errors.New(fmt.Sprintf("Not support decode for %s", m.Desc()))
}

func (m *TemplateDataAck) Desc() string {
	return fmt.Sprintf("FINAL_TEMPLATE_DATA_ACK - id: %d", m.Header.SessId)
}

func (m *TemplateDataAck) RespMsg() []IPDRMsg {

	return nil
}
