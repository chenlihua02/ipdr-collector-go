package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type Error struct {
	Header      MsgHdr
	TimeStamp   uint32
	ErrorCode   uint16
	Description UTF8String
}

func (m *Error) Encode() []byte {

	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m.Header)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	binary.Write(bytesBuffer, endian, m.TimeStamp)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	binary.Write(bytesBuffer, endian, m.ErrorCode)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	b = append(b, m.Description.Encode()...)

	//slice for msgLen
	msgLen := b[4:8]
	binary.Write(bytesBuffer, endian, uint32(len(b)))
	copy(msgLen, bytesBuffer.Bytes())

	return b
}

func (m *Error) Decode(msg []byte) error {

	bytesBuffer := bytes.NewBuffer(msg[:8])

	err := binary.Read(bytesBuffer, endian, &m.Header)
	m.TimeStamp = binary.BigEndian.Uint32(msg[8:12])
	m.ErrorCode = binary.BigEndian.Uint16(msg[12:14])
	m.Description.Length = binary.BigEndian.Uint32(msg[14:18])
	m.Description.Str = make([]byte, m.Description.Length+1)
	copy(m.Description.Str, msg[18:])

	return err
}

func (m *Error) Desc() string {
	return fmt.Sprintf("ERROR - %s, time: %d, code: %d",
		string(m.Description.Str), m.TimeStamp, m.ErrorCode)
}

func (m *Error) RespMsg() []IPDRMsg {
	return nil
}

const (
	ERR_KEEPALIVE_EXPIRED            uint16 = 0
	ERR_MSG_INVALID_FOR_CAPABILITIES uint16 = 1
	ERR_MSG_INVALID_FOR_STATE        uint16 = 2
	ERR_MSG_DECODE_ERROR             uint16 = 3
	ERR_MSG_PROCESS_TERMINATING      uint16 = 4
)

var errorCode map[uint16]string = map[uint16]string{
	ERR_KEEPALIVE_EXPIRED:            "KeepAlive expired",
	ERR_MSG_INVALID_FOR_CAPABILITIES: "Message invalid for capabilities",
	ERR_MSG_INVALID_FOR_STATE:        "Message invalid for state",
	ERR_MSG_DECODE_ERROR:             "Message decode error",
	ERR_MSG_PROCESS_TERMINATING:      "Process terminating",
}

func NewErrorMsg(errCode uint16, desc string) IPDRMsg {

	var h MsgHdr = MsgHdr{
		Version: 2,
		MsgId:   ERROR,
		SessId:  0,
		MsgFlag: 0,
		MsgLen:  8,
	}

	m := &Error{
		Header: h,
	}
	m.TimeStamp = uint32(time.Now().Unix())
	m.ErrorCode = errCode

	// if errCode is predefind, set the desc to predefined too.
	if val, ok := errorCode[errCode]; ok {
		desc = val
	}

	descStr := UTF8String{
		Length: uint32(len(desc)),
		Str:    []byte(desc),
	}

	m.Description = descStr

	return m
}
