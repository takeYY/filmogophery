package user

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
)

type (
	CreateUserUseCase interface {
		Run(ctx context.Context, username string, email string, password string) (*types.Token, error)
	}
	createUserInteractor struct {
		userSvc services.IUserService
	}
)

func NewCreateUserInteractor(
	userSvc services.IUserService,
) CreateUserUseCase {
	return &createUserInteractor{
		userSvc,
	}
}

func (i *createUserInteractor) Run(
	ctx context.Context, username string, email string, password string,
) (*types.Token, error) {
	return i.userSvc.CreateUser(ctx, username, email, password)
}
