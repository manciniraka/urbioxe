package entity

import (
	"time"
)

type NewsCategory string

const (
	NewsCategoryAnnouncement 	NewsCategory = "announcement"
	NewsCategoryNews		 	NewsCategory = "news"
	NewsCategoryEvent		 	NewsCategory = "event"
	NewsCategoryEmergency		NewsCategory = "emergency"
)

type NewsScope string

const (
	NewsScopeGlobal		NewsScope = "global"
	NewsScopeDistrict	NewsScope = "district"
)

type RegionalNews struct {
	ID		  		int64 			`json:"id"`
	DepartmentID 	*int64 			`json:"department_id"`
	DistrictID 		*int64 			`json:"district_id"`
	Title	 		string 			`json:"title"`
	Content 		string 			`json:"content"`
	Category 		NewsCategory 	`json:"category"`
	BannerURL 		*string 			`json:"banner_url"`
	TargetScope 	NewsScope 		`json:"target_scope"`
	IsPinned 		bool 			`json:"is_pinned"`
	CreatedBy 		*int64 			`json:"created_by"`
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
}

type CreateNewsRequest struct {
	DepartmentID 	*int64 			`json:"department_id"`
	DistrictID 		*int64 			`json:"district_id"`
	Title	 		string 			`json:"title"`
	Content 		string 			`json:"content"`
	Category 		NewsCategory 	`json:"category"`
	BannerURL 		*string 			`json:"banner_url"`
	TargetScope 	NewsScope 		`json:"target_scope"`
	IsPinned 		bool 			`json:"is_pinned"`
}

type UpdateNewsRequest struct {
	DepartmentID 	*int64 			`json:"department_id"`
	DistrictID 		*int64 			`json:"district_id"`
	Title	 		string 			`json:"title"`
	Content 		string 			`json:"content"`
	Category 		NewsCategory 	`json:"category"`
	BannerURL 		*string 			`json:"banner_url"`
	TargetScope 	NewsScope 		`json:"target_scope"`
	IsPinned 		bool 			`json:"is_pinned"`
}
