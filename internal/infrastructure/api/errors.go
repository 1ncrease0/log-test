package api

import (
	"errors"
	"net/http"

	"log-parser/internal/application"
)

func mapErrors(err error) (status int, msg string, logErr bool) {
	switch {
	case errors.Is(err, application.ErrArchiveNotFound):
		return http.StatusNotFound, err.Error(), false
	case errors.Is(err, application.ErrInvalidPath):
		return http.StatusBadRequest, err.Error(), false
	case errors.Is(err, application.ErrDuplicateLogPath):
		return http.StatusConflict, err.Error(), false
	case errors.Is(err, application.ErrParseFailed):
		return http.StatusUnprocessableEntity, "invalid or incomplete log archive", false
	case errors.Is(err, application.ErrPersistFailed):
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), true
	default:
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), true
	}
}
