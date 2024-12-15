package session

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type SessionM interface {
	ValidateSession(sessionID string) bool
	CreateSession(sessionId string, t time.Duration, data string) string
	GetSessionData(sessionID string) (string, bool)
	RemoveSession(sessionID string)
}

type Session struct {
	session map[string]sessionData
	mu      sync.Mutex
}

type sessionData struct {
	expiry time.Time
	data   string
}

// NewSession creates a new session and initializes it.
func NewSession() SessionM {
	s := &Session{
		session: make(map[string]sessionData),
	}
	s.init(time.Minute * 15)
	return s
}

// ValidateSession validates a session.
func (s *Session) ValidateSession(sessionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.session[sessionID]
	if !ok || time.Now().After(data.expiry) {
		return false
	}
	return true
}

// CreateSession creates a new session.
// t is the duration for which the session is valid.
func (s *Session) CreateSession(sessionId string, t time.Duration, data string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.session[sessionId] = sessionData{
		expiry: time.Now().Add(t),
		data:   data,
	}
	s.writeSession()
	return s.session[sessionId].expiry.String()
}

// GetSessionData retrieves the data associated with a session.
func (s *Session) GetSessionData(sessionID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.session[sessionID]
	if !ok || time.Now().After(data.expiry) {
		return "", false
	}
	return data.data, true
}

func (s *Session) RemoveSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.session, sessionID)
	s.writeSession()
}

// sessionCleanup cleans up expired sessions.
func (s *Session) sessionCleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.session {
		if time.Now().After(v.expiry) {
			delete(s.session, k)
		}
	}
}

// writeSession to file
func (s *Session) writeSession() {
	file := "session.txt"
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	for k, v := range s.session {
		_, err := f.WriteString(fmt.Sprintf("%s %s %s\n", k, v.expiry.String(), v.data))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// readSession from file
func (s *Session) readSession() {
	file := "session.txt"
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var sessionID, sessionTime, data string
	for {
		_, err := fmt.Fscanf(f, "%s %s %s\n", &sessionID, &sessionTime, &data)
		if err != nil {
			break
		}

		t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", sessionTime)
		if err != nil {
			fmt.Println(err)
			return
		}

		s.session[sessionID] = sessionData{
			expiry: t,
			data:   data,
		}
	}
}

// init initializes the session package.
// t is the period after which the session cleanup process is triggered.
func (s *Session) init(t time.Duration) {
	s.readSession()

	go func() {
		for {
			s.sessionCleanup()
			time.Sleep(t)
		}
	}()
}
