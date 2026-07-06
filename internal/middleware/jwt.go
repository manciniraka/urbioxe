package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/config"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/helper"
)

func AuthMiddleware(
	cfg *config.Config,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return helper.HandleError(
					c,
					errs.ErrUnauthorized,
				)
			}

			splitToken := strings.Split(authHeader, " ")

			if len(splitToken) != 2 || splitToken[0] != "Bearer" {
				return helper.HandleError(
					c,
					errs.ErrUnauthorized,
				)
			}

			tokenString := splitToken[1]

			claims, err := helper.ParseJWT(
				tokenString,
				cfg.JWTSecret,
			)

			if err != nil {
				return helper.HandleError(
					c,
					errs.ErrUnauthorized,
				)
			}

			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
