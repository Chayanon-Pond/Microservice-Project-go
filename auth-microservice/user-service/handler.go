package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

// GetProfile reads the forwarded X-User-Email header and returns the user profile.
func GetProfile(repo *UserRepo) echo.HandlerFunc {
    return func(c echo.Context) error {
        email := c.Request().Header.Get("X-User-Email")
        if email == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing user header"})
        }
        user, err := repo.GetByEmail(email)
        if err != nil {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
        }
        return c.JSON(http.StatusOK, user)
    }
}
