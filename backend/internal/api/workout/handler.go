package workout

import (
	"log"
	"net/http"
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/workout/dto"
	"github.com/lionel/my-fitness-aibo/backend/internal/api/workout/dto/mapper"
	internalworkout "github.com/lionel/my-fitness-aibo/backend/internal/workout"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/messages"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/request"
)

type Handler interface {
	ListWorkouts() http.HandlerFunc
	CreateWorkout() http.HandlerFunc
	GetWorkout() http.HandlerFunc
	UpdateWorkout() http.HandlerFunc
	DeleteWorkout() http.HandlerFunc
}

type handler struct {
	logger         *log.Logger
	workoutService internalworkout.Service
}

func NewHandler(logger *log.Logger, workoutService internalworkout.Service) Handler {
	return &handler{logger: logger, workoutService: workoutService}
}

func (h *handler) ListWorkouts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		planID, ok := request.PathUUID(w, r, "planID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid planID")
			return
		}

		workouts, err := h.workoutService.ListWorkouts(ctx, userID, planID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ListWorkouts", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.WorkoutsFromDomain(workouts))
	}
}

func (h *handler) CreateWorkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		planID, ok := request.PathUUID(w, r, "planID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid planID")
			return
		}

		var req dto.WorkoutRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			render.Error(w, http.StatusBadRequest, "name is required")
			return
		}

		workout, err := h.workoutService.CreateWorkout(ctx, userID, planID, mapper.WorkoutRequestToDomain(req))
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "CreateWorkout", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusCreated, dto.WorkoutFromDomain(workout))
	}
}

func (h *handler) GetWorkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		workoutID, ok := request.PathUUID(w, r, "workoutID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid workoutID")
			return
		}

		workout, err := h.workoutService.GetWorkout(ctx, userID, workoutID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "GetWorkout", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.WorkoutFromDomain(workout))
	}
}

func (h *handler) UpdateWorkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		workoutID, ok := request.PathUUID(w, r, "workoutID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid workoutID")
			return
		}

		var req dto.WorkoutRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			render.Error(w, http.StatusBadRequest, "name is required")
			return
		}

		workout, err := h.workoutService.UpdateWorkout(ctx, userID, workoutID, mapper.WorkoutRequestToDomain(req))
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "UpdateWorkout", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.WorkoutFromDomain(workout))
	}
}

func (h *handler) DeleteWorkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		workoutID, ok := request.PathUUID(w, r, "workoutID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid workoutID")
			return
		}

		if err := h.workoutService.DeleteWorkout(ctx, userID, workoutID); err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "DeleteWorkout", err, render.MapStoreError)
			return
		}

		render.NoContent(w)
	}
}
