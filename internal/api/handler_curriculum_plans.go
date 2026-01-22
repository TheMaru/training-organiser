package api

import (
	"encoding/json"
	"net/http"
)

// @Summary Create new curriculum plan
// @Description creates a new curriculum plan, it's the overarching structure for the curriculum
// @Tags curriculum
// @Accept json
// @Produce json
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
