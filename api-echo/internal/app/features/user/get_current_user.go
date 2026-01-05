package user

import (
	"context"

	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
)

type (
	GetCurrentUserUseCase interface {
		Run(ctx context.Context, operator *model.Users) (*types.User, error)
	}
	getCurrentUserInteractor struct{}
)

func NewGetCurrentUserInteractor() GetCurrentUserUseCase {
	return &getCurrentUserInteractor{}
}

func (i *getCurrentUserInteractor) Run(ctx context.Context, operator *model.Users) (*types.User, error) {
	return &types.User{
		ID:       operator.ID,
		Username: operator.Username,
	}, nil
}
