package message

import (
	"fmt"
	"sync"
)

type MessageObj struct {
	Id           int64  `json:"id"`
	Text         string `json:"text"`
	IsPalindrome *bool  `json:"is-palindrome,omitempty"`
}

type MessageServer struct {
	sync.RWMutex
	latestID int64
	msgStore map[int64]string
}

func NewMessageServer() *MessageServer {
	m := MessageServer{}
	m.msgStore = make(map[int64]string)
	return &m
}

func (m *MessageServer) Add(msg string) MessageObj {
	m.Lock()
	defer m.Unlock()
	m.latestID++
	m.msgStore[m.latestID] = msg
	resp := MessageObj{Id: m.latestID, Text: msg}
	return resp
}

func (m *MessageServer) Get(id int64) (MessageObj, error) {
	m.RLock()
	defer m.RUnlock()
	msg, ok := m.msgStore[id]
	if !ok {
		return MessageObj{}, fmt.Errorf("input text ID %v not found ", id)
	}
	return MessageObj{Id: id, Text: msg}, nil
}

func (m *MessageServer) GetAll() ([]MessageObj, error) {
	m.RLock()
	if len(m.msgStore) == 0 {
		return nil, fmt.Errorf("no messages found")
	}
	defer m.RUnlock()
	msgList := make([]MessageObj, len(m.msgStore))
	index := 0
	for id := range m.msgStore {
		msgList[index] = MessageObj{Id: id, Text: m.msgStore[id]}
		index++
	}
	return msgList, nil
}

func (m *MessageServer) Delete(id int64) error {
	m.Lock()
	defer m.Unlock()
	_, ok := m.msgStore[id]
	if !ok {
		return fmt.Errorf("input text ID %v not found ", id)
	}
	delete(m.msgStore, id)
	return nil
}
