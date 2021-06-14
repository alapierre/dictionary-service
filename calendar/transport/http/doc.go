package http

import (
	"dictionaries-service/calendar"
)

// Package http define endpoints for calendar and calendar type load, save, update and delete

//swagger:parameters saveCalendar updateCalendar
//goland:noinspection GoUnusedType
type saveDtoWrapper struct {

	// Calendar type id
	// in:path
	Type string `json:"type"`

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// Calendar body
	// in:body
	Body calendar.SaveDto
}

//swagger:parameters saveCalendarType updateCalendarType
//goland:noinspection GoUnusedType
type createCalendarTypeWrapper struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// Calendar type body
	// in:body
	Body calendar.DictionaryCalendarType
}

//swagger:parameters loadCalendarTypes
//goland:noinspection GoUnusedType
type calendarTypeRequest struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`
}

//swagger:response calendarResponse
//goland:noinspection GoUnusedType
type calendarResponse struct {
	Day    string            `json:"day"`
	Tenant string            `json:"tenant,omitempty"`
	Name   *string           `json:"name"`
	Kind   *string           `json:"kind,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

//swagger:response calendarResponse
//goland:noinspection GoUnusedType
type calendarResponseWrapper struct {
	// in:body
	Body calendarResponse
}

//swagger:response calendarTypeResponse
//goland:noinspection GoUnusedType
type calendarTypeResponseWrapper struct {

	// in:body
	Body []calendar.DictionaryCalendarType
}

// swagger:parameters deleteCalendarType
//goland:noinspection GoUnusedType
type calendarTypeDeleteWrapper struct {

	// in:path
	Type string `json:"type"`

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`
}
