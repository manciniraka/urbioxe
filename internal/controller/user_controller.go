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

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get authenticated user profile
//	@Tags			User
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	helper.Response
//	@Failure		401	{object}	helper.ErrorResponse
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/users/profile [get]
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

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Update authenticated user profile
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		service.UpdateProfileInput	true	"Profile data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/users/profile [put]
func (uc *UserController) UpdateProfile(c echo.Context) error {
	userID := helper.GetUserID(c)

	var input service.UpdateProfileInput

	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(
			c,
			"invalid request body",
		)
	}

	if err := c.Validate(&input); err != nil {
		return helper.BadRequest(
			c,
			err.Error(),
		)
	}

	user, err := uc.userService.UpdateProfile(
		userID,
		input,
	)

	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"profile updated successfully",
		user,
	)
}

// ChangePassword godoc
//
//	@Summary		Change password
//	@Description	Change authenticated user password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		service.ChangePasswordInput	true	"Password data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		401		{object}	helper.ErrorResponse
//	@Failure		500		{object}	helper.ErrorResponse
//	@Router			/users/password [put]
func (uc *UserController) ChangePassword(c echo.Context) error {
	userID := helper.GetUserID(c)

	var input service.ChangePasswordInput

	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(
			c,
			"invalid request body",
		)
	}

	if err := c.Validate(&input); err != nil {
		return helper.BadRequest(
			c,
			err.Error(),
		)
	}

	err := uc.userService.ChangePassword(
		userID,
		input,
	)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"password changed successfully",
		nil,
	)
}
