package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/auth/auth"

	"github.com/labstack/echo"
	"github.com/master/inbound"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// inboundHandler  represent the http handler for country
type inboundHandler struct {
	inboundUsecase inbound.Usecase
	authUsecase    auth.Usecase
}

// NewinboundHandler will initialize the countrys/ resources endpoint
func NewinboundHandler(e *echo.Echo, us inbound.Usecase, authUsecase auth.Usecase) {
	handler := &inboundHandler{
		inboundUsecase: us,
	}
	e.POST("/master/inbound", handler.Create)
	e.PUT("/master/inbound/:id", handler.UpdateInboundy)
	e.DELETE("/master/inbound/:id", handler.Delete)
	// e.GET("/countrys/:id/credit", handler.GetCreditByID)
	e.GET("/master/inbound/:id", handler.GetDetailID)
	e.GET("/master/inbound", handler.List)
}

func isRequestValid(m *models.NewCommandInbound) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *inboundHandler) Delete(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	id := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	result, err := a.inboundUsecase.Delete(ctx, id, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *inboundHandler) List(c echo.Context) error {

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	search := c.QueryParam("search")

	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qperPage)
	offset = (page - 1) * limit

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.inboundUsecase.List(ctx, page, limit, offset, search)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *inboundHandler) UpdateInboundy(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	id := c.Param("id")
	var inbound models.NewCommandInbound
	inbound.Id = id
	err = c.Bind(&inbound)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.inboundUsecase.Update(ctx, &inbound, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, inbound)
}

func (a *inboundHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	var inbound models.NewCommandInbound
	err = c.Bind(&inbound)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.inboundUsecase.Create(ctx, &inbound, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *inboundHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.inboundUsecase.GetById(ctx, id, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
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
