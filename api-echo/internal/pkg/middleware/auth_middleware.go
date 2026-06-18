package middleware

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/config"
	myJWT "filmogophery/internal/pkg/jwt"
)

func RequireAuthMiddleware(conf *config.Config, userRepo repositories.IUserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log := zerolog.Ctx(c.Request().Context())

			// Authorizationヘッダーからトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			// Bearer プレフィックスを削除
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return echo.ErrUnauthorized
			}

			// JWT トークンを検証
			token, err := jwt.ParseWithClaims(tokenString, &myJWT.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(conf.Token.JWT_SECRET), nil
			})
			if err != nil || !token.Valid {
				log.Error().Msgf("invalid token: %s", err.Error())
				return echo.ErrUnauthorized
			}

			claims, ok := token.Claims.(*myJWT.JWTClaims)
			if !ok {
				return echo.ErrUnauthorized
			}

			user, err := userRepo.FindByID(c.Request().Context(), int32(claims.UserID))
			if err != nil {
				log.Error().Msgf("failed to fetch user: %s", err.Error())
				return echo.ErrInternalServerError
			}
			if user == nil {
				log.Error().Msgf("user(%d) is not found", claims.UserID)
				return echo.ErrNotFound
			}

			c.Set("operator", user)
			return next(c)
		}
	}
}
