package auth

import (
	"app/lib/constant"
	"app/lib/logger"
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
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

func GenerateOtpSecret(ctx context.Context, userId uint, period int) (string, error) {
	identifier := strconv.Itoa(int(userId))
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constant.DefaultIssuer,
		AccountName: identifier,
		Period:      uint(period),
		Digits:      6,
	})
	if err != nil {
		logger.LogError(ctx, "totp.Generate", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"auth", "GenerateOtpSecret"}),
		}...)
		return "", err
	}

	return secret.Secret(), nil
}

func GenerateOtpCode(secret string, period int) (string, error) {
	return totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period: uint(period),
		Digits: 6,
		Skew:   1,
	})
}

func ValidateOtpCode(otpCode, secret string, period int) (bool, error) {
	return totp.ValidateCustom(otpCode, secret, time.Now(), totp.ValidateOpts{
		Period: uint(period),
		Digits: 6,
		Skew:   1,
	})
}
