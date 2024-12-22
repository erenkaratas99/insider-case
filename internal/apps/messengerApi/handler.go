package messengerApi

import (
	"github.com/labstack/echo/v4"
	"insider/configs/appConfigs"
	"insider/internal/apps/messengerApi/entities"
	"insider/internal/clients"
	"insider/pkg"
	"net/http"
	"strconv"
)

type Handler struct {
	e                  *echo.Echo
	messengerService   *MessengerService
	validator          entities.Validator
	jobSchedulerClient *clients.JobSchedulerClient
}

func NewHandler(e *echo.Echo, messengerService *MessengerService, cfg *appConfigs.Configurations) *Handler {
	return &Handler{
		e:                  e,
		messengerService:   messengerService,
		validator:          entities.NewValidator(),
		jobSchedulerClient: clients.NewJobSchedulerClient(cfg),
	}
}

func (h *Handler) RegisterRoutes(prefix string) {
	g := h.e.Group(prefix)
	g.GET("/messenger-job-toggle", h.ToggleMessengerJob)
	g.GET("/", h.GetAll)
	g.GET("/get-two", h.GetTwo)
	g.PUT("/commit/:messageId", h.Commit)
}

func (h *Handler) ToggleMessengerJob(c echo.Context) error {
	command := c.QueryParam("command")

	if command != "start" && command != "stop" {
		return pkg.NewErrorResponse(http.StatusBadRequest, "Invalid command. Use 'start' or 'stop'.", 0).JSON(c)
	}

	res, err := h.jobSchedulerClient.ToggleJob("messenger", command)
	if err != nil {
		return pkg.NewErrorResponse(http.StatusInternalServerError, err.Error(), 0).JSON(c)
	}

	if res.StatusCode() != http.StatusOK {
		return pkg.NewErrorResponse(res.StatusCode(), "there's an error occurred!", 0).JSON(c)
	}

	return pkg.NewSuccessResponse().JSON(c)
}

func (h *Handler) GetAll(c echo.Context) error {

	limit, offset := c.QueryParam("limit"), c.QueryParam("offset")

	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return pkg.NewErrorResponse(http.StatusBadRequest, "Limit must be an integer!", 0).JSON(c)
	}

	o, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return pkg.NewErrorResponse(http.StatusBadRequest, "Offset must be an integer!", 0).JSON(c)
	}

	req := &entities.GetAllRequest{
		Limit:  l,
		Offset: o,
	}

	h.validator.Validate(req)

	messages, err := h.messengerService.GetAll(req)
	if err != nil {
		return pkg.NewErrorResponse(http.StatusInternalServerError, "Couldn't fetch data!", 0).JSON(c)
	}

	return messages.JSON(c)
}

func (h *Handler) GetTwo(c echo.Context) error {
	offset := c.QueryParam("offset")

	o, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return pkg.NewErrorResponse(http.StatusBadRequest, "Offset must be an integer!", 0).JSON(c)
	}

	if o < -1 {
		o = 0
	}

	messages, err := h.messengerService.GetTwo(o)
	if err != nil {
		return pkg.NewErrorResponse(http.StatusInternalServerError, "Couldn't fetch data!", 0).JSON(c)
	}
	return messages.JSON(c)
}

func (h *Handler) Commit(c echo.Context) error {
	messageId := c.Param("messageId")
	if messageId == "" {
		return pkg.NewErrorResponse(http.StatusBadRequest, "Invalid messageId.", 0).JSON(c)
	}

	if err := h.messengerService.CommitMessage(messageId); err != nil {
		return pkg.NewErrorResponse(http.StatusInternalServerError,
			"Couldn't commit the message!", 0).JSON(c)
	}

	return pkg.NewSuccessResponse().JSON(c)
}
