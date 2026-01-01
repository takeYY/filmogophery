// internal/app/services/auth_service.go
package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/jwt"
	"filmogophery/internal/pkg/logger"
)

type (
	IAuthService interface {
		GenerateToken(ctx context.Context, tx *gorm.DB, userID int32, now time.Time) (*types.Token, error)
		RevokeToken(ctx context.Context, tx *gorm.DB, user *model.Users, now time.Time) error
	}
	authService struct {
		tokenGen  jwt.ITokenGenerator
		tokenRepo repositories.ITokenRepository
	}
)

func NewAuthService(
	tokenGen *jwt.ITokenGenerator,
	tokenRepo repositories.ITokenRepository,
) IAuthService {
	return &authService{*tokenGen, tokenRepo}
}

func (s *authService) GenerateToken(
	ctx context.Context, tx *gorm.DB, userID int32, now time.Time,
) (*types.Token, error) {
	expiresIn := int64(3600)
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	accessToken, err := s.tokenGen.GenerateAccessToken(userID, expiresIn)
	if err != nil {
		return nil, err
	}

	refreshToken := s.tokenGen.GenerateRefreshToken()
	refreshTokenHash := hashToken(refreshToken)

	// DBに保存
	token := &model.RefreshTokens{
		UserID:    userID,
		TokenHash: refreshTokenHash,
		ExpiresAt: now.Add(30 * 24 * time.Hour),
	}
	err = s.tokenRepo.Save(ctx, tx, token)
	if err != nil {
		return nil, err
	}

	return &types.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		ExpiresAt:    constant.ToUTC(expiresAt),
	}, nil
}

func (s *authService) RevokeToken(ctx context.Context, tx *gorm.DB, user *model.Users, now time.Time) error {
	logger := logger.GetLogger()

	// 有効なトークンを取得
	activeTokens, err := s.tokenRepo.FindActiveTokenByUserID(ctx, user, now)
	if err != nil {
		logger.Error().Msgf("failed to fetch active tokens: %s", err.Error())
		return responses.InternalServerError()
	}
	if len(activeTokens) == 0 {
		return nil
	}

	tokenIDs := make([]int32, 0, len(activeTokens))
	for _, t := range activeTokens {
		tokenIDs = append(tokenIDs, t.ID)
	}

	// トークンを無効化
	err = s.tokenRepo.Revoke(ctx, tx, tokenIDs, now)
	if err != nil {
		logger.Error().Msgf("failed to update active tokens: %s", err.Error())
		return responses.InternalServerError()
	}

	return nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
