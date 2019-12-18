package main

import (
	"log"
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

	sessions[m.Header.SessId] = s
}

func UpdateSession(m *SessionStart) {
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
	} else {
		log.Printf("Session %d not exist internal when handle start session.\n", sessId)

	}
}

func RemoveSession(m *SessionStop) {
	sessId := m.Header.SessId

	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if s, ok := sessions[sessId]; ok {
		//Didn't remove from map, just mark a flag
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
		id := t.Header.SessId
		sessions[id].LastSeq = t.SequenceNum
		sessions[id].ConfigId = t.ConfigID
		sessions[id].UnackedNum++
		sessions[id].CheckSequenceInterval()
	case *TemplateData:
		AddSession(t)
	case *SessionStart:
		UpdateSession(t)
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
