package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"filmogophery/internal/pkg/logger"
)

// RequestLoggerMiddleware はリクエストごとに requestId を付与した child logger を
// context.Context に埋め込み、レスポンスヘッダー X-Request-ID にも設定します。
// 各層では zerolog.Ctx(ctx) でこのロガーを取り出せます。
func RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := uuid.NewString()

		// zerolog の child logger を生成して context に付与
		log := logger.GetLogger().With().Str("requestId", requestID).Logger()
		ctx := log.WithContext(c.Request().Context())
		c.SetRequest(c.Request().WithContext(ctx))
		c.Response().Header().Set("X-Request-ID", requestID)

		start := time.Now()
		err := next(c)

		// レスポンス後にリクエストログを出力
		event := log.Info()
		if err != nil {
			event = log.Error().Err(err)
		}
		event.
			Str("method", c.Request().Method).
			Str("uri", c.Request().RequestURI).
			Int("status", c.Response().Status).
			Dur("durationMs", time.Since(start)).
			Msg("request completed")

		return err
	}
}
