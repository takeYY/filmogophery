package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/hasher"
	"filmogophery/internal/pkg/logger"
)

type (
	IUserService interface {
		// --- Create --- //

		// ユーザーを登録する
		CreateUser(ctx context.Context, username string, email string, password string) (*types.Token, error)

		// --- Read --- //

		// --- Update --- //

		// --- Delete --- //
	}
	userService struct {
		db       *gorm.DB
		authSvc  IAuthService
		hasher   hasher.IPasswordHasher
		userRepo repositories.IUserRepository
	}
)

func NewUserService(
	db *gorm.DB,
	authSvc IAuthService,
	hasher *hasher.IPasswordHasher,
	userRepo repositories.IUserRepository,
) IUserService {
	return &userService{
		db,
		authSvc,
		*hasher,
		userRepo,
	}
}

// ユーザーを登録する
func (s *userService) CreateUser(
	ctx context.Context, username string, email string, password string,
) (*types.Token, error) {
	logger := logger.GetLogger()
	now := time.Now()

	// パスワードをハッシュ化
	pwdHash, err := s.hasher.Hash(password)
	if err != nil {
		logger.Error().Msgf("failed to hash password: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	token := &types.Token{}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// ユーザーを作成
		user := &model.Users{
			Username:     username,
			Email:        email,
			PasswordHash: string(pwdHash),
			LastLoginAt:  &now,
		}
		err := s.userRepo.Save(ctx, tx, user)
		if err != nil {
			if err.Error() == "duplicated key not allowed" {
				logger.Error().Msg("duplicated user")
				errors := make(map[string][]string)
				errors["username"] = []string{"username is already taken"}
				return responses.ConflictError("user", errors)
			}
			logger.Error().Msgf("failed to create user: %s", err.Error())
			return responses.InternalServerError()
		}

		// トークンを登録
		token, err = s.authSvc.GenerateToken(ctx, tx, user.ID)
		if err != nil {
			logger.Error().Msgf("failed to create refresh token: %s", err.Error())
			return responses.InternalServerError()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
