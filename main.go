package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"


	_authHttpDeliver "github.com/auth/auth/delivery/http"
	_authUsecase "github.com/auth/auth/usecase"

	_userHttpDeliver "github.com/auth/user/delivery/http"
	_userRepo "github.com/auth/user/repository"
	_userUsecase "github.com/auth/user/usecase"

	_inboundHttpDeliver "github.com/master/inbound/delivery/http"
	_inboundRepo "github.com/master/inbound/repository"
	_inboundUsecase "github.com/master/inbound/usecase"

	_outboundHttpDeliver "github.com/master/outbound/delivery/http"
	_outboundRepo "github.com/master/outbound/repository"
	_outboundUsecase "github.com/master/outbound/usecase"

	_productHttpDeliver "github.com/master/product/delivery/http"
	_productRepo "github.com/master/product/repository"
	_productUsecase "github.com/master/product/usecase"

	_stockHttpDeliver "github.com/master/stock/delivery/http"
	_stockRepo "github.com/master/stock/repository"
	_stockUsecase "github.com/master/stock/usecase"

	_orderHttpDeliver "github.com/order/order/delivery/http"
	_orderRepo "github.com/order/order/repository"
	_orderUsecase "github.com/order/order/usecase"

	_echoMiddleware "github.com/labstack/echo/middleware"


)

func main() {

	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)


	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	//middL := middleware.InitMiddleware()
	//e.Use(middL.CORS)
	e.Use(_echoMiddleware.CORSWithConfig(_echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	userRepo := _userRepo.NewUserRepository(dbConn)
	inboundRepo := _inboundRepo.NewInboundRepository(dbConn)
	outboundRepo := _outboundRepo.NewOutboundRepository(dbConn)
	productRepo := _productRepo.NewProductRepository(dbConn)
	stockRepo := _stockRepo.NewStockRepository(dbConn)
	orderRepo := _orderRepo.NewOrderRepository(dbConn)

	timeoutContext := 30 * time.Second

	authUsecase := _authUsecase.NewauthUsecase()
	userUsecase := _userUsecase.NewuserUsecase(timeoutContext,userRepo)
	inboundUsecase := _inboundUsecase.NewinboundUsecase(userUsecase,inboundRepo,timeoutContext)
	outboundUsecase := _outboundUsecase.NewoutboundUsecase(userUsecase,outboundRepo,timeoutContext)
	productUsecase := _productUsecase.NewproductUsecase(userUsecase,productRepo,timeoutContext)
	stockUsecase := _stockUsecase.NewstockUsecase(userUsecase,stockRepo,timeoutContext)
	orderUsecase := _orderUsecase.NeworderUsecase(userUsecase,orderRepo,timeoutContext,outboundRepo,stockRepo)

	_authHttpDeliver.NewuserHandler(e,authUsecase,userUsecase)
	_userHttpDeliver.NewuserHandler(e,userUsecase)
	_inboundHttpDeliver.NewinboundHandler(e,inboundUsecase,authUsecase)
	_outboundHttpDeliver.NewoutboundHandler(e,outboundUsecase,authUsecase)
	_productHttpDeliver.NewproductHandler(e,productUsecase,authUsecase)
	_stockHttpDeliver.NewstockHandler(e,stockUsecase,authUsecase)
	_orderHttpDeliver.NeworderHandler(e,orderUsecase,authUsecase)

	log.Fatal(e.Start(":9090"))
}
