package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type TemplateData struct {
	Header   MsgHdr
	ConfigID uint16
	Flags    uint8
}

func (m *TemplateData) Encode() []byte {
	log.Printf("Not support to Encode!\n")
	return nil
}

func (m *TemplateData) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])

	err := binary.Read(bytesBuffer, endian, &m.Header)
	m.ConfigID = binary.BigEndian.Uint16(msg[8:10])
	m.Flags = msg[10]

	return err
}

func (m *TemplateData) Desc() string {
	return fmt.Sprintf("TEMPLATE_DATA - id: %d", m.Header.SessId)
}

func (m *TemplateData) RespMsg() []IPDRMsg {
	msgs := []IPDRMsg{}

	var h MsgHdr = MsgHdr{
		Version: 2,
		MsgId:   FINAL_TEMPLATE_DATA_ACK,
		SessId:  m.Header.SessId,
		MsgFlag: 0,
		MsgLen:  8,
	}

	templateAck := &TemplateDataAck{
		Header: h,
	}

	msgs = append(msgs, templateAck)

	return msgs
}
