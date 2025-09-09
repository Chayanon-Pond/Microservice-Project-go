package main

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo/v4"
)

type addTodoReq struct {
    Text string `json:"text"`
}

func ListTodosHandler(repo *TodoRepo) echo.HandlerFunc {
    return func(c echo.Context) error {
        owner := c.Request().Header.Get("X-User-Email")
        if owner == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing user"})
        }
        todos := repo.ListByOwner(owner)
        return c.JSON(http.StatusOK, todos)
    }
}

func AddTodoHandler(repo *TodoRepo) echo.HandlerFunc {
    return func(c echo.Context) error {
        owner := c.Request().Header.Get("X-User-Email")
        if owner == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing user"})
        }
        var req addTodoReq
        if err := c.Bind(&req); err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
        }
        t := repo.Add(owner, req.Text)
        return c.JSON(http.StatusCreated, t)
    }
}

func DeleteTodoHandler(repo *TodoRepo) echo.HandlerFunc {
    return func(c echo.Context) error {
        owner := c.Request().Header.Get("X-User-Email")
        if owner == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing user"})
        }
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
        }
        ok := repo.Delete(owner, id)
        if !ok {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "not found or not allowed"})
        }
        return c.NoContent(http.StatusNoContent)
    }
}
