package progress

import (
	"log"
	"net/http"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/progress/dto"
	internalexercise "github.com/lionel/my-fitness-aibo/backend/internal/exercise"
	internalprogress "github.com/lionel/my-fitness-aibo/backend/internal/progress"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/request"
)

type Handler interface {
	ExerciseProgress() http.HandlerFunc
}

type handler struct {
	logger          *log.Logger
	exerciseService internalexercise.Service
	progressService internalprogress.Service
}

func NewHandler(
	logger *log.Logger,
	exerciseService internalexercise.Service,
	progressService internalprogress.Service,
) Handler {
	return &handler{
		logger:          logger,
		exerciseService: exerciseService,
		progressService: progressService,
	}
}

func (h *handler) ExerciseProgress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		exerciseID, ok := request.PathUUID(w, r, "exerciseID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid exerciseID")
			return
		}

		if _, err := h.exerciseService.GetExercise(ctx, userID, exerciseID); err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ExerciseProgress", err, render.MapStoreError)
			return
		}

		points, err := h.progressService.ExerciseProgress(ctx, userID, exerciseID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ExerciseProgress", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.ProgressPointsFromDomain(points))
	}
}
