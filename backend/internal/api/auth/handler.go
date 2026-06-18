package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/auth/dto"
	"github.com/lionel/my-fitness-aibo/backend/internal/api/auth/dto/mapper"
	internalauth "github.com/lionel/my-fitness-aibo/backend/internal/auth"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/messages"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/request"
)

type Handler interface {
	Signup() http.HandlerFunc
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
	Me() http.HandlerFunc
}

type handler struct {
	logger      *log.Logger
	authService internalauth.Service
}

func NewHandler(logger *log.Logger, authService internalauth.Service) Handler {
	return &handler{logger: logger, authService: authService}
}

func (h *handler) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req dto.CredentialsRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}

		creds := mapper.CredentialsRequestToDomain(req)
		if !strings.Contains(creds.Email, "@") {
			render.Error(w, http.StatusBadRequest, "a valid email is required")
			return
		}
		if len(creds.Password) < 8 {
			render.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
			return
		}

		result, err := h.authService.Signup(ctx, creds)
		if err != nil {
			if errors.Is(err, commonErrors.ErrConflict) {
				render.Error(w, http.StatusConflict, messages.ConflictMsg)
				return
			}
			if msg, ok := mapSignupInternalError(err); ok {
				render.Error(w, http.StatusInternalServerError, msg)
				return
			}
			render.HandleServiceErrorResponse(h.logger, w, r, "Signup", err, mapErrorsToStatusCodeAndMessages)
			return
		}

		render.Json(w, http.StatusCreated, mapper.AuthResultToResponse(result))
	}
}

func (h *handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req dto.CredentialsRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}

		creds := mapper.CredentialsRequestToDomain(req)
		result, err := h.authService.Login(ctx, creds)
		if err != nil {
			if errors.Is(err, internalauth.ErrInvalidCredentials) {
				render.Error(w, http.StatusUnauthorized, "invalid email or password")
				return
			}
			if msg, ok := mapLoginInternalError(err); ok {
				render.Error(w, http.StatusInternalServerError, msg)
				return
			}
			render.HandleServiceErrorResponse(h.logger, w, r, "Login", err, mapErrorsToStatusCodeAndMessages)
			return
		}

		render.Json(w, http.StatusOK, mapper.AuthResultToResponse(result))
	}
}

func (h *handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		render.NoContent(w)
	}
}

func (h *handler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := commoncontext.UserID(ctx)
		if !ok {
			render.Error(w, http.StatusUnauthorized, "missing or malformed authorization header")
			return
		}

		result, err := h.authService.Me(ctx, userID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "Me", err, mapErrorsToStatusCodeAndMessages)
			return
		}

		render.Json(w, http.StatusOK, mapper.MeResultToResponse(result))
	}
}

func mapSignupInternalError(err error) (string, bool) {
	if strings.Contains(err.Error(), "hash password") {
		return "could not hash password", true
	}
	if strings.Contains(err.Error(), "issue token") {
		return "could not issue token", true
	}

	return "", false
}

func mapLoginInternalError(err error) (string, bool) {
	if strings.Contains(err.Error(), "issue token") {
		return "could not issue token", true
	}

	return "", false
}

func mapErrorsToStatusCodeAndMessages(err error) (int, string) {
	switch {
	case errors.Is(err, commonErrors.ErrNotFound):
		return http.StatusNotFound, messages.NotFoundMsg
	case errors.Is(err, commonErrors.ErrConflict):
		return http.StatusConflict, messages.ConflictMsg
	default:
		return http.StatusInternalServerError, messages.InternalServerErrorMsg
	}
}
