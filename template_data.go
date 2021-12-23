package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type FieldDescriptor struct {
	TypeID    uint32
	FieldID   uint32
	FieldName UTF8String
	IsEnabled byte
}

type TemplateBlock struct {
	TemplateID uint16
	SchemaName UTF8String
	TypeName   UTF8String
	Fields     []FieldDescriptor
}

type TemplateData struct {
	Header    MsgHdr
	ConfigID  uint16
	Flags     uint8
	Templates []TemplateBlock
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

	numTemplates := binary.BigEndian.Uint32(msg[11:15])
	msg = msg[15:]
	var msgLen uint32
	//log.Printf("num temps: %d\n", numTemplates)
	for i := uint32(0); i < numTemplates; i++ {
		tb := TemplateBlock{}
		tb.TemplateID = binary.BigEndian.Uint16(msg[:2])
		msg = msg[2:]
		tb.SchemaName, msgLen = DecodeUTF8String(msg)
		msg = msg[msgLen:]
		tb.TypeName, msgLen = DecodeUTF8String(msg)
		msg = msg[msgLen:]
		numFields := binary.BigEndian.Uint32(msg[:4])
		msg = msg[4:]
		for j := uint32(0); j < numFields; j++ {
			f := FieldDescriptor{}
			f.TypeID = binary.BigEndian.Uint32(msg[:4])
			f.FieldID = binary.BigEndian.Uint32(msg[4:8])
			msg = msg[8:]
			f.FieldName, msgLen = DecodeUTF8String(msg)
			msg = msg[msgLen:]
			f.IsEnabled = msg[0]
			msg = msg[1:]
			tb.Fields = append(tb.Fields, f)
		}
		m.Templates = append(m.Templates, tb)
	}

	//log.Printf("Decode Template data: % +v\n", m)

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