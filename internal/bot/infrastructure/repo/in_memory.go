package repo

import (
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"sync"
	"time"
)

type InMemorySessionRepo struct {
	mu       sync.RWMutex
	sessions map[int64]*domain.UserSession
}

func NewInMemorySessionRepo() *InMemorySessionRepo {
	return &InMemorySessionRepo{sessions: make(map[int64]*domain.UserSession)}
}

func (r *InMemorySessionRepo) GetOrCreateSession(userID int64) *domain.UserSession {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sess, ok := r.sessions[userID]; ok {
		return sess
	}
	sess := &domain.UserSession{
		UserID:       userID,
		CurrentState: domain.StateDefault,
		TempData:     make(map[string]interface{}),
		UpdatedAt:    time.Now(),
	}
	r.sessions[userID] = sess
	return sess
}

func (r *InMemorySessionRepo) SaveSession(sess *domain.UserSession) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sess.UpdatedAt = time.Now()
	r.sessions[sess.UserID] = sess
}
