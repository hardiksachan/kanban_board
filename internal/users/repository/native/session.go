package native

import (
	"kanban_board/internal/users"
	"kanban_board/internal/users/core/domain"
	"strconv"
	"time"
)

type SessionStore struct {
	sessions []*domain.Session
	id       int
}

func NewSessionStore() *SessionStore {
	return &SessionStore{id: 0}
}

func (s *SessionStore) Create(userId string) (*domain.Session, error) {
	session := &domain.Session{
		ID:        strconv.Itoa(s.id),
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Minute * 10),
	}

	s.id++
	s.sessions = append(s.sessions, session)

	return session, nil
}

func (s *SessionStore) Get(sessionId string) (*domain.Session, error) {
	op := "SessionStore.Get"

	pos, err := find(s.sessions, sessionId)
	if err != nil {
		return nil, &users.Error{Op: op, Err: err}
	}

	return s.sessions[pos], nil
}

func (s *SessionStore) Delete(sessionId string) error {
	op := "SessionStore.Delete"

	pos, err := find(s.sessions, sessionId)
	if err != nil {
		return &users.Error{Op: op, Err: err}
	}

	s.sessions = append(s.sessions[:pos], s.sessions[pos+1:]...)
	return nil
}

func find(sessions []*domain.Session, sessionId string) (int, error) {
	for i, session := range sessions {
		if session.ID == sessionId {
			return i, nil
		}
	}

	return -1, &users.Error{Code: users.ENOTFOUND, Message: "session does not exist"}
}
