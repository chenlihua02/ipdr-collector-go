package main

import (
	"encoding/binary"
	"fmt"
	"github.com/stellar/go-xdr/xdr"
)

type TypeID uint32

const (
	//Basic Type
	INT       TypeID = 0x00000021
	UINT      TypeID = 0x00000022
	LONG      TypeID = 0x00000023
	ULONG     TypeID = 0x00000024
	FLOAT     TypeID = 0x00000025
	DOUBLE    TypeID = 0x00000026
	HEXBINARY TypeID = 0x00000027
	STRING    TypeID = 0x00000028
	BOOLEAN   TypeID = 0x00000029
	BYTE      TypeID = 0x0000002a
	UBYTE     TypeID = 0x0000002b
	SHORT     TypeID = 0x0000002c
	USHORT    TypeID = 0x0000002d
	//Derived Type
	DATETIME     TypeID = 0x00000122
	DATETIMEMSEC TypeID = 0x00000224
	IPV4ADDR     TypeID = 0x00000322
	IPV6ADDR     TypeID = 0x00000427
	IPADDR       TypeID = 0x00000827
	UUID         TypeID = 0x00000527
	DATETIMEUSEC TypeID = 0x00000623
	MACADDR      TypeID = 0x00000723
)

func XdrDecode(typeID TypeID, input []byte) (string, error) {

	var ret string

	switch typeID {
	case SHORT, USHORT:
		ret = fmt.Sprintf("%d", binary.BigEndian.Uint16(input))
	case INT, UINT, DATETIME:
		ret = fmt.Sprintf("%d", binary.BigEndian.Uint32(input))
	case LONG, ULONG, DATETIMEMSEC, DATETIMEUSEC:
		ret = fmt.Sprintf("%d", binary.BigEndian.Uint64(input))
	case FLOAT:
		var f float32
		xdr.Unmarshal(input, &f)
		ret = fmt.Sprintf("%f", f)
	case DOUBLE:
		var d float64
		xdr.Unmarshal(input, &d)
		ret = fmt.Sprintf("%f", d)
	case HEXBINARY:
		input = input[4:]
		ret = fmt.Sprintf("[% x]", input)
	case STRING:
		input = input[4:]
		ret = string(input)
	case IPV4ADDR:
		ret = fmt.Sprintf("%d.%d.%d.%d", input[0], input[1], input[2], input[3])
	case IPV6ADDR:
		input = input[4:]
		ret = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x",
			binary.BigEndian.Uint16(input),
			binary.BigEndian.Uint16(input[2:]),
			binary.BigEndian.Uint16(input[4:]),
			binary.BigEndian.Uint16(input[6:]),
			binary.BigEndian.Uint16(input[8:]),
			binary.BigEndian.Uint16(input[10:]),
			binary.BigEndian.Uint16(input[12:]),
			binary.BigEndian.Uint16(input[14:]))
	case IPADDR:
		var length = binary.BigEndian.Uint32(input)
		input = input[4:]
		if length == 4 {
			ret = fmt.Sprintf("%d.%d.%d.%d", input[0], input[1], input[2], input[3])
		} else {
			ret = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x:%02x:%02x",
				binary.BigEndian.Uint16(input),
				binary.BigEndian.Uint16(input[2:]),
				binary.BigEndian.Uint16(input[4:]),
				binary.BigEndian.Uint16(input[6:]),
				binary.BigEndian.Uint16(input[8:]),
				binary.BigEndian.Uint16(input[10:]),
				binary.BigEndian.Uint16(input[12:]),
				binary.BigEndian.Uint16(input[14:]))
		}
	case UUID:
		input = input[4:]
		ret = fmt.Sprintf("%02x%02x-%02x%02x-%02x%02x-%02x%02x",
			binary.BigEndian.Uint16(input),
			binary.BigEndian.Uint16(input[2:]),
			binary.BigEndian.Uint16(input[4:]),
			binary.BigEndian.Uint16(input[6:]),
			binary.BigEndian.Uint16(input[8:]),
			binary.BigEndian.Uint16(input[10:]),
			binary.BigEndian.Uint16(input[12:]),
			binary.BigEndian.Uint16(input[14:]))
	case MACADDR:
		input = input[2:]
		ret = fmt.Sprintf("%02x%02x.%02x%02x.%02x%02x",
			input[0], input[1], input[2], input[3], input[4], input[5])
	case BOOLEAN:
		if input[0] == 0 {
			ret = "false"
		} else {
			ret = "true"
		}
	case BYTE, UBYTE:
		ret = fmt.Sprintf("%d", input[0])
	}

	return ret, nil
}

func XdrTypeLength(typeID TypeID, len []byte) uint32 {
	switch typeID {
	case INT, UINT, FLOAT, DATETIME, IPV4ADDR:
		return 4
	case LONG, ULONG, DOUBLE, DATETIMEMSEC, DATETIMEUSEC, MACADDR:
		return 8
	case HEXBINARY, STRING, IPADDR:
		return binary.BigEndian.Uint32(len[:4]) + 4
	case BOOLEAN, BYTE, UBYTE:
		return 1
	case SHORT, USHORT:
		return 2
	case IPV6ADDR, UUID:
		return 20

	}

	return 0
}