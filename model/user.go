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
	IsActive          bool           `json:"is_active"`
	IsVerified        bool           `json:"is_verified"`
	CreatedAt         time.Time      `json:"created_at"`
	CreatedBy         uint           `json:"created_by"`
	UpdatedAt         time.Time      `json:"updated_at"`
	UpdatedBy         uint           `json:"updated_by"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at"`
	DeletedBy         uint           `json:"deleted_by"`
}
