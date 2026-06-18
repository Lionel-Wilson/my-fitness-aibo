package exercise

import (
	"log"
	"net/http"
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/exercise/dto"
	"github.com/lionel/my-fitness-aibo/backend/internal/api/exercise/dto/mapper"
	internalexercise "github.com/lionel/my-fitness-aibo/backend/internal/exercise"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/messages"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/request"
)

type Handler interface {
	ListExercises() http.HandlerFunc
	CreateExercise() http.HandlerFunc
	GetExercise() http.HandlerFunc
	UpdateExercise() http.HandlerFunc
	DeleteExercise() http.HandlerFunc
}

type handler struct {
	logger          *log.Logger
	exerciseService internalexercise.Service
}

func NewHandler(logger *log.Logger, exerciseService internalexercise.Service) Handler {
	return &handler{logger: logger, exerciseService: exerciseService}
}

func (h *handler) ListExercises() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		workoutID, ok := request.PathUUID(w, r, "workoutID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid workoutID")
			return
		}

		exercises, err := h.exerciseService.ListExercises(ctx, userID, workoutID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ListExercises", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.ExercisesFromDomain(exercises))
	}
}

func (h *handler) CreateExercise() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		workoutID, ok := request.PathUUID(w, r, "workoutID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid workoutID")
			return
		}

		var req dto.ExerciseRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			render.Error(w, http.StatusBadRequest, "name is required")
			return
		}

		exercise, err := h.exerciseService.CreateExercise(ctx, userID, workoutID, mapper.ExerciseRequestToDomain(req))
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "CreateExercise", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusCreated, dto.ExerciseFromDomain(exercise))
	}
}

func (h *handler) GetExercise() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		exerciseID, ok := request.PathUUID(w, r, "exerciseID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid exerciseID")
			return
		}

		exercise, err := h.exerciseService.GetExercise(ctx, userID, exerciseID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "GetExercise", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.ExerciseFromDomain(exercise))
	}
}

func (h *handler) UpdateExercise() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		exerciseID, ok := request.PathUUID(w, r, "exerciseID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid exerciseID")
			return
		}

		var req dto.ExerciseRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			render.Error(w, http.StatusBadRequest, "name is required")
			return
		}

		exercise, err := h.exerciseService.UpdateExercise(ctx, userID, exerciseID, mapper.ExerciseRequestToDomain(req))
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "UpdateExercise", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.ExerciseFromDomain(exercise))
	}
}

func (h *handler) DeleteExercise() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		exerciseID, ok := request.PathUUID(w, r, "exerciseID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid exerciseID")
			return
		}

		if err := h.exerciseService.DeleteExercise(ctx, userID, exerciseID); err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "DeleteExercise", err, render.MapStoreError)
			return
		}

		render.NoContent(w)
	}
}
