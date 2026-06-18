package cycle

import (
	"log"
	"net/http"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/cycle/dto"
	"github.com/lionel/my-fitness-aibo/backend/internal/api/cycle/dto/mapper"
	internalcycle "github.com/lionel/my-fitness-aibo/backend/internal/cycle"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/messages"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/request"
)

type Handler interface {
	ListCycles() http.HandlerFunc
	CreateCycle() http.HandlerFunc
	UpdateCycle() http.HandlerFunc
	DeleteCycle() http.HandlerFunc
}

type handler struct {
	logger       *log.Logger
	cycleService internalcycle.Service
}

func NewHandler(logger *log.Logger, cycleService internalcycle.Service) Handler {
	return &handler{logger: logger, cycleService: cycleService}
}

func (h *handler) ListCycles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		planID, ok := request.PathUUID(w, r, "planID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid planID")
			return
		}

		cycles, err := h.cycleService.ListCycles(ctx, userID, planID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ListCycles", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.CyclesFromDomain(cycles))
	}
}

func (h *handler) CreateCycle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		planID, ok := request.PathUUID(w, r, "planID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid planID")
			return
		}

		var req dto.CreateCycleRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}

		cycle, err := h.cycleService.CreateCycle(ctx, userID, planID, mapper.CreateCycleRequestToDomain(req))
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "CreateCycle", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusCreated, dto.CycleFromDomain(cycle))
	}
}

func (h *handler) UpdateCycle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		cycleID, ok := request.PathUUID(w, r, "cycleID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid cycleID")
			return
		}

		var req dto.UpdateCycleRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}

		cycle, err := h.cycleService.UpdateCycle(ctx, userID, cycleID, mapper.UpdateCycleRequestToDomain(req))
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "UpdateCycle", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.CycleFromDomain(cycle))
	}
}

func (h *handler) DeleteCycle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		cycleID, ok := request.PathUUID(w, r, "cycleID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid cycleID")
			return
		}

		if err := h.cycleService.DeleteCycle(ctx, userID, cycleID); err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "DeleteCycle", err, render.MapStoreError)
			return
		}

		render.NoContent(w)
	}
}
