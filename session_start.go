package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type SessionStart struct {
	Header              MsgHdr
	ExporterBootTime    uint32
	FirstRecordSeqNum   uint64
	DroppedRecordCount  uint64
	Primary             byte
	AckTimeInterval     uint32
	AckSequenceInterval uint32
	DocumentID          []byte
}

func (m *SessionStart) Encode() []byte {
	log.Printf("Not support to Encode!\n")
	return nil
}

func (m *SessionStart) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])

	err := binary.Read(bytesBuffer, endian, &m.Header)

	m.ExporterBootTime = binary.BigEndian.Uint32(msg[8:12])
	m.FirstRecordSeqNum = binary.BigEndian.Uint64(msg[12:20])
	m.DroppedRecordCount = binary.BigEndian.Uint64(msg[20:28])
	m.Primary = msg[28]
	m.AckTimeInterval = binary.BigEndian.Uint32(msg[29:33])
	m.AckSequenceInterval = binary.BigEndian.Uint32(msg[33:37])
	m.DocumentID = make([]byte, 16)
	copy(m.DocumentID, msg[37:])

	return err
}

func (m *SessionStart) Desc() string {
	return fmt.Sprintf("SESSION_START - id: %d", m.Header.SessId)
}

func (m *SessionStart) RespMsg() []IPDRMsg {

	return nil
}
