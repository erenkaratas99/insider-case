package handlers

import (
	"github.com/labstack/echo/v4"
	"insider/pkg"
	"net/http"
)

func (h *Handler) RegisterRoutesForMessengerJob() {
	g := h.e.Group(h.prefix + "/messenger")
	g.GET("/start", h.Start)
	g.GET("/stop", h.Stop)
	g.GET("/is-working", h.IsWorking)
}

func (h *Handler) Start(c echo.Context) error {
	if jobDto := h.mj.Start(); jobDto.Status != "success" {
		return pkg.NewErrorResponse(http.StatusBadRequest, jobDto.Message, jobDto.Code).JSON(c)
	}

	return pkg.NewSuccessResponse().JSON(c)
}

func (h *Handler) Stop(c echo.Context) error {
	if jobDto := h.mj.Stop(); jobDto.Status != "success" {
		return pkg.NewErrorResponse(http.StatusBadRequest, jobDto.Message, jobDto.Code).JSON(c)
	}

	return pkg.NewSuccessResponse().JSON(c)
}

func (h *Handler) IsWorking(c echo.Context) error {
	if h.mj.IsRunning() {
		return pkg.NewSuccessResponse("Working").JSON(c)
	}

	return pkg.NewSuccessResponse("Not Working").JSON(c)
}
