// internal/app/services/auth_service.go
package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/jwt"
)

type (
	IAuthService interface {
		GenerateToken(ctx context.Context, tx *gorm.DB, userID int32) (*types.Token, error)
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
	ctx context.Context, tx *gorm.DB, userID int32,
) (*types.Token, error) {
	expiresIn := int64(3600)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

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
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	err = s.tokenRepo.Save(ctx, tx, token)
	if err != nil {
		return nil, err
	}

	return &types.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiresAt:    constant.ToUTC(expiresAt),
	}, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
