package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

type ConnectResponse struct {
	Header       MsgHdr
	Capabilities uint32
	KaInterval   uint32
	VendorId     UTF8String
}

func (m *ConnectResponse) Encode() []byte {
	log.Printf("Not support to Encode!\n")
	return nil
}

func (m *ConnectResponse) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])

	err := binary.Read(bytesBuffer, endian, &m.Header)
	m.Capabilities = binary.BigEndian.Uint32(msg[8:12])
	m.KaInterval = binary.BigEndian.Uint32(msg[12:16])
	m.VendorId.Length = binary.BigEndian.Uint32(msg[16:20])
	m.VendorId.Str = make([]byte, m.VendorId.Length)
	copy(m.VendorId.Str, msg[20:])

	return err
}

func (m *ConnectResponse) Desc() string {
	return "CONNECT_RESPONSE"
}

func (m *ConnectResponse) RespMsg() []IPDRMsg {

	msgs := []IPDRMsg{}
	var h MsgHdr = MsgHdr{
		Version: 2,
		MsgId:   GET_SESSIONS,
		SessId:  0,
		MsgFlag: 0,
		MsgLen:  10,
	}

	getSessions := &GetSessions{
		Header:    h,
		RequestId: 0,
	}

	msgs = append(msgs, getSessions)
	return msgs

}
