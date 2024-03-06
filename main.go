package main

import (
	handlerAuth "dietku-backend/cmd/auth/handler"
	"dietku-backend/cmd/log"
	handlerUser "dietku-backend/cmd/user/handler"
	"dietku-backend/config"
	"dietku-backend/docs"
	"dietku-backend/version"
	"fmt"
	"github.com/labstack/echo/v4"
	echoswagger "github.com/swaggo/echo-swagger"
	"net/http"
	"time"
)

// @title Dietku Backend API
// @description Dietku Backend API
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	defer log.RecoverWithTrace()

	http.DefaultClient.Timeout = 30 * time.Second

	conf, err := config.InitConfigApp(".env")
	if err != nil {
		fmt.Println(err)
	}

	db := config.ConnectMongo()

	e := echo.New()
	log.SetLogger(e)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	docs.SwaggerInfo.Version = version.Version
	docs.SwaggerInfo.Host = conf.SwaggerHost
	e.GET("/swagger/*", echoswagger.WrapHandler)

	handlerAuth.NewAuthHandler(e, db)
	handlerUser.NewUserApi(e, db)

	server := fmt.Sprintf("%v:3000", conf.AppHost)
	if conf.AppPort != "" {
		server = fmt.Sprintf("%v:%v", conf.AppHost, conf.AppPort)
	}

	log.Fatal(e.Start(server))
}
