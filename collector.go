package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	sendChan = make(chan []byte, 1)
	rcvChan  = make(chan []byte, 10)
	run      bool
)

func msgSanityCheck(msg []byte) error {

	if (len(msg)) < 8 {
		log.Printf("in sanity check, msg is not correct\n")
		return nil
	}

	if msg[0] != 2 {
		return errors.New(fmt.Sprintf("Msg %x version %d error, should be 2!",
			msg[1], msg[0]))
	}

	msg_len := binary.BigEndian.Uint32(msg[4:8])

	if msg_len != uint32(len(msg)) {
		return errors.New(fmt.Sprintf("Msg length error! (in msg: %d, msg len: %d)",
			msg_len, len(msg)))
	}

	return nil
}

func msgDecode(msg []byte) (IPDRMsg, error) {
	var rcvdMsg IPDRMsg
	var err error

	switch MessageID(msg[1]) {
	case CONNECT_RESPONSE:
		rcvdMsg = &ConnectResponse{}
	case TEMPLATE_DATA:
		rcvdMsg = &TemplateData{}
	case SESSION_START:
		rcvdMsg = &SessionStart{}
	case SESSION_STOP:
		rcvdMsg = &SessionStop{}
	case DATA:
		rcvdMsg = &Data{}
	case KEEP_ALIVE:
		rcvdMsg = &KeepAlive{}
	case ERROR:
		rcvdMsg = &Error{}
	case GET_SESSIONS_RESPONSE:
		rcvdMsg = &GetSessionsResponse{}
	default:
		log.Printf("Unsupported msg 0x%x\n", msg[1])
		return nil, errors.New("Unsupported msg.")
	}

	err = rcvdMsg.Decode(msg)
	return rcvdMsg, err
}

func handleRcvMsg(msg []byte) ([]IPDRMsg, error) {
	rcvdMsg, err := msgDecode(msg)
	if err != nil {
		return nil, err
	}

	log.Printf("Rcvd %s\n", rcvdMsg.Desc())
	RcvMsg(rcvdMsg)
	UpdateLastKaRcvdTime()
	return rcvdMsg.RespMsg(), nil
}

func receiveMsg(msg []byte) {
	err := msgSanityCheck(msg)

	if err != nil {
		log.Fatalf("Msg Error: %s", err)
		return
	}
	nextMsgs, err := handleRcvMsg(msg)
	if nextMsgs != nil {
		for _, nextMsg := range nextMsgs {
			log.Printf("Send %s\n", nextMsg.Desc())
			sendChan <- nextMsg.Encode()
		}
	}
}

func sendMsg(conn net.Conn, msg []byte) {
	_, err := conn.Write(msg)

	UpdateLastKaSendTime()

	if err != nil {
		log.Fatal("Send msg error:", err)
	}
}

func RcvMsgHandlerRoutine() {
	var buf_remain []byte
	var buf_len_remain uint32
	var msg_len uint32
	i := 0
	for {
		select {
		case buf := <-rcvChan:
			buf = append(buf_remain, buf...)
			//log.Printf("Handle msg %d, len %d\n", i, len(buf))
			i++
			buf_len_remain = uint32(len(buf))
			for {
				if len(buf) < 8 {
					break
				}
				msg_len = binary.BigEndian.Uint32(buf[4:8])
				if buf_len_remain < msg_len {
					break
				}
				receiveMsg(buf[0:msg_len])
				buf = buf[msg_len:]
				buf_remain = buf
				buf_len_remain = uint32(len(buf))
			}
		}
	}
}

func ReceiverRoutine(conn net.Conn) {
	i := 0
	for {
		if !run {
			break
		}
		buf := make([]byte, 65536)
		cnt, err := conn.Read(buf)
		if err == io.EOF {
			log.Printf("EOF Error!")
			break
		}
		if err != nil {
			log.Printf("Failed to read data from net: %s\n", err)
			continue
		}
		//log.Printf("Rcv buf %d, len %d\n", i, cnt)
		i++
		rcvChan <- buf[:cnt]
	}
}

func SenderRoutine(conn net.Conn) {
	for {
		select {
		case buf := <-sendChan:
			go sendMsg(conn, buf)
		}
	}
}

func SendMsgToExporter(buf []byte) {
	sendChan <- buf
}

func convertToIntIP(ip string) (uint32, error) {
	ips := strings.Split(ip, ".")
	E := errors.New("Not A IP.")
	if len(ips) != 4 {
		return 0, E
	}
	var intIP uint32
	for k, v := range ips {
		i, err := strconv.Atoi(v)
		if err != nil || i > 255 {
			return 0, E
		}
		intIP = intIP | uint32(i<<uint32(8*(3-k)))
	}
	return intIP, nil
}

func StartCollector(address string, port uint16, clientName string, version uint8, ka uint32) {

	var data []byte = []byte(clientName)
	var h MsgHdr = MsgHdr{
		Version: version,
		MsgId:   CONNECT,
		SessId:  0,
		MsgFlag: 0,
	}
	var vId UTF8String = UTF8String{
		Length: uint32(len(clientName)),
		Str:    data,
	}

	initAddr, err := convertToIntIP(address)
	if err != nil {
		log.Fatalf("Convert to Int IP error: %v\n", err)
		return
	}

	SetKaRecvInterval(ka)
	connect := &Connect{
		Header:       h,
		InitAddr:     initAddr,
		InitPort:     port,
		Capabilities: 2,
		KaInterval:   ka,
		VendorId:     vId,
	}

	log.Printf("Send %s\n", connect.Desc())

	sendChan <- connect.Encode()
}

func main() {
	err := ReadConfig()
	if err != nil {
		log.Fatalf("Read config file error: %v\n", err)
		return
	}

	server := GetServerAddr()
	timeout := GetConnectTimeout()

	conn, err := net.DialTimeout("tcp", server, time.Second*time.Duration(timeout))
	if err != nil {
		log.Fatalf("Fail to connect, %s\n", err)
		return
	}
	defer conn.Close()

	SessionMgrInit()

	go RcvMsgHandlerRoutine()
	go ReceiverRoutine(conn)
	go SenderRoutine(conn)

	StartCollector(GetConnectParam())

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	run = true
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		}
	}

	log.Print("Exiting collector")

}
