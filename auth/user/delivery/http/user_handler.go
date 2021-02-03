package http

import (
	"context"
	"net/http"

	"github.com/auth/user"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// stockHandler  represent the http handler for country
type userHandler struct {
	userUsecase user.Usecase
}

// NewstockHandler will initialize the countrys/ resources endpoint
func NewuserHandler(e *echo.Echo, userUsecase user.Usecase) {
	handler := &userHandler{
		userUsecase: userUsecase,
	}
	e.POST("/users", handler.Create)
}

func (a *userHandler) Create(c echo.Context) error {
	var userCommand models.NewCommandUser
	if err := c.Bind(&userCommand); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")

	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	user, error := a.userUsecase.Create(ctx, userCommand)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, user)
}
func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	case models.ErrConflict:
		return http.StatusBadRequest
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
