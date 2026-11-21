package http

import "gateway/internal/adapter/http/handlers"

type Handlers struct {
	General  *handlers.General
	Schedule *handlers.Schedule
	Archive  *handlers.Archive
	Day      *handlers.Day
	Page     *handlers.Page
}
