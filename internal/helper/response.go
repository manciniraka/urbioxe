package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Success(
	c echo.Context,
	message string,
	data any,
) error {

	response := echo.Map{
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	return c.JSON(
		http.StatusOK,
		response,
	)
}

func Created(
	c echo.Context,
	message string,
	data any,
) error {

	response := echo.Map{
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	return c.JSON(
		http.StatusCreated,
		response,
	)
}

func BadRequest(
	c echo.Context,
	message string,
) error {

	return c.JSON(
		http.StatusBadRequest,
		echo.Map{
			"message": message,
		},
	)
}
