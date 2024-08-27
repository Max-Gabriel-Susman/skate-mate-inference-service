package web

import (
	"net"
	"sync"
	"time"
)

type Session struct {
	Conn          net.Conn
	ID            string
	UserID        string
	LastActive    time.Time
	IncomingChan  chan string
	OutgoingChan  chan string
	BroadcastChan chan string
	FeedbackChan  chan bool
	Mutex         sync.Mutex
}

func NewSession(conn net.Conn, id, userID string) *Session {
	return &Session{
		Conn:          conn,
		ID:            id,
		UserID:        userID,
		LastActive:    time.Now(),
		IncomingChan:  make(chan string),
		OutgoingChan:  make(chan string),
		BroadcastChan: make(chan string),
		FeedbackChan:  make(chan bool),
	}
}
