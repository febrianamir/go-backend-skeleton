package usecase

import (
	"app/lib"
	"app/lib/logger"
	"app/model"
	"app/request"
	"app/response"
	"context"
	"time"

	"go.uber.org/zap"
)

func (usecase *Usecase) GetUsers(ctx context.Context, req request.GetUsers) (res response.GetUsers, err error) {
	users, total, err := usecase.repo.GetUsers(ctx, req)
	if err != nil {
		return res, err
	}

	res.Data = []response.UserList{}
	for _, user := range users {
		res.Data = append(res.Data, response.NewUserList(user))
	}
	res.Total = uint(total)
	return res, nil
}

func (usecase *Usecase) GetUser(ctx context.Context, req request.GetUser) (res response.UserDetailed, err error) {
	user, err := usecase.repo.GetUser(ctx, req)
	if err != nil {
		return res, err
	}
	if user.ID == 0 {
		notFoundError := lib.ErrorNotFound
		notFoundError.Message = "User Not Found"
		return res, notFoundError
	}

	return response.NewUserDetailed(user), nil
}

func (usecase *Usecase) CreateUser(ctx context.Context, req request.CreateUser) (err error) {
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
			zap.Strings("tags", []string{"usecase", "user"}),
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
		IsVerified:        true,
		CreatedAt:         timeNow,
		UpdatedAt:         timeNow,
	})
	return err
}

func (usecase *Usecase) UpdateUser(ctx context.Context, req request.UpdateUser) (err error) {
	user, err := usecase.repo.GetUser(ctx, request.GetUser{
		ID: req.ID,
	})
	if err != nil {
		return err
	}
	if user.ID == 0 {
		notFoundError := lib.ErrorNotFound
		notFoundError.Message = "User Not Found"
		return notFoundError
	}

	if user.Email != req.Email {
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
	}

	timeNow := time.Now()
	user.Name = req.Name
	user.Email = req.Email
	user.PhoneNumber = req.PhoneNumber
	user.UpdatedAt = timeNow
	_, err = usecase.repo.UpdateUser(ctx, user)
	return err
}

func (usecase *Usecase) DeleteUser(ctx context.Context, req request.DeleteUser) (err error) {
	return usecase.repo.DeleteUser(ctx, req.ID)
}
