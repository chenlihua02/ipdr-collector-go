package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	msgChan                     = make(chan IPDRMsg)
	sessions                    = make(map[byte]*Session)
	sessionMutex   sync.RWMutex = sync.RWMutex{}
	kaSendInterval uint32       = 300
	kaRecvInterval uint32       = 300 //default 300 seconds
	lastKaSendTime time.Time
)

type Field struct {
	TypeID    uint32
	FieldID   uint32
	FieldName string
	IsEnabled bool
}

type Template struct {
	TemplateID uint16
	SchemaName string
	TypeName   string
	Fields     []*Field
	FileName   string
	Output     *os.File
}

type Session struct {
	Id                  byte
	Type                byte
	AckSequenceInterval uint32
	AckTimeInterval     uint32
	ConfigId            uint16
	UnackedNum          uint32
	LastSeq             uint64
	LastAckedTime       time.Time
	Started             bool
	DocID               []byte
	Templates           []*Template
}

func (t *Template) appendRecord(r []byte) error {
	var length uint32
	var typeId TypeID
	var err error
	first := true
	r = r[4:]
	bufferedWriter := bufio.NewWriter(t.Output)
	for _, f := range t.Fields {
		typeId = TypeID(f.TypeID)
		length = XdrTypeLength(typeId, r)
		str, _ := XdrDecode(typeId, r[:length])
		if first {
			_, err = bufferedWriter.WriteString(str)
			first = false
		} else {
			_, err = bufferedWriter.WriteString("," + str)
		}
		if err != nil {
			log.Printf("appendRecord err: %s\n", err)
			break
		}
		//log.Printf("F.ID %x F.FieldID %d, F.Name %s, length %d, r: % x, str %s\n",
		//	f.TypeID, f.FieldID, f.FieldName, length, r[:length], str)
		//log.Printf("remain len %d, len %d\n", len(r), length)
		if uint32(len(r)) > length {
			r = r[length:]
		} else {
			break
		}
	}
	_, err = bufferedWriter.WriteString("\n")
	if err != nil {
		log.Printf("appendRecord err: %s\n", err)
	}
	bufferedWriter.Flush()
	bufferedWriter.Reset(bufferedWriter)
	return nil
}

func (s *Session) sendAck() {

	s.UnackedNum = 0
	s.LastAckedTime = time.Now()
	msg := NewDataAckMsg(s.ConfigId, s.Id, s.LastSeq)
	log.Printf("Send %s\n", msg.Desc())
	SendMsgToExporter(msg.Encode())
}

func (s *Session) CheckSequenceInterval() {
	if s.UnackedNum >= s.AckSequenceInterval {
		s.sendAck()
	}
}

func (s *Session) CheckAckTimeInterval() {
	if time.Since(s.LastAckedTime) >= time.Duration(s.AckTimeInterval)*time.Second {
		s.sendAck()
	}
}

func UpdateLastKaSendTime() {
	lastKaSendTime = time.Now()
}

var kaRecvTimer = time.NewTimer(time.Second * time.Duration(kaRecvInterval))

func UpdateLastKaRcvdTime() {
	kaRecvTimer.Reset(time.Second * time.Duration(kaRecvInterval))
}

func SetKaRecvInterval(ka uint32) {
	kaRecvInterval = ka + 2
	log.Printf("Set KA recv interval to %d\n", ka)
}

func CheckKeepAliveInterval() {
	//Send KA
	if time.Since(lastKaSendTime) >= time.Duration(kaSendInterval)*time.Second {
		UpdateLastKaSendTime()
		msg := NewKeepAliveMsg()
		log.Printf("Send %s\n", msg.Desc())
		SendMsgToExporter(msg.Encode())
	}

}

func AddSession(m *TemplateData) {
	sessId := m.Header.SessId

	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	s := &Session{
		Id:       sessId,
		ConfigId: m.ConfigID,
	}

	for _, tb := range m.Templates {
		t := &Template{}
		t.TemplateID = tb.TemplateID
		t.SchemaName = string(tb.SchemaName.Str)
		t.TypeName = string(tb.TypeName.Str)
		t.Output = nil
		t.FileName = ""
		log.Printf("set sess %d template %d name to null", sessId, t.TemplateID)
		for _, fd := range tb.Fields {
			f := &Field{}
			f.TypeID = fd.TypeID
			f.FieldID = fd.FieldID
			f.FieldName = string(fd.FieldName.Str)
			if fd.IsEnabled == 0 {
				f.IsEnabled = false
			} else {
				f.IsEnabled = true
			}
			t.Fields = append(t.Fields, f)
		}
		s.Templates = append(s.Templates, t)
	}

	//log.Printf("Add session % +v\n", s)

	sessions[m.Header.SessId] = s
}

func createFileTemplate(t *Template, s *Session) {
	fileName := fmt.Sprintf("IPDR_RECORD_%s_S%d_T%d_%s.csv",
		time.Now().Format("2006-01-02-15-04-05"), s.Id, t.TemplateID, t.TypeName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Err: %s\n", err)
		return
	}
	t.FileName = fileName
	log.Printf("set sess %d template %d name to %s", s.Id, t.TemplateID, fileName)
	t.Output = file
	bufferedWriter := bufio.NewWriter(file)
	first := true
	for _, f := range t.Fields {
		if first {
			_, err = bufferedWriter.WriteString(f.FieldName)
			first = false
		} else {
			_, err = bufferedWriter.WriteString("," + f.FieldName)
		}
	}
	_, err = bufferedWriter.WriteString("\n")
	bufferedWriter.Flush()
	bufferedWriter.Reset(bufferedWriter)

}

func createFiles(s *Session) {
	for _, t := range s.Templates {
		createFileTemplate(t, s)
	}
}

func closeFiles(s *Session) {
	for _, t := range s.Templates {
		t.Output.Close()
		t.FileName = ""
	}
}

func StartSession(m *SessionStart) {
	sessId := m.Header.SessId

	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if s, ok := sessions[sessId]; ok {
		s.Type = 0
		s.AckSequenceInterval = m.AckSequenceInterval
		s.AckTimeInterval = m.AckTimeInterval
		s.UnackedNum = 0
		s.LastSeq = 0
		s.LastAckedTime = time.Now()
		s.Started = true
		s.DocID = make([]byte, 16)
		copy(s.DocID, m.DocumentID)
		createFiles(s)
	} else {
		log.Printf("Session %d not exist internal when handle start session.\n", sessId)

	}
}

func UpdateSession(d *Data) {
	sessId := d.Header.SessId
	if s, ok := sessions[sessId]; ok {
		sessions[sessId].LastSeq = d.SequenceNum
		sessions[sessId].ConfigId = d.ConfigID
		sessions[sessId].UnackedNum++
		sessions[sessId].CheckSequenceInterval()
		for _, t := range s.Templates {
			if t.TemplateID == d.TemplateID {
				t.appendRecord(d.Record)
				break
			}
		}
	}
}

func RemoveSession(m *SessionStop) {
	sessId := m.Header.SessId

	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if s, ok := sessions[sessId]; ok {
		//Didn't remove from map, just mark a flag
		closeFiles(s)
		s.Started = false
	}
}

func RcvMsg(msg IPDRMsg) {
	msgChan <- msg
}

func handleTimerEvt() {
	CheckKeepAliveInterval()
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	for _, s := range sessions {
		if s.Started {
			s.CheckAckTimeInterval()
		}
	}
}

func handleKaTimeout() {
	msg := NewErrorMsg(ERR_KEEPALIVE_EXPIRED, "")
	log.Printf("Send %s\n", msg.Desc())
	SendMsgToExporter(msg.Encode())
}

func handleMsg(msg IPDRMsg) {

	switch t := msg.(type) {
	case *Data:
		UpdateSession(t)
	case *TemplateData:
		AddSession(t)
	case *SessionStart:
		StartSession(t)
	case *SessionStop:
		RemoveSession(t)
	case *KeepAlive:
	case *ConnectResponse:
		kaSendInterval = t.KaInterval
		log.Printf("Set KA send interval to %d\n", kaSendInterval)
		if kaSendInterval >= 5 {
			//Send KA 2 seconds before interval.
			kaSendInterval -= 2
		}
	}

}

func SessionMgrInit() {
	var sessTimer = time.NewTimer(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case <-sessTimer.C:
				handleTimerEvt()
				sessTimer.Reset(time.Second * time.Duration(1))
			case <-kaRecvTimer.C:
				handleKaTimeout()
				kaRecvTimer.Reset(time.Second * time.Duration(kaRecvInterval))
			case msg := <-msgChan:
				handleMsg(msg)
			}
		}
	}()
}