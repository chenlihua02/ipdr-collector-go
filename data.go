package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type Data struct {
	Header      MsgHdr
	TemplateID  uint16
	ConfigID    uint16
	Flags       byte
	SequenceNum uint64
	Record      []byte
}

func (m *Data) Encode() []byte {
	log.Printf("Not support to Encode!\n")
	return nil
}

func (m *Data) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])

	err := binary.Read(bytesBuffer, endian, &m.Header)
	m.TemplateID = binary.BigEndian.Uint16(msg[8:10])
	m.ConfigID = binary.BigEndian.Uint16(msg[10:12])
	m.Flags = msg[12]
	m.SequenceNum = binary.BigEndian.Uint64(msg[13:21])
	//data_len = msg_len - header_len - 17
	m.Record = make([]byte, m.Header.MsgLen-8-17+1)

	copy(m.Record, msg[21:])

	//fmt.Printf("Got seq %d\n", lastSeqNum)
	return err
}

func (m *Data) Desc() string {
	return fmt.Sprintf("DATA - id: %d, seq: %d",
		m.Header.SessId, m.SequenceNum)
}

var i int

func (m *Data) RespMsg() []IPDRMsg {
	return nil
}
