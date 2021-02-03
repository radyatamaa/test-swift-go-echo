package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/auth/auth"

	"github.com/labstack/echo"
	"github.com/master/stock"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// stockHandler  represent the http handler for country
type stockHandler struct {
	stockUsecase stock.Usecase
	authUsecase  auth.Usecase
}

// NewstockHandler will initialize the countrys/ resources endpoint
func NewstockHandler(e *echo.Echo, us stock.Usecase, authUsecase auth.Usecase) {
	handler := &stockHandler{
		stockUsecase: us,
	}
	e.GET("/master/stock", handler.List)
}

//func isRequestValid(m *models.NewCommandStock) (bool, error) {
//	validate := validator.New()
//	err := validate.Struct(m)
//	if err != nil {
//		return false, err
//	}
//	return true, nil
//}

func (a *stockHandler) List(c echo.Context) error {

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	productId := c.QueryParam("product_id")
	types := c.QueryParam("types")

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
	bound := 0
	if types == "inbound"{
		bound = 1
	}else if types == "outbound"{
		bound = 2
	}
	result, err := a.stockUsecase.List(ctx, page, limit, offset, productId,bound)
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
