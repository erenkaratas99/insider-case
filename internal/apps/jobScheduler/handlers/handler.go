package handlers

import (
	"github.com/labstack/echo/v4"
	"insider/internal/apps/jobScheduler/jobs"
)

type Handler struct {
	e      *echo.Echo
	prefix string
	mj     *jobs.MessengerJob
}

func NewHandler(e *echo.Echo, prefix string, mj *jobs.MessengerJob) *Handler {
	return &Handler{
		e:      e,
		prefix: prefix,
		mj:     mj,
	}
}
