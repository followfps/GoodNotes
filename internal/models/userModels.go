package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name       string
	Email      string `gorm:"unique"`
	Password   string
	BucketName string
	UserID     uuid.UUID `gorm:"unique"`
}

func (u *Users) GetName() string {
	return u.Name
}

func (u *Users) GetEmail() string {
	return u.Email
}

func (u *Users) GetPass() string {
	return u.Password
}

func (u *Users) GetBucketName() string {
	return u.BucketName
}
