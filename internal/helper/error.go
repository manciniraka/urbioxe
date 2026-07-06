package helper

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/errs"
)

var logger = slog.Default()

func HandleError(
	c echo.Context,
	err error,
) error {

	status := errs.StatusCode(err)

	if status == http.StatusInternalServerError {

		logger.Error(
			"request failed",
			"path", c.Path(),
			"method", c.Request().Method,
			"error", err,
		)

		return c.JSON(
			status,
			echo.Map{
				"message": "internal server error",
			},
		)

	}

	logger.Warn(
		"request failed",
		"path", c.Path(),
		"method", c.Request().Method,
		"ip", c.RealIP(),
		"status", status,
		"error", err,
	)

	return c.JSON(
		status,
		echo.Map{
			"message": err.Error(),
		},
	)

}
