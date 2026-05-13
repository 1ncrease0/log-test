package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"log-parser/internal/application"
)

func TestMapErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
		wantLogErr bool
	}{
		{
			name:       "archive_not_found",
			err:        application.ErrArchiveNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    application.ErrArchiveNotFound.Error(),
			wantLogErr: false,
		},
		{
			name:       "invalid_path",
			err:        application.ErrInvalidPath,
			wantStatus: http.StatusBadRequest,
			wantMsg:    application.ErrInvalidPath.Error(),
			wantLogErr: false,
		},
		{
			name:       "duplicate_path",
			err:        application.ErrDuplicateLogPath,
			wantStatus: http.StatusConflict,
			wantMsg:    application.ErrDuplicateLogPath.Error(),
			wantLogErr: false,
		},
		{
			name:       "parse_failed",
			err:        application.ErrParseFailed,
			wantStatus: http.StatusUnprocessableEntity,
			wantMsg:    "invalid or incomplete log archive",
			wantLogErr: false,
		},
		{
			name:       "persist_failed",
			err:        application.ErrPersistFailed,
			wantStatus: http.StatusInternalServerError,
			wantMsg:    http.StatusText(http.StatusInternalServerError),
			wantLogErr: true,
		},
		{
			name:       "parse_failed_wrapped",
			err:        errors.Join(application.ErrParseFailed, errors.New("inner")),
			wantStatus: http.StatusUnprocessableEntity,
			wantMsg:    "invalid or incomplete log archive",
			wantLogErr: false,
		},
		{
			name:       "unknown",
			err:        errors.New("oops"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    http.StatusText(http.StatusInternalServerError),
			wantLogErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			st, msg, logErr := mapErrors(tc.err)
			require.Equal(t, tc.wantStatus, st)
			require.Equal(t, tc.wantMsg, msg)
			require.Equal(t, tc.wantLogErr, logErr)
		})
	}
}
