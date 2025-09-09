package main

import (
    "log"

    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    repo := NewUserRepo()

    e.GET("/profile", GetProfile(repo))

    log.Println("user-service listening on :8084")
    e.Logger.Fatal(e.Start(":8084"))
}
