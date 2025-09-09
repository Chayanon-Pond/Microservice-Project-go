package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(c echo.Context, db *sql.DB) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	if err := RegisterUser(db, req.Email, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, "Registered")
}

func LoginHandler(c echo.Context, db *sql.DB) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	user, err := LoginUser(db, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid credentials")
	}
	token, _ := GenerateToken(user.Email)
	refresh, _ := GenerateRefreshToken(user.Email)
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  token,
		"refresh_token": refresh,
	})
}

func ValidateTokenHandler(c echo.Context) error {
	tokenStr := c.QueryParam("token")
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid token")
	}
	return c.JSON(http.StatusOK, claims)
}

func GenerateToken(email string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"email": email,
		"exp":   jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 นาที
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(email string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"email": email,
		"exp":   jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 วัน
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
