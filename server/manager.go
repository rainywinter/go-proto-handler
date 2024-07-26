package server

import (
	"fmt"
	"sync"
)

type ConnManager struct {
	conns map[ConnID]*Conn
	lock  sync.RWMutex
}

func NewConnManager() *ConnManager {
	m := &ConnManager{}
	m.conns = make(map[ConnID]*Conn)

	return m
}

func (m *ConnManager) Get(id ConnID) *Conn {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.conns[id]
}

func (m *ConnManager) Add(id ConnID, conn *Conn) {
	fmt.Println("ConnManager:Add", "id", id)
	m.lock.Lock()
	old := m.conns[id]
	m.conns[id] = conn
	m.lock.Unlock()

	if old != nil {
		fmt.Println("ConnManager:Add remove old conn", "id", id)
		_ = old.Close()
	}
}

func (m *ConnManager) Remve(id ConnID) {
	m.lock.Lock()
	conn := m.conns[id]
	delete(m.conns, id)
	m.lock.Unlock()

	if conn != nil {
		_ = conn.Close()
	}
}

func (m *ConnManager) UnexpectClose(id ConnID) {
	fmt.Println("ConnManager:UnexpectClose", "id", id)
	// todo, if bug happends, then should save the *conn, compare the conn value before delete
	m.lock.Lock()
	delete(m.conns, id)
	m.lock.Unlock()
}

func (m *ConnManager) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, c := range m.conns {
		_ = c.Close()
	}
}
