package domain

import "errors"

var (
	ErrGroupNotFound     = errors.New("group not found")
	ErrScheduleNotEnough = errors.New("not enough schedules")
	ErrScheduleMalformed = errors.New("schedule malformed")
)

type Schedule struct {
	Id       int64
	Date     string
	Schedule string
}

type Archive struct {
	Id       int64
	Date     string
	Schedule string
}

type GroupSchedule struct {
	Today    string
	Tomorrow string
}
