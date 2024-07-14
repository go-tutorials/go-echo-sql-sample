package handler

import (
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/labstack/echo/v4"

	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) All(c echo.Context) error {
	res, err := h.service.All(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Load(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	res, err := h.service.Load(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if res == nil {
		return c.JSON(http.StatusNotFound, res)
	}
	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Create(c echo.Context) error {
	var user model.User
	er1 := c.Bind(&user)

	defer c.Request().Body.Close()
	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	res, er2 := h.service.Create(c.Request().Context(), &user)
	if er2 != nil {
		return c.String(http.StatusInternalServerError, er2.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusCreated, user)
	} else {
		return c.JSON(http.StatusConflict, res)
	}

}

func (h *UserHandler) Update(c echo.Context) error {
	var user model.User
	er1 := c.Bind(&user)
	defer c.Request().Body.Close()

	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		return c.String(http.StatusBadRequest, "Id not match")
	}

	res, er2 := h.service.Update(c.Request().Context(), &user)
	if er2 != nil {
		return c.String(http.StatusInternalServerError, er2.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusOK, user)
	} else if res == 0 {
		return c.JSON(http.StatusNotFound, res)
	} else {
		return c.JSON(http.StatusConflict, res)
	}
}

func (h *UserHandler) Patch(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	r := c.Request()
	var user model.User
	userType := reflect.TypeOf(user)
	_, jsonMap, _ := core.BuildMapField(userType)
	body, er0 := core.BuildMapAndStruct(r, &user)
	if er0 != nil {
		return c.String(http.StatusInternalServerError, er0.Error())
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		return c.String(http.StatusBadRequest, "Id not match")
	}
	json, er1 := core.BodyToJsonMap(r, user, body, []string{"id"}, jsonMap)
	if er1 != nil {
		return c.String(http.StatusInternalServerError, er1.Error())
	}

	res, er2 := h.service.Patch(r.Context(), json)
	if er2 != nil {
		return c.String(http.StatusInternalServerError, er2.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusOK, json)
	} else if res == 0 {
		return c.JSON(http.StatusNotFound, res)
	} else {
		return c.JSON(http.StatusConflict, res)
	}
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "Id cannot be empty")
	}

	res, err := h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if res > 0 {
		return c.JSON(http.StatusOK, res)
	} else {
		return c.JSON(http.StatusNotFound, res)
	}
}
