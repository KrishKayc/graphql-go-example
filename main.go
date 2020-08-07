package main

import (
	"bookCab/graphql"

	_ "bookCab/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Book Cab API
// @version 1.0
// @description Server to book cabs.

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.basic BasicAuth
func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//initialize persitance and cache from config/secrets manager
	db, cache, err := config("config.json")
	handleErr(e, err)

	defer db.Close()
	defer cache.Close()

	//graphql handler creates fields and the queries, mutations available in those fields
	h, err := graphql.NewHandler(db, cache)
	handleErr(e, err)

	e.POST("/graphql", echo.WrapHandler(h))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := e.Start(":3000"); err != nil {
		e.Logger.Fatal(err)
	}
}

func handleErr(e *echo.Echo, err error) {
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
}
