package auth

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/gen/model"
)

type (
	LogoutUseCase interface {
		Run(ctx context.Context, operator *model.Users) error
	}
	logoutInteractor struct {
		userSvc services.IUserService
	}
)

func NewLogoutInteractor(userSvc services.IUserService) LogoutUseCase {
	return &logoutInteractor{userSvc}
}

func (i *logoutInteractor) Run(ctx context.Context, operator *model.Users) error {
	return i.userSvc.LogoutUser(ctx, operator)
}
