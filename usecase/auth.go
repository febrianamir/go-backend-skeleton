package usecase

import (
	"app/lib"
	"app/lib/logger"
	"app/model"
	"app/request"
	"context"
	"time"

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

	timeNow := time.Now()
	_, err = usecase.repo.CreateUser(ctx, model.User{
		Name:              req.Name,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		EncryptedPassword: encryptedPassword,
		IsActive:          true,
		IsVerified:        false,
		CreatedAt:         timeNow,
		CreatedBy:         0,
		UpdatedAt:         timeNow,
		UpdatedBy:         0,
	})
	return err
}
