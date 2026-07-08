package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(
	userService service.UserService,
) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) GetProfile(c echo.Context) error {
	userID := helper.GetUserID(c)

	user, err := uc.userService.GetProfile(userID)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"profile retrieved successfully",
		user,
	)
}
