package main

import (
	"bytes"
	"encoding/binary"
)

var (
	endian = binary.BigEndian
)

type MessageID byte

const (
	CONNECT                  MessageID = 0x05
	CONNECT_RESPONSE         MessageID = 0x06
	DISCONNECT               MessageID = 0x07
	FLOW_START               MessageID = 0x01
	FLOW_STOP                MessageID = 0x03
	SESSION_START            MessageID = 0x08
	SESSION_STOP             MessageID = 0x09
	KEEP_ALIVE               MessageID = 0x40
	TEMPLATE_DATA            MessageID = 0x10
	MODIFY_TEMPLATE          MessageID = 0x1a
	MODIFY_TEMPLATE_RESPONSE MessageID = 0x1b
	FINAL_TEMPLATE_DATA_ACK  MessageID = 0x13
	START_NEGOTIATION        MessageID = 0x1d
	START_NEGOTIATION_REJECT MessageID = 0x1e
	GET_SESSIONS             MessageID = 0x14
	GET_SESSIONS_RESPONSE    MessageID = 0x15
	GET_TEMPLATES            MessageID = 0x16
	GET_TEMPLATES_RESPONSE   MessageID = 0x17
	DATA                     MessageID = 0x20
	DATA_ACK                 MessageID = 0x21
	ERROR                    MessageID = 0x23
)

type IPDRMsg interface {
	// Encode the struct to byte array.
	Encode() []byte
	// Decode the byte array to struct.
	Decode([]byte) error
	// Msg desc
	Desc() string
	//Construct response messages if have.
	RespMsg() []IPDRMsg
}

type UTF8String struct {
	Length uint32
	Str    []byte
}

func (s *UTF8String) Encode() []byte {

	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, s.Length)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	b = append(b, s.Str...)

	return b
}

type MsgHdr struct {
	Version byte
	MsgId   MessageID
	SessId  byte
	MsgFlag byte
	MsgLen  uint32
}
