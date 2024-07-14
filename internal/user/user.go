package user

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"go-service/internal/user/handler"
	"go-service/internal/user/repository/adapter"
	"go-service/internal/user/service"
)

type UserTransport interface {
	All(echo.Context) error
	Load(echo.Context) error
	Create(echo.Context) error
	Update(echo.Context) error
	Patch(echo.Context) error
	Delete(echo.Context) error
}

func NewUserHandler(db *sql.DB) (UserTransport, error) {
	userRepository, err := adapter.NewUserAdapter(db)
	if err != nil {
		return nil, err
	}
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	return userHandler, nil
}
