package http

import (
	"context"
	"github.com/auth/auth"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/master/product"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// productHandler  represent the http handler for country
type productHandler struct {
	productUsecase product.Usecase
	authUsecase auth.Usecase
}

// NewproductHandler will initialize the countrys/ resources endpoint
func NewproductHandler(e *echo.Echo, us product.Usecase,authUsecase auth.Usecase) {
	handler := &productHandler{
		productUsecase: us,
	}
	e.POST("/master/product", handler.Create)
	e.PUT("/master/product/:id", handler.UpdateProducty)
	e.DELETE("/master/product/:id", handler.Delete)
	// e.GET("/countrys/:id/credit", handler.GetCreditByID)
	e.GET("/master/product/:id", handler.GetDetailID)
	e.GET("/master/product", handler.List)
}

func isRequestValid(m *models.NewCommandProduct) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *productHandler) Delete(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	id := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	result, err := a.productUsecase.Delete(ctx, id, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *productHandler) List(c echo.Context) error {

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

	result, err := a.productUsecase.List(ctx, page, limit, offset, search)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *productHandler) UpdateProducty(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	id := c.Param("id")
	var product models.NewCommandProduct
	product.Id = id
	err = c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.productUsecase.Update(ctx, &product, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, product)
}

func (a *productHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	var product models.NewCommandProduct
	err = c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.productUsecase.Create(ctx, &product, userId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *productHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return	c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.productUsecase.GetById(ctx, id, userId)
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
