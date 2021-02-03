package usecase

import (
	"fmt"
	"github.com/auth/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/helper"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/twinj/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

type authUsecase struct {

}
// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewauthUsecase() auth.Usecase {
	return &authUsecase{

	}
}
func (a authUsecase) CreateToken(userid string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (a authUsecase) CreateAuth(userid string, td *models.TokenDetails) error {
	//at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	//rt := time.Unix(td.RtExpires, 0)
	//now := time.Now()

	//errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	//if errAccess != nil {
	//	return errAccess
	//}
	//errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	//if errRefresh != nil {
	//	return errRefresh
	//}

	helper.SetCache(td.AccessUuid,userid)
	helper.SetCache(td.RefreshUuid,userid)
	return nil
}

func (a authUsecase) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (a authUsecase) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := a.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a authUsecase) TokenValid(r *http.Request) error {
	token, err := a.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func (a authUsecase) ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := a.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId := fmt.Sprintf("%.f", claims["user_id"])
		
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:   userId,
		}, nil
	}
	return nil, err
}

func (a authUsecase) FetchAuth(authD *models.AccessDetails) (string, error) {
	userid, _ := helper.GetCache(authD.AccessUuid)
	//if err != nil {
	//	return 0, err
	//}
	//userID, _ := strconv.ParseUint(userid, 10, 64)
	return userid, nil
}

func (a authUsecase) DeleteAuth(givenUuid string) (int64, error) {
	helper.Cache.Delete(givenUuid)
	//if err != nil {
	//	return 0, err
	//}
	return 1, nil
}

func (a authUsecase) TokenAuthMiddleware() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := a.TokenValid(c.Request())
		if err != nil {
			return	c.JSON(http.StatusUnauthorized, err.Error())
			//c.Abort()

		}
		return nil
		//c.Next()
	}
}


