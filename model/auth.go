package model

import (
	"time"

	"gorm.io/gorm"
)

type UserAuth struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	UserID                uint           `json:"user_id"`
	AccessToken           string         `json:"access_token"`
	RefreshToken          string         `json:"refresh_token"`
	IDToken               string         `json:"id_token"`
	AccessTokenExpiredAt  time.Time      `json:"access_token_expired_at"`
	RefreshTokenExpiredAt time.Time      `json:"refresh_token_expired_at"`
	IDTokenExpiredAt      time.Time      `json:"id_token_expired_at"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at"`
}
