package exerciselog

import (
	"log"
	"net/http"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/exerciselog/dto"
	"github.com/lionel/my-fitness-aibo/backend/internal/api/exerciselog/dto/mapper"
	internalexercise "github.com/lionel/my-fitness-aibo/backend/internal/exercise"
	internalexerciselog "github.com/lionel/my-fitness-aibo/backend/internal/exerciselog"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/messages"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/request"
)

type Handler interface {
	ListExerciseLogs() http.HandlerFunc
	UpsertExerciseLog() http.HandlerFunc
}

type handler struct {
	logger            *log.Logger
	exerciseService   internalexercise.Service
	exerciseLogService internalexerciselog.Service
}

func NewHandler(
	logger *log.Logger,
	exerciseService internalexercise.Service,
	exerciseLogService internalexerciselog.Service,
) Handler {
	return &handler{
		logger:             logger,
		exerciseService:    exerciseService,
		exerciseLogService: exerciseLogService,
	}
}

func (h *handler) ListExerciseLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		exerciseID, ok := request.PathUUID(w, r, "exerciseID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid exerciseID")
			return
		}

		if _, err := h.exerciseService.GetExercise(ctx, userID, exerciseID); err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ListExerciseLogs", err, render.MapStoreError)
			return
		}

		logs, err := h.exerciseLogService.ListExerciseLogs(ctx, userID, exerciseID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ListExerciseLogs", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.ExerciseLogsFromDomain(logs))
	}
}

func (h *handler) UpsertExerciseLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		exerciseID, ok := request.PathUUID(w, r, "exerciseID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid exerciseID")
			return
		}

		cycleID, ok := request.PathUUID(w, r, "cycleID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid cycleID")
			return
		}

		var req dto.UpsertLogRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}

		logEntry, err := h.exerciseLogService.UpsertExerciseLog(ctx, userID, mapper.UpsertLogRequestToDomain(exerciseID, cycleID, req))
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "UpsertExerciseLog", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.ExerciseLogFromDomain(logEntry))
	}
}
