package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TheMaru/training-organiser/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary Create new curriculum plan
// @Description creates a new curriculum plan, it's the overarching structure for the curriculum
// @Tags curriculum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body CreateCurriculumPlanRequest true "Curriculum Plan Request json"
// @Success 201 {object} CurriculumPlanResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /curriculum/plans [post]
func (cfg *ApiConfig) HandleCreateCurriculumPlan(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := CreateCurriculumPlanRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	dbPlan, err := cfg.DB.CreateCurriculumPlan(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create Plan", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, PlanResponse{
		ID:        dbPlan.ID,
		CreatedAt: dbPlan.CreatedAt,
		Name:      dbPlan.Name,
	})
}

// @Summary Get curriculum plans
// @Description Returns all curriculum plans
// @Tags curriculum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} PlanResponse
// @Failure 500 {object} ErrorResponse
// @Router /curriculum/plans [get]
func (cfg *ApiConfig) HandleListCurriculumPlans(w http.ResponseWriter, r *http.Request) {
	dbPlans, err := cfg.DB.ListCurriculumPlans(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not load plans", err)
		return
	}

	plans := make([]PlanResponse, 0)
	for _, plan := range dbPlans {
		plans = append(plans, PlanResponse{
			ID:        plan.ID,
			CreatedAt: plan.CreatedAt,
			Name:      plan.Name,
		})
	}

	respondWithJSON(w, http.StatusOK, plans)
}

// @Summary Add a slot to a plan
// @Description Plan is the overaching structure and can have multiple slots for topics, this adds a topic to a plan
// @Tags curriculum
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body AddSlotToPlanRequest true "struct AddSlotToPlanRequest"
// @Success 201 {object} AddSlotToPlanResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /curriculum/plans/{id}/slots [post]
func (cfg *ApiConfig) HandleAddSlotToPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid plan ID", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := AddSlotToPlanRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if params.DurationUnit != "weeks" && params.DurationUnit != "months" {
		respondWithError(w, http.StatusBadRequest, "duration_unit must be 'weeks' or 'months'", fmt.Errorf("Wrong duration_unit in request: %v", params.DurationUnit))
		return
	}

	if params.DurationValue < 1 {
		respondWithError(w, http.StatusBadRequest, "duration_value must be positive", fmt.Errorf("Wrong duration_value in request: %v", params.DurationValue))
		return
	}

	slot, err := cfg.DB.AddSlotToPlan(r.Context(), database.AddSlotToPlanParams{
		PlanID:        planID,
		TopicID:       params.TopicID,
		SequenceOrder: params.SequenceOrder,
		DurationValue: params.DurationValue,
		DurationUnit:  params.DurationUnit,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not add slot", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, AddSlotToPlanResponse{
		PlanID:        slot.PlanID,
		TopicID:       slot.TopicID,
		SequenceOrder: slot.SequenceOrder,
		DurationValue: slot.DurationValue,
		DurationUnit:  slot.DurationUnit,
	})
}
