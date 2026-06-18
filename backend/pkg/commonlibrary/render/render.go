package render

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
	commonMappers "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/mappers"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/messages"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

func Json(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set(headerContentType, contentTypeJSON)

	body, err := json.Marshal(payload)
	if err != nil {
		statusCode = http.StatusInternalServerError
		body = []byte(fmt.Sprintf(`{"error":"%s"}`, err))
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}

func Error(w http.ResponseWriter, status int, msg string) {
	Json(w, status, commonMappers.ToSimpleErrorResponse(msg))
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func ErrorCausedByTimeoutOrClientCancellation(w http.ResponseWriter, r *http.Request, logger *log.Logger, err error) bool {
	switch {
	case errors.Is(err, context.Canceled):
		logger.Printf("client canceled request: path=%s", r.URL.Path)
		return true
	case errors.Is(err, context.DeadlineExceeded):
		Json(w, http.StatusGatewayTimeout, commonMappers.ToSimpleErrorResponse("request timed out"))
		return true
	default:
		return false
	}
}

func MapStoreError(err error) (int, string) {
	switch {
	case errors.Is(err, commonErrors.ErrNotFound):
		return http.StatusNotFound, messages.NotFoundMsg
	case errors.Is(err, commonErrors.ErrConflict):
		return http.StatusConflict, messages.ConflictMsg
	default:
		return http.StatusInternalServerError, messages.InternalServerErrorMsg
	}
}

func HandleServiceErrorResponse(
	logger *log.Logger,
	w http.ResponseWriter,
	r *http.Request,
	handlerName string,
	err error,
	errorsToStatusCodeAndMessageMapper func(err error) (int, string),
) {
	if ErrorCausedByTimeoutOrClientCancellation(w, r, logger, err) {
		return
	}

	statusCode, errMsg := errorsToStatusCodeAndMessageMapper(err)

	switch {
	case statusCode == http.StatusInternalServerError:
		logger.Printf("%s failure: %v", handlerName, err)
	case statusCode >= 400:
		logger.Printf("%s failure: %v", handlerName, err)
	}

	Json(w, statusCode, commonMappers.ToSimpleErrorResponse(errMsg))
}
