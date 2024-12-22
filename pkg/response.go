package pkg

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type BaseResponse struct {
	HttpStatus   int         `json:"status"`
	InternalCode int         `json:"internalCode"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(data ...interface{}) *BaseResponse {
	var value interface{}
	if len(data) > 0 {
		value = data[0]
	} else {
		value = "succeeded"
	}
	return &BaseResponse{
		HttpStatus:   http.StatusOK,
		InternalCode: 0,
		Data:         value,
	}
}

// NewErrorResponse is a helper to create an error types with code and message.
func NewErrorResponse(httpStatus int, message string, internalCode int) *BaseResponse {
	return &BaseResponse{
		HttpStatus:   httpStatus,
		InternalCode: internalCode,
		Message:      message,
	}
}

func (br *BaseResponse) JSON(c echo.Context) error {
	return c.JSON(br.HttpStatus, br)
}
