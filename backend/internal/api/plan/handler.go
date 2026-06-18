package plan

import (
	"log"
	"net/http"
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/plan/dto"
	"github.com/lionel/my-fitness-aibo/backend/internal/api/plan/dto/mapper"
	internalplan "github.com/lionel/my-fitness-aibo/backend/internal/plan"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/messages"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/request"
)

type Handler interface {
	ListPlans() http.HandlerFunc
	CreatePlan() http.HandlerFunc
	GetPlan() http.HandlerFunc
	UpdatePlan() http.HandlerFunc
	DeletePlan() http.HandlerFunc
}

type handler struct {
	logger      *log.Logger
	planService internalplan.Service
}

func NewHandler(logger *log.Logger, planService internalplan.Service) Handler {
	return &handler{logger: logger, planService: planService}
}

func (h *handler) ListPlans() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		plans, err := h.planService.ListPlans(ctx, userID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "ListPlans", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.PlansFromDomain(plans))
	}
}

func (h *handler) CreatePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		var req dto.PlanRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			render.Error(w, http.StatusBadRequest, "name is required")
			return
		}

		input, err := mapper.PlanRequestToDomain(req)
		if err != nil {
			render.Error(w, http.StatusBadRequest, "dates must be in YYYY-MM-DD format")
			return
		}

		plan, err := h.planService.CreatePlan(ctx, userID, input)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "CreatePlan", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusCreated, dto.PlanFromDomain(plan))
	}
}

func (h *handler) GetPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		planID, ok := request.PathUUID(w, r, "planID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid planID")
			return
		}

		plan, err := h.planService.GetPlan(ctx, userID, planID)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "GetPlan", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.PlanFromDomain(plan))
	}
}

func (h *handler) UpdatePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		planID, ok := request.PathUUID(w, r, "planID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid planID")
			return
		}

		var req dto.PlanRequest
		if err := request.DecodeAndValidate(r.Body, &req); err != nil {
			render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			render.Error(w, http.StatusBadRequest, "name is required")
			return
		}

		input, err := mapper.PlanRequestToDomain(req)
		if err != nil {
			render.Error(w, http.StatusBadRequest, "dates must be in YYYY-MM-DD format")
			return
		}

		plan, err := h.planService.UpdatePlan(ctx, userID, planID, input)
		if err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "UpdatePlan", err, render.MapStoreError)
			return
		}

		render.Json(w, http.StatusOK, dto.PlanFromDomain(plan))
	}
}

func (h *handler) DeletePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, _ := commoncontext.UserID(ctx)

		planID, ok := request.PathUUID(w, r, "planID")
		if !ok {
			render.Error(w, http.StatusBadRequest, "invalid planID")
			return
		}

		if err := h.planService.DeletePlan(ctx, userID, planID); err != nil {
			render.HandleServiceErrorResponse(h.logger, w, r, "DeletePlan", err, render.MapStoreError)
			return
		}

		render.NoContent(w)
	}
}
