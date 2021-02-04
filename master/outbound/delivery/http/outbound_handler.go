package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/auth/auth"

	"github.com/labstack/echo"
	"github.com/master/outbound"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// outboundHandler  represent the http handler for country
type outboundHandler struct {
	outboundUsecase outbound.Usecase
	authUsecase     auth.Usecase
}

// NewoutboundHandler will initialize the countrys/ resources endpoint
func NewoutboundHandler(e *echo.Echo, us outbound.Usecase, authUsecase auth.Usecase) {
	handler := &outboundHandler{
		outboundUsecase: us,

		authUsecase:authUsecase,
	}
	//e.POST("/master/outbound", handler.Create)
	//e.PUT("/master/outbound/:id", handler.UpdateOutboundy)
	//e.DELETE("/master/outbound/:id", handler.Delete)
	// e.GET("/countrys/:id/credit", handler.GetCreditByID)
	//e.GET("/master/outbound/:id", handler.GetDetailID)
	e.GET("/master/outbound", handler.List)
}

func isRequestValid(m *models.NewCommandOutbound) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
//func (a *outboundHandler) Delete(c echo.Context) error {
//	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	userId, err := a.authUsecase.FetchAuth(tokenAuth)
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	id := c.Param("id")
//	ctx := c.Request().Context()
//	if ctx == nil {
//		ctx = context.Background()
//	}
//	result, err := a.outboundUsecase.Delete(ctx, id, userId)
//	if err != nil {
//		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
//	}
//	return c.JSON(http.StatusOK, result)
//}

func (a *outboundHandler) List(c echo.Context) error {

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

	result, err := a.outboundUsecase.List(ctx, page, limit, offset, search)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

//func (a *outboundHandler) UpdateOutboundy(c echo.Context) error {
//	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	userId, err := a.authUsecase.FetchAuth(tokenAuth)
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	id := c.Param("id")
//	var outbound models.NewCommandOutbound
//	outbound.Id = id
//	err = c.Bind(&outbound)
//	if err != nil {
//		return c.JSON(http.StatusUnprocessableEntity, err.Error())
//	}
//
//	ctx := c.Request().Context()
//	if ctx == nil {
//		ctx = context.Background()
//	}
//
//	err = a.outboundUsecase.Update(ctx, &outbound, userId)
//	if err != nil {
//		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
//	}
//	return c.JSON(http.StatusOK, outbound)
//}
//
//func (a *outboundHandler) Create(c echo.Context) error {
//	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	userId, err := a.authUsecase.FetchAuth(tokenAuth)
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	var outbound models.NewCommandOutbound
//	err = c.Bind(&outbound)
//	if err != nil {
//		return c.JSON(http.StatusUnprocessableEntity, err.Error())
//	}
//
//	ctx := c.Request().Context()
//	if ctx == nil {
//		ctx = context.Background()
//	}
//
//	res, err := a.outboundUsecase.Create(ctx, &outbound, userId)
//	if err != nil {
//		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
//	}
//	return c.JSON(http.StatusOK, res)
//}
//
//func (a *outboundHandler) GetDetailID(c echo.Context) error {
//	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	userId, err := a.authUsecase.FetchAuth(tokenAuth)
//	if err != nil {
//		return c.JSON(http.StatusUnauthorized, "unauthorized")
//
//	}
//	id := c.Param("id")
//
//	ctx := c.Request().Context()
//	if ctx == nil {
//		ctx = context.Background()
//	}
//
//	result, err := a.outboundUsecase.GetById(ctx, id, userId)
//	if err != nil {
//		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
//	}
//	return c.JSON(http.StatusOK, result)
//}

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
