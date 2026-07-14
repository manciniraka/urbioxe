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

	userID, _ := c.Get("user_id").(int64)
	// test if not login
	if userID == 0 {
		userID = 6
	}

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

func (rc *ReportController) GetAll(c echo.Context) error {
	userID, _ := c.Get("user_id").(int64)
	role, _ := c.Get("role").(string)

	// test if not login
	if userID == 0 {
		userID = 3
	}
	if role == "" {
		role = "officer"
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	districtID, _ := strconv.ParseInt(c.QueryParam("district_id"), 10, 64)
	status := c.QueryParam("status")

	param := service.GetReportsParam{
		UserID:     userID,
		Role:       role,
		DistrictID: districtID,
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
	reportID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "Report ID not valid")
	}

	userID, _ := c.Get("user_id").(int64)
	role, _ := c.Get("role").(string)

	// test if not login
	if userID == 0 {
		userID = 3
	}
	if role == "" {
		role = "officer"
	}

	result, err := rc.svc.GetReportByID(reportID, userID, role)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success get report", result)
}

func (rc *ReportController) Update(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "Report ID not valid")
	}

	userID, _ := c.Get("user_id").(int64)
	// test not login
	if userID == 0 {
		userID = 1
	}

	var input service.UpdateReportInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	result, err := rc.svc.UpdateReport(reportID, userID, input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success update report", result)
}

func (rc *ReportController) Assign(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "Report ID not valid")
	}

	adminUserID, _ := c.Get("user_id").(int64)
	role, _ := c.Get("role").(string)
	// test login as admin
	if adminUserID == 0 {
		adminUserID = 4
	}
	if role == "" {
		role = "department_admin"
	}

	var input service.AssignReportInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	err = rc.svc.AssignReport(reportID, adminUserID, role, input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "Success assign officer to this report", nil)
}

func (rc *ReportController) Start(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "Report ID not valid")
	}

	officerUserID, _ := c.Get("user_id").(int64)
	role, _ := c.Get("role").(string)

	// test officer login
	if officerUserID == 0 {
		officerUserID = 3
	}
	if role == "" {
		role = "officer"
	}

	var input service.UpdateStatusReportInput
	_ = c.Bind(&input)

	err = rc.svc.StartReport(reportID, officerUserID, role, input.Notes)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "Report handled started. Status: in progress", nil)
}

func (rc *ReportController) Resolve(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "Report ID not valid")
	}

	officerUserID, _ := c.Get("user_id").(int64)
	role, _ := c.Get("role").(string)
	// test user
	if officerUserID == 0 {
		officerUserID = 3
	}
	if role == "" {
		role = "officer"
	}

	form, err := c.MultipartForm()
	if err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	notes := c.FormValue("notes")
	files := form.File["images"]

	err = rc.svc.ResolveReport(reportID, officerUserID, role, notes, files)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "Report mark as resolve", nil)
}

func (rc *ReportController) Reject(c echo.Context) error {
	idParam := c.Param("id")
	reportID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return helper.BadRequest(c, "Report ID not valid")
	}

	adminUserID, _ := c.Get("user_id").(int64)
	role, _ := c.Get("role").(string)

	// test without login
	if adminUserID == 0 {
		adminUserID = 4
	}
	if role == "" {
		role = "department_admin"
	}

	var input service.UpdateStatusReportInput
	if err := c.Bind(&input); err != nil {
		return helper.BadRequest(c, "format body not valid")
	}

	err = rc.svc.RejectReport(reportID, adminUserID, role, input)
	if err != nil {
		return helper.HandleError(c, err)
	}

	return helper.Success(c, "success reject report", nil)
}
