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

// @Summary Toggle Messenger Job
// @Description Starts or stops the Messenger Job based on the provided command by sending http requests to the job scheduler.
// @Tags MessengerAPI
// @Accept  json
// @Produce  json
// @Param command query string true "Command to toggle job. Use 'start' or 'stop'"
// @Success 200 {object} pkg.BaseResponse "Job toggled successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid command parameter"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /messenger-job-toggle [get]
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

// @Summary Get All Messages
// @Description Retrieves all messenger messages with pagination.
// @Tags MessengerAPI
// @Accept  json
// @Produce  json
// @Param limit query integer true "Number of items to retrieve" minimum(1)
// @Param offset query integer true "Number of items to skip" minimum(0)
// @Success 200 {array} pkg.BaseResponse "List of messages under field of 'data'"
// @Failure 400 {object} pkg.BaseResponse "Invalid query parameters"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router / [get]
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

// @Summary Get Two Messages for Messenger Job
// @Description Retrieves two messenger messages based on the provided offset, tailored for the Messenger Job.
// @Tags MessengerAPI
// @Accept  json
// @Produce  json
// @Param offset query integer true "Offset for retrieving messages" minimum(-1)
// @Success 200 {array} pkg.BaseResponse "List of two messages, projection of 'to' and 'content'"
// @Failure 400 {object} pkg.BaseResponse "Invalid query parameter"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /get-two [get]
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

// @Summary Commit a Message
// @Description Commits a messenger message based on the provided message ID, if success, changes the message status to sent.
// @Tags MessengerAPI
// @Accept  json
// @Produce  json
// @Param messageId path string true "ID of the message to commit"
// @Success 200 {object} pkg.BaseResponse "Message committed successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid message ID"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /commit/{messageId} [put]
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
