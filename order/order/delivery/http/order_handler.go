package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/auth/auth"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/order/order"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// orderHandler  represent the http handler for country
type orderHandler struct {
	orderUsecase order.Usecase
	authUsecase  auth.Usecase
}

// NeworderHandler will initialize the countrys/ resources endpoint
func NeworderHandler(e *echo.Echo, us order.Usecase, authUsecase auth.Usecase) {
	handler := &orderHandler{
		orderUsecase: us,
		authUsecase:authUsecase,
	}
	e.POST("/order/order", handler.Create)
	e.PUT("/order/order-status/:id", handler.UpdateStatus)
	e.GET("/order/order", handler.List)
}

func isRequestValid(m *models.NewCommandOrder) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *orderHandler) List(c echo.Context) error {

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	referenceNumber := c.QueryParam("reference_number")

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

	result, err := a.orderUsecase.List(ctx, page, limit, offset, referenceNumber)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *orderHandler) UpdateStatus(c echo.Context) error {

	id := c.Param("id")
	var order models.NewCommandOrderStatus
	order.ReferenceNumber = id
	err := c.Bind(&order)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.orderUsecase.UpdateStatus(ctx, &order)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	result := models.ResponseDelete{
		Id:      order.ReferenceNumber,
		Message: "Success",
	}
	return c.JSON(http.StatusOK, result)
}

func (a *orderHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	var order models.NewCommandOrder
	err = c.Bind(&order)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.orderUsecase.Create(ctx, &order, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
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
