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

		// ユーザーのログイン
		LoginUser(ctx context.Context, email, password string) (*types.Token, error)

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
			PasswordHash: pwdHash,
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
		token, err = s.authSvc.GenerateToken(ctx, tx, user.ID, now)
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

// ユーザーのログイン
func (s *userService) LoginUser(ctx context.Context, email, password string) (*types.Token, error) {
	logger := logger.GetLogger()
	now := time.Now()

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Error().Msgf("failed to fetch user: %s", err.Error())
		return nil, responses.InternalServerError()
	}
	if user == nil {
		logger.Error().Msgf("user(%s) is not found", email)
		errors := make(map[string][]string)
		errors["user"] = []string{"user is not found"}
		return nil, responses.NotFoundError("user", errors)
	}

	// パスワードを検証
	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		logger.Error().Msgf("invalid password for user(%s): %s", email, err.Error())
		errors := make(map[string][]string)
		errors["user"] = []string{"invalid email or password"}
		return nil, responses.UnauthorizedError(errors)
	}

	token := &types.Token{}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 最終ログイン日時を更新
		user.LastLoginAt = &now
		err := s.userRepo.Update(ctx, tx, user)
		if err != nil {
			logger.Error().Msgf("failed to update user: %s", err.Error())
			return responses.InternalServerError()
		}

		// 有効なトークンを無効化
		err = s.authSvc.RevokeToken(ctx, tx, user, now)
		if err != nil {
			return err
		}

		// 新しい有効なトークンを生成
		token, err = s.authSvc.GenerateToken(ctx, tx, user.ID, now)
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
