package model

import (
	"time"

	"gorm.io/gorm"
)

type UserVerification struct {
	ID        uint           `json:"id"`
	Type      string         `json:"type"`
	UserID    uint           `json:"user_id"`
	Code      string         `json:"code"`
	ExpiredAt *time.Time     `json:"expired_at"`
	UsedAt    *time.Time     `json:"used_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

const (
	UserVerificationTypeVerifyAccount = "VERIFY_ACCOUNT"
	UserVerificationTypeResetPassword = "RESET_PASSWORD"
)
