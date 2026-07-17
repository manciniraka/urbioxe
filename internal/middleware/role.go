package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/errs"
	"github.com/manciniraka/urbioxe/internal/helper"
)

func RequireRoles(
	roles ...entity.UserRole,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := helper.GetUserRole(c)

			for _, role := range roles {

				if userRole == string(role) {
					return next(c)
				}
			}

			return helper.HandleError(
				c,
				errs.ErrForbidden,
			)
		}
	}
}
