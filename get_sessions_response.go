package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

type SessionBlock struct {
	SessId              byte
	Reserved            byte
	SessName            UTF8String
	SessDesc            UTF8String
	AckTimeInterval     uint32
	AckSequenceInterval uint32
}

type GetSessionsResponse struct {
	Header        MsgHdr
	RequestId     uint16
	BlockLength   uint32
	SessionBlocks []SessionBlock
}

func (m *GetSessionsResponse) Encode() []byte {
	log.Printf("Not support to Encode!\n")
	return nil
}

func (m *GetSessionsResponse) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])

	err := binary.Read(bytesBuffer, endian, &m.Header)
	m.RequestId = binary.BigEndian.Uint16(msg[8:10])
	m.BlockLength = binary.BigEndian.Uint32(msg[10:14])

	msg = msg[14:]
	for {
		s := SessionBlock{}
		s.SessId = msg[0]
		s.Reserved = msg[1]
		msg = msg[2:]

		length := binary.BigEndian.Uint32(msg[:4])
		s.SessName.Length = length
		s.SessName.Str = make([]byte, length)
		copy(s.SessName.Str, msg[4:length+4])
		msg = msg[length+4:]

		length = binary.BigEndian.Uint32(msg[:4])
		s.SessDesc.Length = length
		s.SessDesc.Str = make([]byte, length)
		copy(s.SessDesc.Str, msg[4:length+4])
		msg = msg[length+4:]

		s.AckTimeInterval = binary.BigEndian.Uint32(msg[:4])
		s.AckSequenceInterval = binary.BigEndian.Uint32(msg[4:8])
		msg = msg[8:]

		m.SessionBlocks = append(m.SessionBlocks, s)

		if len(msg) == 0 {
			break
		}
	}

	return err
}

func (m *GetSessionsResponse) Desc() string {
	return "GET_SESSIONS_RESPONSE"
}

func (m *GetSessionsResponse) RespMsg() []IPDRMsg {
	msgs := []IPDRMsg{}
	sessIds := GetSessionList()

	for _, id := range sessIds {

		var h MsgHdr = MsgHdr{
			Version: 2,
			MsgId:   FLOW_START,
			SessId:  id,
			MsgFlag: 0,
			MsgLen:  8,
		}

		flowStart := &FlowStart{
			Header: h,
		}

		msgs = append(msgs, flowStart)

	}

	return msgs
}
