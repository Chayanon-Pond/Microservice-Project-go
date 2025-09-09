package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	db := initDB()
	defer db.Close()

	e := echo.New()
	e.POST("/register", func(c echo.Context) error { return RegisterHandler(c, db) })
	e.POST("/login", func(c echo.Context) error { return LoginHandler(c, db) })
	e.GET("/validate", ValidateTokenHandler)

	e.Logger.Fatal(e.Start(":8081"))
}
