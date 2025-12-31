package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"filmogophery/internal/pkg/config"
)

type (
	ITokenGenerator interface {
		GenerateAccessToken(userID int32, expiresInSeconds int64) (string, error)
		GenerateRefreshToken() string
	}
	tokenGenerator struct {
		secret string
	}

	JWTClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}
)

func NewTokenGenerator(conf *config.Config) ITokenGenerator {
	return &tokenGenerator{conf.Token.JWT_SECRET}
}

func (g *tokenGenerator) GenerateAccessToken(userID int32, expiresInSeconds int64) (string, error) {
	claims := JWTClaims{
		UserID: int(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresInSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secret))
}

func (g *tokenGenerator) GenerateRefreshToken() string {
	return uuid.NewString()
}
