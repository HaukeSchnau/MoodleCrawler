package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Course struct {
	Id           string   `json:"id"`
	NumSections  int      `json:"numsections"`
	SectionList  []string `json:"sectionlist"`
	EditMode     bool     `json:"editmode"`
	Highlighted  string   `json:"highlighted"`
	MaxSections  string   `json:"maxsections"`
	Baseurl      string   `json:"baseurl"`
	StateKey     string   `json:"statekey"`
	MaxBytes     int64    `json:"maxbytes"`
	MaxBytesText string   `json:"maxbytestext"`
}

type Section struct {
	Id               string   `json:"id"`
	Section          int      `json:"section"`
	Number           int      `json:"number"`
	Title            string   `json:"title"`
	HasSummary       bool     `json:"hassummary"`
	RawTitle         string   `json:"rawtitle"`
	CmList           []string `json:"cmlist"`
	Visible          bool     `json:"visible"`
	SectionUrl       string   `json:"sectionurl"`
	Current          bool     `json:"current"`
	IndexCollapsed   bool     `json:"indexcollapsed"`
	ContentCollapsed bool     `json:"contentcollapsed"`
	HasRestrictions  bool     `json:"hasrestrictions"`
	BulkEditable     bool     `json:"bulkeditable"`
}

type Cm struct {
	Id                string      `json:"id"`
	Anchor            string      `json:"anchor"`
	Name              string      `json:"name"`
	Visible           bool        `json:"visible"`
	Stealth           bool        `json:"stealth"`
	SectionId         string      `json:"sectionid"`
	SectionNumber     int         `json:"sectionnumber"`
	UserVisible       bool        `json:"uservisible"`
	HasCmRestrictions bool        `json:"hascmrestrictions"`
	ModName           string      `json:"modname"`
	Indent            int         `json:"indent"`
	GroupMode         json.Number `json:"groupmode"`
	Module            string      `json:"module"`
	Plugin            string      `json:"plugin"`
	AccessVisible     bool        `json:"accessvisible"`
	Url               string      `json:"url"`
	IsTrackedUser     bool        `json:"istrackeduser"`
	AllowStealth      bool        `json:"allowstealth"`
}

type CourseState struct {
	Course  Course    `json:"course"`
	Section []Section `json:"section"`
	Cm      []Cm      `json:"cm"`
}

type GetCoursePayloadArgs struct {
	CourseId int `json:"courseid"`
}

func getCourseResources(client *http.Client, sessionKey string, courseId int) CourseState {
	args := GetCoursePayloadArgs{
		CourseId: courseId,
	}

	data := makeServiceRequest[string](client, args, sessionKey, "core_courseformat_get_state")

	var courseState CourseState
	err := json.Unmarshal([]byte(data), &courseState)

	if err != nil {
		log.Fatalf("Error while parsing course format %s", err)
	}

	return courseState
}
