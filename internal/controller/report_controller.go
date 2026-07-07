package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manciniraka/urbioxe/internal/entity"
	"github.com/manciniraka/urbioxe/internal/helper"
	"github.com/manciniraka/urbioxe/internal/service"
)

type ReportController struct {
	svc service.ReportService
}

func NewReportController(
	reportService service.ReportService,
) *ReportController {
	return &ReportController{
		svc: reportService,
	}
}

func (rc *ReportController) Create(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return helper.BadRequest(c, "error reading form data")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	districtID, _ := strconv.ParseInt(c.FormValue("incident_district_id"), 10, 64)

	userID := int64(1) // helper.GetUserID(c)

	files := form.File["images"]
	if len(files) == 0 {
		return helper.BadRequest(c, "upload minimal 1 image")
	}

	reportInput := entity.Report{
		UserID:             userID,
		CategoryID:         categoryID,
		IncidentDistrictID: districtID,
		Title:              title,
		Description:        description,
	}

	result, err := rc.svc.CreateReport(reportInput, files)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Created(
		c,
		"success create report",
		result,
	)
}
