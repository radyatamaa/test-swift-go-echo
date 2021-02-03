package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/auth/auth"
	_userUsecase "github.com/auth/user"
	"github.com/dgrijalva/jwt-go"
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
	userUsecase _userUsecase.Usecase
	authUsecase auth.Usecase
}

// NewstockHandler will initialize the countrys/ resources endpoint
func NewuserHandler(e *echo.Echo, authUsecase auth.Usecase, userUsecase _userUsecase.Usecase) {
	handler := &userHandler{
		authUsecase: authUsecase,
		userUsecase: userUsecase,
	}
	e.POST("/auth/login", handler.Login)
	e.POST("/auth/logout", handler.Logout)
	e.POST("/auth/get-user-info", handler.CreateTodo)
	e.POST("/auth/token/refresh", handler.Refresh)
}

//var user = models.UserAuth{
//	ID:            1,
//	Username: "username",
//	Password: "password",
//	Phone: "49123454322", //this is a random number
//}

func (a *userHandler) Login(c echo.Context) error {
	var u models.UserAuth
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")

	}
	//compare the user from the request, with the one we defined:
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	user, err := a.userUsecase.ValidateUser(ctx, u.Email, u.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Please provide valid login details")

	}
	ts, err := a.authUsecase.CreateToken(user.Id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())

	}
	saveErr := a.authUsecase.CreateAuth(user.Id, ts)
	if saveErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return c.JSON(http.StatusOK, tokens)
}
func (a *userHandler) CreateTodo(c echo.Context) error {
	a.authUsecase.TokenAuthMiddleware()
	var td *models.Todo
	if err := c.Bind(&td); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "invalid json")

	}
	tokenAuth, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	userId, err := a.authUsecase.FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	td.UserID = userId

	//you can proceed to save the Todo to a database
	//but we will just return it to the caller here:
	return c.JSON(http.StatusCreated, td)
}
func (a *userHandler) Logout(c echo.Context) error {
	a.authUsecase.TokenAuthMiddleware()
	au, err := a.authUsecase.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	deleted, delErr := a.authUsecase.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		return c.JSON(http.StatusUnauthorized, "unauthorized")

	}
	return c.JSON(http.StatusOK, "Successfully logged out")
}
func (a *userHandler) Refresh(c echo.Context) error {
	mapToken := map[string]string{}
	if err := c.Bind(&mapToken); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())

	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Refresh token expired")

	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return c.JSON(http.StatusUnauthorized, err)

	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, err)

		}
		userId := fmt.Sprintf("%.f", claims["user_id"])

		//Delete the previous Refresh Token
		deleted, delErr := a.authUsecase.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			return c.JSON(http.StatusUnauthorized, "unauthorized")

		}
		//Create new pairs of refresh and access tokens
		ts, createErr := a.authUsecase.CreateToken(userId)
		if createErr != nil {
			return c.JSON(http.StatusForbidden, createErr.Error())

		}
		//save the tokens metadata to redis
		saveErr := a.authUsecase.CreateAuth(userId, ts)
		if saveErr != nil {
			return c.JSON(http.StatusForbidden, saveErr.Error())

		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return c.JSON(http.StatusCreated, tokens)
	} else {
		return c.JSON(http.StatusUnauthorized, "refresh expired")
	}
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
