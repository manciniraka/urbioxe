package controller

import (
	"github.com/labstack/echo/v4"

	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(
	authService service.AuthService,
) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (ac *AuthController) Register(c echo.Context) error {
	var input service.RegisterInput

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

	user, err := ac.authService.Register(input)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Created(
		c,
		"user registered successfully",
		user,
	)
}

func (ac *AuthController) Login(
	c echo.Context,
) error {

	var input service.LoginInput

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

	token, err := ac.authService.Login(input)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"login successfully",
		LoginResponse{
			Token: token,
		},
	)
}
