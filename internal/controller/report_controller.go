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
	categoryID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 64)
	districtID, _ := strconv.ParseUint(c.FormValue("incident_district_id"), 10, 64)

	userID := helper.GetUserID(c)

	files := form.File["images"]
	if len(files) == 0 {
		return helper.BadRequest(c, "upload minimal 1 image")
	}

	reportInput := entity.Report{
		UserID:             userID,
		CategoryID:         uint(categoryID),
		IncidentDistrictID: uint(districtID),
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

func (rc *ReportController) GetAll(c echo.Context) error {
	userID := helper.GetUserID(c)
	role := helper.GetUserRole(c)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	districtID, _ := strconv.ParseUint(c.QueryParam("district_id"), 10, 64)
	status := c.QueryParam("status")

	param := service.GetReportsParam{
		UserID:     userID,
		Role:       role,
		DistrictID: uint(districtID),
		Status:     status,
		Page:       page,
		Limit:      limit,
	}

	result, err := rc.svc.GetAllReports(param)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get reports", echo.Map{
		"meta": echo.Map{
			"page":       result.Page,
			"limit":      result.Limit,
			"total_data": result.TotalData,
			"total_page": result.TotalPage,
		},
		"data": result.Data,
	})
}

func (rc *ReportController) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "report id is not valid")
	}

	userID := helper.GetUserID(c)
	role := helper.GetUserRole(c)

	result, err := rc.svc.GetReportByID(uint(reportID), userID, role)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get report", result)
}

func (rc *ReportController) Update(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "report id is not valid")
	}

	userID := helper.GetUserID(c)

	var input service.UpdateReportInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	result, err := rc.svc.UpdateReport(uint(reportID), userID, input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success update report", result)
}

func (rc *ReportController) Assign(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "report id is not valid")
	}

	adminUserID := helper.GetUserID(c)
	role := helper.GetUserRole(c)

	var input service.AssignReportInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	err = rc.svc.AssignReport(uint(reportID), adminUserID, role, input)
	if err != nil {
		if err.Error() == "staff_id required" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "Success assign officer to this report", nil)
}

func (rc *ReportController) Start(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "report id is not valid")
	}

	officerUserID := helper.GetUserID(c)
	role := helper.GetUserRole(c)

	var input service.UpdateStatusReportInput
	_ = c.Bind(&input)

	err = rc.svc.StartReport(uint(reportID), officerUserID, role, input.Notes)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "Report handled started. Status: in progress", nil)
}

func (rc *ReportController) Resolve(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "report id is not valid")
	}

	officerUserID := helper.GetUserID(c)
	role := helper.GetUserRole(c)

	form, err := c.MultipartForm()
	if err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	notes := c.FormValue("notes")
	files := form.File["images"]

	err = rc.svc.ResolveReport(uint(reportID), officerUserID, role, notes, files)
	if err != nil {
		if err.Error() == "upload minimal 1 image" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "Report mark as resolve", nil)
}

func (rc *ReportController) Reject(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "report id is not valid")
	}

	adminUserID := helper.GetUserID(c)
	role := helper.GetUserRole(c)

	var input service.UpdateStatusReportInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	err = rc.svc.RejectReport(uint(reportID), adminUserID, role, input)
	if err != nil {
		if err.Error() == "notes required!" {
			return helper.BadRequest(c, err.Error())
		}
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success reject report", nil)
}

func (rc *ReportController) Verify(c echo.Context) error {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "report id is not valid")
	}

	userID := helper.GetUserID(c)
	role := helper.GetUserRole(c)

	var input service.UpdateStatusReportInput
	_ = c.Bind(&input)

	err = rc.svc.VerifyReport(uint(id64), userID, role, input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success verified report", nil)
}
