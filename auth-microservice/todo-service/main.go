package main

import "github.com/labstack/echo/v4"

func main() {
    e := echo.New()
    repo := NewTodoRepo()

    e.GET("/todos", ListTodosHandler(repo))
    e.POST("/todos", AddTodoHandler(repo))
    e.DELETE("/todos/:id", DeleteTodoHandler(repo))

    e.Logger.Fatal(e.Start(":8083"))
}
