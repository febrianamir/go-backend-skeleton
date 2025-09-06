package response

import (
	"app/model"
	"time"
)

type Login struct {
	IsNeedMfa             bool         `json:"is_need_mfa"`
	UserID                uint         `json:"user_id"`
	User                  UserDetailed `json:"user"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`
}

func NewLogin(auth model.UserAuth, user model.User, isNeedMfa bool) Login {
	return Login{
		IsNeedMfa:             isNeedMfa,
		UserID:                auth.UserID,
		User:                  NewUserDetailed(user),
		AccessToken:           auth.AccessToken,
		RefreshToken:          auth.RefreshToken,
		AccessTokenExpiredAt:  auth.AccessTokenExpiredAt,
		RefreshTokenExpiredAt: auth.RefreshTokenExpiredAt,
		CreatedAt:             auth.CreatedAt,
		UpdatedAt:             auth.UpdatedAt,
	}
}
