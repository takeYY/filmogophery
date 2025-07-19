package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	genreHandlers "filmogophery/internal/app/features/genre/handlers"
	"filmogophery/internal/app/features/health"
	mediaHandlers "filmogophery/internal/app/features/media/handlers"
	movieHandlers "filmogophery/internal/app/features/movie/handlers"
	"filmogophery/internal/app/routers"
)

func RegisterV0Routes() fx.Option {
	asV0Route := func(h any) any {
		return fx.Annotate(
			h,
			fx.As(new(routers.IRoute)),
			fx.ResultTags(`group:"v0-routers"`),
		)
	}

	return fx.Module(
		"v0-route",
		fx.Provide(
			fx.Private,
			asV0Route(health.NewHealthHandler),             // health
			asV0Route(movieHandlers.NewMockedMovieHandler), // movie
			asV0Route(genreHandlers.NewMockedGenreHandler), // genre
			asV0Route(mediaHandlers.NewMockedMediaHandler), // media

			fx.Annotate(
				getRouters,
				fx.ParamTags(`group:"v0-routers"`),
			),
		),
		fx.Invoke(func(e *echo.Echo, hs []routers.IRoute, m ...echo.MiddlewareFunc) {
			g := e.Group("v0", m...)
			for _, h := range hs {
				h.Register(g)
			}
		}),
	)
}
