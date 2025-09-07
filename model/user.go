package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	PhoneNumber       string         `json:"phone_number"`
	EncryptedPassword string         `json:"encrypted_password"`
	OtpSecret         string         `json:"otp_secret"`
	IsActive          bool           `json:"is_active"`
	IsVerified        bool           `json:"is_verified"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at"`
}
