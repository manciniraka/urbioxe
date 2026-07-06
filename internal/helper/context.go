package helper

import "github.com/labstack/echo/v4"

func GetUserID(c echo.Context) string {

	userID, _ := c.Get("user_id").(string)

	return userID

}

func GetUserRole(c echo.Context) string {

	role, _ := c.Get("role").(string)

	return role

}
