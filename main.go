package main

import (
	handlerAuth "dietku-backend/cmd/auth/handler"
	handlerBlog "dietku-backend/cmd/blog/handler"
	"dietku-backend/cmd/log"
	handlerUser "dietku-backend/cmd/user/handler"
	"dietku-backend/config"
	"dietku-backend/docs"
	"dietku-backend/version"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"
	"net/http"
	"strings"
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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(conf.AllowOrigins, ","),
		AllowCredentials: true,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/google", func(c echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(`
<!doctype html>
<html>
<head>
    <title>Google SignIn</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css"> <!-- load bulma css -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"> <!-- load fontawesome -->
    <style>
        body{ padding-top:70px; }
    </style>
</head>
<body>
<div class="container">
    <div class="jumbotron text-center text-success">
        <h1><span class="fa fa-lock"></span> Social Authentication</h1>
        <p>Login or Register with:</p>
        <a href="https://dietku-api.up.railway.app/api/login-google" class="btn btn-danger"><span class="fa fa-google"></span> Sign In with Google</a>
    </div>
</div>
</body>
</html>`))
	})

	docs.SwaggerInfo.Version = version.Version
	docs.SwaggerInfo.Host = conf.SwaggerHost
	e.GET("/swagger/*", echoswagger.WrapHandler)

	handlerAuth.NewAuthHandler(e, db, conf)
	handlerUser.NewUserApi(e, db)
	handlerBlog.NewBlogApi(e, db)

	server := fmt.Sprintf("%v:3000", conf.AppHost)
	if conf.AppPort != "" {
		server = fmt.Sprintf("%v:%v", conf.AppHost, conf.AppPort)
	}

	log.Fatal(e.Start(server))
}
