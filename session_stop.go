package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type SessionStop struct {
	Header     MsgHdr
	ReasonCode uint16
	ReasonInfo UTF8String
}

func (m *SessionStop) Encode() []byte {
	log.Printf("Not support to Encode!\n")
	return nil
}

func (m *SessionStop) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])

	err := binary.Read(bytesBuffer, endian, &m.Header)
	m.ReasonCode = binary.BigEndian.Uint16(msg[8:10])
	m.ReasonInfo.Length = binary.BigEndian.Uint32(msg[10:14])
	m.ReasonInfo.Str = make([]byte, m.ReasonInfo.Length+1)
	copy(m.ReasonInfo.Str, msg[14:])

	return err
}

func (m *SessionStop) Desc() string {
	return fmt.Sprintf("SESSION_STOP - id: %d", m.Header.SessId)
}

func (m *SessionStop) RespMsg() []IPDRMsg {

	msgs := []IPDRMsg{}

	var h MsgHdr = MsgHdr{
		Version: 2,
		MsgId:   FLOW_STOP,
		SessId:  m.Header.SessId,
		MsgFlag: 0,
	}

	var str = "Exporter Handler Shutdown"
	var data []byte = []byte(str)
	var reasonInfo UTF8String = UTF8String{
		Length: uint32(len(str)),
		Str:    data,
	}

	flowStop := &FlowStop{
		Header:     h,
		ReasonCode: 0,
		ReasonInfo: reasonInfo,
	}

	msgs = append(msgs, flowStop)

	return msgs
}
