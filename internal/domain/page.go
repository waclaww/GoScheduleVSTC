package domain

import (
	"errors"
	"time"
)

var (
	ErrPageNotFound = errors.New("page not found")
	ErrInvalidPage  = errors.New("invalid page")
)

type Page struct {
	Name      string
	Date      time.Time
	Title     string
	UpdatedAt time.Time
}

type GroupPages struct {
	Group    string
	Today    *GroupDay
	Tomorrow *GroupDay
}

type PageResult struct {
	Name         string
	Date         time.Time
	PreviousDate *time.Time
	Changed      bool
}
