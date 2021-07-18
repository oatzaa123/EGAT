package main

import (
	//	"go-callcenter/controllers/announcement"
	//	"go-callcenter/controllers/division"
	//	"go-callcenter/controllers/group"
	//	"go-callcenter/controllers/setting"
	//	"go-callcenter/controllers/worktype"

	"go-callcenter/controllers/charger"
	"go-callcenter/controllers/stations"
	"go-callcenter/controllers/users"
	"go-callcenter/database"
	"go-callcenter/keys"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	database.InitDatabase()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowOrigin},
	}))

	// privateRoute := e.Group("/private")
	publicRoute := e.Group("/public")
	privateRoute := e.Group("/private")

	// middleware
	privateRoute.Use(middleware.JWT([]byte(keys.Secret)))

	// public route
	users.PublicRoute(publicRoute)

	//private route
	//group.PrivateRoute(privateRoute)
	//announcement.PrivateRoute(privateRoute)
	//setting.PrivateRoute(privateRoute)
	//division.PrivateRoute(privateRoute)
	charger.PrivateRoute(privateRoute)
	users.PrivateRoute(privateRoute)
	stations.PrivateRoute(privateRoute)
	//worktype.PrivateRoute(privateRoute)

	e.Logger.Fatal(e.Start(":3100"))

}
