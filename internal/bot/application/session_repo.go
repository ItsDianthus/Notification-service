package application

import "go-ItsDianthus-NotificationLink/internal/bot/domain"

type SessionRepo interface {
	GetOrCreate(userID int64) *domain.UserSession
	Save(session *domain.UserSession)
}
