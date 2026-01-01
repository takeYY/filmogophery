package auth

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
)

type (
	LoginUseCase interface {
		Run(ctx context.Context, email, password string) (*types.Token, error)
	}

	loginInteractor struct {
		userSvc services.IUserService
	}
)

func NewLoginInteractor(
	userSvc services.IUserService,
) LoginUseCase {
	return &loginInteractor{
		userSvc,
	}
}

func (i *loginInteractor) Run(ctx context.Context, email, password string) (*types.Token, error) {
	return i.userSvc.LoginUser(ctx, email, password)
}
