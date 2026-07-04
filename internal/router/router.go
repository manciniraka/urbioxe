package router

import "github.com/labstack/echo/v4"

func InitRouter(
	e *echo.Echo,
) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(
			200,
			echo.Map{
				"message": "ubioxe service is running",
			},
		)
	})
}
