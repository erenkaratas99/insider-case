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

// @Summary Start Messenger Job
// @Description Initiates the Messenger Job and returns the operation status.
// @Tags MessengerJob
// @Accept json
// @Produce json
// @Success 200 {object} pkg.BaseResponse "success"
// @Failure 400 {object} pkg.BaseResponse "bad request"
// @Router /messenger/start [get]
func (h *Handler) Start(c echo.Context) error {
	if jobDto := h.mj.Start(); jobDto.Status != "success" {
		return pkg.NewErrorResponse(http.StatusBadRequest, jobDto.Message, jobDto.Code).JSON(c)
	}

	return pkg.NewSuccessResponse().JSON(c)
}

// @Summary Stop Messenger Job
// @Description Terminates the Messenger Job and returns the operation status.
// @Tags MessengerJob
// @Accept json
// @Produce json
// @Success 200 {object} pkg.BaseResponse "success"
// @Failure 400 {object} pkg.BaseResponse "Invalid request"
// @Router /messenger/stop [get]
func (h *Handler) Stop(c echo.Context) error {
	if jobDto := h.mj.Stop(); jobDto.Status != "success" {
		return pkg.NewErrorResponse(http.StatusBadRequest, jobDto.Message, jobDto.Code).JSON(c)
	}

	return pkg.NewSuccessResponse().JSON(c)
}

// @Summary Check Messenger Job Status
// @Description Checks whether the Messenger Job is currently running.
// @Tags MessengerJob
// @Accept json
// @Produce json
// @Success 200 {object} pkg.BaseResponse "Current job status"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /messenger/is-working [get]
func (h *Handler) IsWorking(c echo.Context) error {
	if h.mj.IsRunning() {
		return pkg.NewSuccessResponse("Working").JSON(c)
	}

	return pkg.NewSuccessResponse("Not Working").JSON(c)
}
