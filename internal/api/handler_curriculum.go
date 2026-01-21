package api

import (
	"encoding/json"
	"net/http"

	"github.com/TheMaru/training-organiser/internal/database"
)

// @Summary Create new curriculum topic
// @Description Creates a new (general) topic for the curriculum
// @Tags curriculum
// @Accept json
// @Produce json
// @Param request body CreateTopicRequest true "Topic"
// @Success 201 {object} TopicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /curriculum/topics [post]
func (cfg *ApiConfig) HandleCreateCurriculumTopic(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := CreateTopicRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	dbTopic, err := cfg.DB.CreateCurriculumTopic(r.Context(), database.CreateCurriculumTopicParams{
		Name:        params.Name,
		Description: stringToPointer(params.Description),
		ColorCode:   stringToPointer(params.ColorCode),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create topic", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTopicToResponse(dbTopic))
}

func databaseTopicToResponse(dbTopic database.CurriculumTopic) TopicResponse {
	var desc *string
	if dbTopic.Description.Valid {
		desc = &dbTopic.Description.String
	}

	var color *string
	if dbTopic.ColorCode.Valid {
		color = &dbTopic.ColorCode.String
	}

	return TopicResponse{
		ID:          dbTopic.ID,
		Name:        dbTopic.Name,
		Description: desc,
		ColorCode:   color,
		CreatedAt:   dbTopic.CreatedAt,
	}
}

// @Summary Get curriculum topics
// @Description Returns all curriculum topics
// @Tags curriculum
// @Accept json
// @Produce json
// @Success 200 {array} TopicResponse
// @Failure 500 {object} ErrorResponse
// @Router /curriculum/topics [get]
func (cfg *ApiConfig) HandleListCurriculumTopics(w http.ResponseWriter, r *http.Request) {
	dbTopics, err := cfg.DB.ListCurriculumTopics(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not receive topics", err)
	}

	topics := make([]TopicResponse, 0)
	for _, topic := range dbTopics {
		topics = append(topics, databaseTopicToResponse(topic))
	}

	respondWithJSON(w, http.StatusOK, topics)
}
