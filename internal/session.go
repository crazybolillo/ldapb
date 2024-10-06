package internal

import (
	"crypto/sha256"
	"net"
	"sync"
)

type Session struct {
	store map[[32]byte]string
	lock  sync.Mutex
}

func NewSession() *Session {
	return &Session{
		store: make(map[[32]byte]string),
	}
}

func connectionID(conn net.Conn) [32]byte {
	return sha256.Sum256([]byte(conn.LocalAddr().String() + conn.RemoteAddr().String()))
}

func (s *Session) Get(conn net.Conn) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.store[connectionID(conn)]
}

func (s *Session) Add(conn net.Conn, dn string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.store[connectionID(conn)] = dn
}

func (s *Session) Remove(conn net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.store, connectionID(conn))
}
