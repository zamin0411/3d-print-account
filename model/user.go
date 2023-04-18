package model

import (
	"3d-print-account/enum/sex"

	"github.com/gofrs/uuid"

	"gorm.io/gorm"
)

// Account data
type User struct {
	ID        uuid.UUID `gorm:"primaryKey;column:user_id" json:"id"`
	Email     string    `gorm:"column:user_email" json:"email"`
	Password  string    `gorm:"column:user_password" json:"password"`
	Language  string    `gorm:"user_language" json:"language"`
	Username  string    `gorm:"column:user_username" json:"username"`
	FirstName string    `gorm:"column:user_first_name" json:"firstName"`
	LastName  string    `gorm:"column:user_last_name" json:"lastName"`
	Sex       sex.Sex   `gorm:"type:enum('Female', 'Male', 'Other')" json:"sex"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == uuid.Nil {
		user.ID, _ = uuid.NewV1()
	}
	return
}
