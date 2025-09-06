package auth

import "github.com/golang-jwt/jwt/v5"

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
