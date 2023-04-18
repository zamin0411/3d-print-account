package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginSession struct {
	ID    uuid.UUID  `gorm:"column:login_session_id"`
	Token string     `gorm:"column:login_session_token"`
	Time  *time.Time `gorm:"column:login_session_time"`
}

func (session *LoginSession) BeforeCreate(tx *gorm.DB) (err error) {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	return
}
