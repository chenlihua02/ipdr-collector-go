package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type DataAck struct {
	Header      MsgHdr
	ConfigID    uint16
	SequenceNum uint64
}

func (m *DataAck) Encode() []byte {
	b := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})

	binary.Write(bytesBuffer, endian, m)
	b = append(b, bytesBuffer.Bytes()...)
	bytesBuffer.Reset()

	return b
}

func (m *DataAck) Decode(msg []byte) error {

	return errors.New(fmt.Sprintf("Not support decode for %s", m.Desc()))
}

func (m *DataAck) Desc() string {
	return fmt.Sprintf("DATA_ACK - id: %d, seq: %d",
		m.Header.SessId, m.SequenceNum)
}

func (m *DataAck) RespMsg() []IPDRMsg {

	return nil
}

func NewDataAckMsg(configId uint16, sessId byte, seqNum uint64) IPDRMsg {
	var h MsgHdr = MsgHdr{
		Version: 2,
		MsgId:   DATA_ACK,
		SessId:  sessId,
		MsgFlag: 0,
		MsgLen:  18,
	}

	dataAck := &DataAck{
		Header:      h,
		ConfigID:    configId,
		SequenceNum: seqNum,
	}

	return dataAck
}
