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

// GetAll godoc
//
//	@Summary		Get all news
//	@Description	Retrieve all regional news
//	@Tags			News
//	@Produce		json
//	@Success		200	{object}	helper.Response
//	@Failure		500	{object}	helper.ErrorResponse
//	@Router			/news [get]
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

// GetByID godoc
//
//	@Summary		Get news by ID
//	@Description	Retrieve news by ID
//	@Tags			News
//	@Produce		json
//	@Param			id	path		int	true	"News ID"
//	@Success		200	{object}	helper.Response
//	@Failure		400	{object}	helper.ErrorResponse
//	@Failure		404	{object}	helper.ErrorResponse
//	@Router			/news/{id} [get]
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

// Create godoc
//
//	@Summary		Create news
//	@Description	Create a new news
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Param			news	body		service.CreateNewsInput	true	"News data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Router			/news [post]
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

// Update godoc
//
//	@Summary		Update news
//	@Description	Update an existing news
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"News ID"
//	@Param			news	body		service.UpdateNewsInput	true	"News data"
//	@Success		200		{object}	helper.Response
//	@Failure		400		{object}	helper.ErrorResponse
//	@Failure		404		{object}	helper.ErrorResponse
//	@Router			/news/{id} [put]
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

// Delete godoc
//
//	@Summary		Delete news
//	@Description	Delete an existing news
//	@Tags			News
//	@Produce		json
//	@Param			id	path		int	true	"News ID"
//	@Success		200	{object}	helper.Response
//	@Failure		400	{object}	helper.ErrorResponse
//	@Failure		404	{object}	helper.ErrorResponse
//	@Router			/news/{id} [delete]
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
