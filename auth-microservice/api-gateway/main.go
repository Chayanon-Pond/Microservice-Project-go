package main

import (
	"fmt"
	"io"
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
		// Allow OPTIONS preflight and public paths (e.g., auth service endpoints)
		if c.Request().Method == http.MethodOptions {
			return next(c)
		}
		publicPrefixes := []string{"/auth/login", "/auth/register"}
		for _, p := range publicPrefixes {
			if strings.HasPrefix(c.Request().URL.Path, p) {
				return next(c)
			}
		}

		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, "Missing or invalid token")
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(os.Getenv("JWT_SECRET"))

		// Parse with claims and restrict accepted signing methods
		token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			// ensure signing method is HMAC-based (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid token")
		}

		// attach claims to context for downstream handlers
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
			if email, ok := claims["email"].(string); ok {
				// forward email to downstream services as header
				c.Request().Header.Set("X-User-Email", email)
			}
		}
		return next(c)
	}
}

// ProxyToService เป็นตัวอย่าง proxy (forward) request ไป service อื่น
func ProxyToService(target string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// build proxied request (preserve query string)
		req, err := http.NewRequest(c.Request().Method, target+c.Request().URL.Path, c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Proxy error")
		}
		req.Header = c.Request().Header.Clone()
		req.URL.RawQuery = c.Request().URL.RawQuery
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.JSON(http.StatusBadGateway, "Service unavailable")
		}
		defer resp.Body.Close()
		// copy response headers
		for k, vv := range resp.Header {
			c.Response().Header().Del(k)
			for _, v := range vv {
				c.Response().Header().Add(k, v)
			}
		}
		c.Response().WriteHeader(resp.StatusCode)
		_, copyErr := io.Copy(c.Response().Writer, resp.Body)
		if copyErr != nil {
			return c.JSON(http.StatusInternalServerError, "Error copying response body")
		}
		return nil
	}
}
