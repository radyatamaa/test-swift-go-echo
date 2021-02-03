package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/models"
	"net/http"
)

type Usecase interface {
	CreateToken(userid string) (*models.TokenDetails, error)
	CreateAuth(userid string, td *models.TokenDetails) error
	ExtractToken(r *http.Request) string
	VerifyToken(r *http.Request) (*jwt.Token, error)
	TokenValid(r *http.Request) error
	ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error)
	FetchAuth(authD *models.AccessDetails) (string, error)
	DeleteAuth(givenUuid string) (int64,error)
	TokenAuthMiddleware() echo.HandlerFunc
}
