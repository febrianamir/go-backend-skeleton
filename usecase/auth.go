package usecase

import (
	"app/lib"
	"app/lib/constant"
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
			CreatedBy:         0,
			UpdatedAt:         timeNow,
			UpdatedBy:         0,
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

		return usecase.repo.PublishTask(ctx, constant.TaskTypeEmailSend, request.SendEmailPayload{
			To:           []string{user.Email},
			TemplateName: "register_verification.html",
			TemplateData: map[string]any{
				"code": userVerification.Code,
			},
			Subject: "Register Verification",
		})
	})
}
