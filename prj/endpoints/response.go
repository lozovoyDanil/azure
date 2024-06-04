package endpoints

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	ErrorMessage string `json:"message"`
}

func newErrorResponse(code int, msg string) error {
	slog.Error(msg)
	return echo.NewHTTPError(code, errorResponse{ErrorMessage: msg})
}
