package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type UserCtxKey struct{}

// https://auth0.com/docs/secure/tokens/access-tokens#sample-access-token
type AccessTokenClaims struct {
	Sub     string   `json:"sub"`
	Iss     string   `json:"iss"`
	Aud     []string `json:"aud"`
	Exp     uint     `json:"exp"`
	Iat     uint     `json:"iat"`
	IDToken string   `json:"id_token"`
}

// https://auth0.com/docs/secure/tokens/access-tokens#sample-access-token
// https://openid.net/specs/openid-connect-core-1_0.html#IDToken
type IDTokenClaims struct {
	jwt.RegisteredClaims
	IsMfaToken bool `json:"is_mfa_token"`
	UserID     uint `json:"user_id"`
}

func NewFromCtx(ctx context.Context, idTokenClaim *IDTokenClaims) context.Context {
	return context.WithValue(ctx, UserCtxKey{}, idTokenClaim)
}

func GetAuthFromCtx(ctx context.Context) *IDTokenClaims {
	if idTokenClaim, ok := ctx.Value(UserCtxKey{}).(*IDTokenClaims); ok {
		return idTokenClaim
	}
	return nil
}
