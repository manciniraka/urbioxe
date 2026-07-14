package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type NewsController struct {
	newsService service.NewsService
}

func NewNewsController(
	newsService service.NewsService,
) *NewsController {
	return &NewsController{
		newsService: newsService,
	}
}

func (nc *NewsController) GetAll(c echo.Context) error {
	news, err := nc.newsService.GetAll()
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"news retrieved successfully",
		news,
	)
}

func (nc *NewsController) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid news id",
		)
	}

	news, err := nc.newsService.GetByID(id)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"news retrieved successfully",
		news,
	)
}

func (nc *NewsController) Create(c echo.Context) error {
	userID := helper.GetUserID(c)

	var input service.CreateNewsInput

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

	news, err := nc.newsService.Create(
		int64(userID),
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
		"news created successfully",
		news,
	)
}

func (nc *NewsController) Update(c echo.Context) error {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid news id",
		)
	}

	var input service.UpdateNewsInput

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

	news, err := nc.newsService.Update(
		id,
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
		"news updated successfully",
		news,
	)
}

func (nc *NewsController) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		return helper.BadRequest(
			c,
			"invalid news id",
		)
	}

	err = nc.newsService.Delete(id)
	if err != nil {
		return helper.HandleError(
			c,
			err,
		)
	}

	return helper.Success(
		c,
		"news deleted successfully",
		nil,
	)
}
