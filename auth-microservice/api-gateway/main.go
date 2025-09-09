package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Use(JWTMiddleware)

	// Proxy route ตัวอย่าง (forward ไป todo-service)
	e.Any("/todo/*", ProxyToService("http://localhost:8083"))

	e.Logger.Fatal(e.Start(":8080"))
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, "Missing or invalid token")
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid token")
		}
		return next(c)
	}
}

// ProxyToService เป็นตัวอย่าง proxy (forward) request ไป service อื่น
func ProxyToService(target string) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, err := http.NewRequest(c.Request().Method, target+c.Request().URL.Path, c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Proxy error")
		}
		req.Header = c.Request().Header
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.JSON(http.StatusBadGateway, "Service unavailable")
		}
		defer resp.Body.Close()
		return c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
	}
}
