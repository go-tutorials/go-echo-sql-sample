package app

import (
	"context"
	"github.com/labstack/echo/v4"
)

func Route(ctx context.Context, e *echo.Echo, config Config) error {
	app, err := NewApp(ctx, config)
	if err != nil {
		return err
	}

	e.GET("/health", app.Health.Check)

	userPath := "/users"
	e.GET(userPath, app.Handler.All)
	e.GET(userPath+"/:id", app.Handler.Load)
	e.POST(userPath, app.Handler.Create)
	e.PUT(userPath+"/:id", app.Handler.Update)
	e.PATCH(userPath+"/:id", app.Handler.Patch)
	e.DELETE(userPath+"/:id", app.Handler.Delete)

	return nil
}
