package main

import (
	"net/http"
)

type GetCoursesPayloadArgs struct {
	Classification   string `json:"classification"`
	CustomFieldName  string `json:"customfieldname"`
	CustomFieldValue string `json:"customfieldvalue"`
	Limit            int    `json:"limit"`
	Offset           int    `json:"offset"`
	Sort             string `json:"sort"`
}

type FullCourse struct {
	Id                       int    `json:"id"`
	FullName                 string `json:"fullname"`
	Shortname                string `json:"shortname"`
	IdNumber                 string `json:"idnumber"`
	Summary                  string `json:"summary"`
	SummaryFormat            int    `json:"summaryformat"`
	StartDate                int    `json:"startdate"`
	EndDate                  int    `json:"enddate"`
	Visible                  bool   `json:"visible"`
	ShowActivityDates        bool   `json:"showactivitydates"`
	ShowCompletionConditions bool   `json:"showcompletionconditions"`
	PdfExportFont            string `json:"pdfexportfont"`
	FullNameDisplay          string `json:"fullnamedisplay"`
	ViewUrl                  string `json:"viewurl"`
	CourseImage              string `json:"courseimage"`
	Progress                 int    `json:"progress"`
	HasProgress              bool   `json:"hasprogress"`
	IsFavourite              bool   `json:"isfavourite"`
	Hidden                   bool   `json:"hidden"`
	ShowShortName            bool   `json:"showshortname"`
	CourseCategory           string `json:"coursecategory"`
}

type GetCoursesResponse struct {
	Courses    []FullCourse `json:"courses"`
	NextOffset int          `json:"nextoffset"`
}

func getCourses(client *http.Client, sessionKey string) GetCoursesResponse {
	args := GetCoursesPayloadArgs{
		Classification:   "all",
		CustomFieldName:  "",
		CustomFieldValue: "",
		Limit:            24,
		Offset:           0,
		Sort:             "fullname",
	}

	return makeServiceRequest[GetCoursesResponse](client, args, sessionKey, "core_course_get_enrolled_courses_by_timeline_classification")
}
