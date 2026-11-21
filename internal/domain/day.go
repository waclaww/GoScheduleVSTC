package domain

import (
	"errors"
	"time"
)

var (
	ErrDayNotFound = errors.New("day not found")
	ErrInvalidDay  = errors.New("invalid day")
)

type Day struct {
	Id         int64
	Date       time.Time
	Title      string
	SourceHash string
	FetchedAt  time.Time
	Groups     []string
	Slots      []BellSlot
	Lessons    []Lesson
}

type BellSlot struct {
	Slot     int
	StartsAt string
	EndsAt   string
}

type Lesson struct {
	Group   string
	Slot    int
	Subject string
	Room    string
}

type GroupDay struct {
	Date    time.Time
	Title   string
	Group   string
	Lessons []GroupLesson
}

type GroupLesson struct {
	Slot     int
	StartsAt string
	EndsAt   string
	Subject  string
	Room     string
}

type UpsertResult struct {
	Created bool
	Changed bool
}
