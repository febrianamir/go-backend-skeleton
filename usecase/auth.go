package usecase

import (
	"app/lib"
	"app/lib/auth"
	"app/lib/constant"
	"app/lib/logger"
	"app/model"
	"app/request"
	"app/response"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (usecase *Usecase) Register(ctx context.Context, req request.Register) (err error) {
	checkUserEmail, err := usecase.repo.GetUser(ctx, request.GetUser{
		Email: req.Email,
	})
	if err != nil {
		return err
	}
	if checkUserEmail.ID > 0 {
		validationError := lib.ErrorValidation
		validationError.ErrDetails = map[string]any{
			"email": "Email already registered",
		}
		return validationError
	}

	encryptedPassword, err := lib.GeneratePasswordHash(req.Password)
	if err != nil {
		logger.LogError(ctx, "Error GeneratePasswordHash", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"usecase", "Register"}),
		}...)
		return lib.ErrorInternalServer
	}

	return usecase.repo.Transaction(ctx, func(ctx context.Context) error {
		timeNow := time.Now()
		user, err := usecase.repo.CreateUser(ctx, model.User{
			Name:              req.Name,
			Email:             req.Email,
			PhoneNumber:       req.PhoneNumber,
			EncryptedPassword: encryptedPassword,
			IsActive:          true,
			IsVerified:        false,
			CreatedAt:         timeNow,
			UpdatedAt:         timeNow,
		})
		if err != nil {
			return err
		}

		userVerification, err := usecase.repo.CreateUserVerification(ctx, model.UserVerification{
			Type:      model.UserVerificationTypeVerifyAccount,
			UserID:    user.ID,
			Code:      lib.GenerateUUID(),
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		})
		if err != nil {
			return err
		}

		err = usecase.repo.PublishTask(ctx, constant.TaskTypeEmailSend, request.SendEmailPayload{
			To:           []string{user.Email},
			TemplateName: "register_verification.html",
			TemplateData: map[string]any{
				"code": userVerification.Code,
			},
			Subject: "Register Verification",
		})
		if err != nil {
			return err
		}

		return usecase.repo.SetVerificationDelayCache(ctx, user.ID, model.UserVerificationTypeVerifyAccount)
	})
}

func (usecase *Usecase) RegisterResendVerification(ctx context.Context, req request.RegisterResendVerification) (err error) {
	user, err := usecase.repo.GetUser(ctx, request.GetUser{
		Email: req.Email,
	})
	if err != nil {
		return err
	}
	if user.ID == 0 {
		notFoundError := lib.ErrorNotFound
		notFoundError.Message = "User Not Found"
		return notFoundError
	}

	verificationDelayCache, remainingTtl, err := usecase.repo.GetVerificationDelayCacheWithTtl(ctx, user.ID, model.UserVerificationTypeVerifyAccount)
	if err != nil {
		return err
	}
	if verificationDelayCache != "" {
		verificationDelayError := lib.ErrorVerificationDelay
		verificationDelayError.ErrDetails = map[string]any{
			"remaining_ttl": remainingTtl / time.Second,
		}
		return verificationDelayError
	}

	userVerification, err := usecase.repo.GetUserVerification(ctx, request.GetUserVerification{})
	if err != nil {
		return err
	}
	if userVerification.ID == 0 {
		notFoundError := lib.ErrorNotFound
		notFoundError.Message = "User Verification Not Found"
		return notFoundError
	}

	return usecase.repo.Transaction(ctx, func(ctx context.Context) error {
		err = usecase.repo.PublishTask(ctx, constant.TaskTypeEmailSend, request.SendEmailPayload{
			To:           []string{user.Email},
			TemplateName: "register_verification.html",
			TemplateData: map[string]any{
				"code": userVerification.Code,
			},
			Subject: "Register Verification",
		})
		if err != nil {
			return err
		}

		return usecase.repo.SetVerificationDelayCache(ctx, user.ID, model.UserVerificationTypeVerifyAccount)
	})
}

func (usecase *Usecase) VerifyAccount(ctx context.Context, req request.VerifyAccount) (response.Auth, error) {
	userVerification, err := usecase.repo.GetUserVerification(ctx, request.GetUserVerification{
		Code: req.Code,
	})
	if err != nil {
		return response.Auth{}, err
	}
	if userVerification.ID == 0 {
		notFoundError := lib.ErrorNotFound
		notFoundError.Message = "User Verification Not Found"
		return response.Auth{}, notFoundError
	}

	if userVerification.ExpiredAt != nil && !userVerification.ExpiredAt.IsZero() {
		return response.Auth{}, lib.ErrorVerificationInactive
	}

	if userVerification.UsedAt != nil && !userVerification.UsedAt.IsZero() {
		return response.Auth{}, lib.ErrorVerificationInactive
	}

	user, err := usecase.repo.GetUser(ctx, request.GetUser{
		ID: userVerification.UserID,
	})
	if err != nil {
		return response.Auth{}, err
	}
	if user.ID == 0 {
		notFoundError := lib.ErrorNotFound
		notFoundError.Message = "User Not Found"
		return response.Auth{}, notFoundError
	}

	var auth model.UserAuth
	err = usecase.repo.Transaction(ctx, func(ctx context.Context) error {
		timeNow := time.Now()

		user.IsActive = true
		user.IsVerified = true
		user.UpdatedAt = timeNow
		_, err := usecase.repo.UpdateUser(ctx, user)
		if err != nil {
			return err
		}

		userVerification.ExpiredAt = &timeNow
		userVerification.UsedAt = &timeNow
		userVerification.UpdatedAt = timeNow
		_, err = usecase.repo.UpdateUserVerification(ctx, userVerification)
		if err != nil {
			return err
		}

		// Generate non-mfa auth
		auth, _, err = usecase.generateAuth(ctx, user, false)
		if err != nil {
			return err
		}

		return usecase.repo.SetMfaFlag(ctx, user.ID)
	})
	if err != nil {
		return response.Auth{}, err
	}

	return response.NewAuth(auth, user, false), nil
}

func (usecase *Usecase) Login(ctx context.Context, req request.Login) (res response.Auth, err error) {
	user, err := usecase.repo.GetUser(ctx, request.GetUser{
		Email: req.Email,
	})
	if err != nil {
		return res, err
	}
	if user.ID == 0 {
		notFoundError := lib.ErrorNotFound
		notFoundError.Message = "User Not Found"
		return res, notFoundError
	}

	err = lib.CompareHashAndPassword(user.EncryptedPassword, req.Password)
	if err != nil {
		return res, lib.ErrorWrongCredential
	}

	isNeedMfa, err := usecase.isNeedMfa(ctx, user.ID)
	if err != nil {
		return res, err
	}

	auth, _, err := usecase.generateAuth(ctx, user, isNeedMfa)
	if err != nil {
		return res, err
	}

	return response.NewAuth(auth, user, isNeedMfa), nil
}

func (usecase *Usecase) GetIDToken(ctx context.Context, accessToken string) (string, error) {
	accessTokenClaims, err := usecase.repo.GetAccessToken(ctx, accessToken)
	if err != nil {
		return "", err
	}
	return accessTokenClaims.IDToken, err
}

func (usecase *Usecase) ParseIDToken(ctx context.Context, idToken string) (*auth.IDTokenClaims, error) {
	token, err := jwt.ParseWithClaims(idToken, &auth.IDTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(usecase.config.ID_TOKEN_HMAC_KEY), nil
	})
	if err != nil {
		logger.LogError(ctx, "Error jwt.ParseWithClaims", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"usecase", "ParseIDToken"}),
		}...)
		return nil, err
	}

	if idTokenClaim, ok := token.Claims.(*auth.IDTokenClaims); ok {
		return idTokenClaim, nil
	}

	logger.LogError(ctx, "Failed parse id token claims", []zap.Field{
		zap.Strings("tags", []string{"usecase", "ParseIDToken"}),
	}...)
	return nil, errors.New("Failed parse id token claims")
}

func (usecase *Usecase) isNeedMfa(ctx context.Context, userId uint) (bool, error) {
	mfaFlag, err := usecase.repo.GetMfaFlag(ctx, userId)
	if err != nil {
		return true, err
	}

	if mfaFlag == "" {
		return true, nil
	}

	return false, nil
}

func (usecase *Usecase) generateAuth(ctx context.Context, user model.User, isNeedMfa bool) (model.UserAuth, bool, error) {
	accessToken, idToken, accessTokenExp, idTokenExp, err := usecase.generateAuthToken(ctx, user, isNeedMfa)
	if err != nil {
		return model.UserAuth{}, false, err
	}

	auth := model.UserAuth{
		UserID:               user.ID,
		AccessToken:          accessToken,
		IDToken:              idToken,
		AccessTokenExpiredAt: accessTokenExp,
		IDTokenExpiredAt:     idTokenExp,
	}

	if !isNeedMfa {
		auth, err = usecase.repo.CreateAuth(ctx, auth)
		if err != nil {
			return model.UserAuth{}, false, err
		}
	}

	return auth, isNeedMfa, nil
}

func (usecase *Usecase) generateAuthToken(ctx context.Context, user model.User, isMfaToken bool) (accessToken, idToken string, accessTokenExp, idTokenExp time.Time, err error) {
	timeNow := time.Now()

	idTokenExp = timeNow.Add(time.Duration(usecase.config.ID_TOKEN_TTL) * time.Second)
	idToken, err = usecase.generateIDToken(ctx, auth.IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(idTokenExp),
			IssuedAt:  jwt.NewNumericDate(timeNow),
			NotBefore: jwt.NewNumericDate(timeNow),
			Issuer:    constant.DefaultIssuer,
			Subject:   strconv.Itoa(int(user.ID)),
			Audience:  []string{constant.DefaultAudience},
		},
		UserID:     user.ID,
		IsMfaToken: isMfaToken,
	})
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	accessTokenTtl := usecase.config.ACCESS_TOKEN_TTL
	if isMfaToken {
		accessTokenTtl = usecase.config.MFA_ACCESS_TOKEN_TTL
	}

	accessTokenExp = timeNow.Add(time.Duration(accessTokenTtl) * time.Second)
	accessToken, err = usecase.generateAccessToken(ctx, auth.AccessTokenClaims{
		Sub:     strconv.Itoa(int(user.ID)),
		Exp:     uint(accessTokenExp.Unix()),
		IDToken: idToken,
	})
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	// TODO: generate refresh token

	return accessToken, idToken, accessTokenExp, idTokenExp, err
}

func (usecase *Usecase) generateAccessToken(ctx context.Context, claims auth.AccessTokenClaims) (string, error) {
	if claims.Sub == "" {
		err := errors.New("Sub must be set")
		logger.LogError(ctx, "Error generateAccessToken", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"usecase", "generateAccessToken"}),
		}...)
		return "", err
	}

	claims.Iat = uint(time.Now().Unix())

	if claims.Exp == 0 {
		err := errors.New("Exp must be set")
		logger.LogError(ctx, "Error generateAccessToken", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"usecase", "generateAccessToken"}),
		}...)
		return "", err
	}

	if claims.Iss == "" {
		claims.Iss = constant.DefaultIssuer
	}

	if len(claims.Aud) == 0 {
		claims.Aud = append(claims.Aud, constant.DefaultAudience)
	}

	accessToken := uuid.New().String()
	err := usecase.repo.SetAccessToken(ctx, accessToken, claims)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (usecase *Usecase) generateIDToken(ctx context.Context, claims auth.IDTokenClaims) (string, error) {
	idToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedIDToken, err := idToken.SignedString([]byte(usecase.config.ID_TOKEN_HMAC_KEY))
	if err != nil {
		logger.LogError(ctx, "error idToken.SignedString", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"usecase", "generateIDToken"}),
		}...)
		return "", err
	}
	return signedIDToken, nil
}
