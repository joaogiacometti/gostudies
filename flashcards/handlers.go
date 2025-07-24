package flashcards

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/joaogiacometti/gostudies/constants"
	"github.com/joaogiacometti/gostudies/jsonutils"
	"github.com/joaogiacometti/gostudies/pgstore"
)

type FlashcardHandler struct {
	service  *FlashcardService
	sessions *scs.SessionManager
}

func NewFlashcardHandler(service *FlashcardService, sessions *scs.SessionManager) *FlashcardHandler {
	return &FlashcardHandler{
		service:  service,
		sessions: sessions,
	}
}

func (h *FlashcardHandler) HandleCreateFlashcard(w http.ResponseWriter, r *http.Request) {
	data, err := jsonutils.DecodeValidJson[RequestCreateFlashcard](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": err.Error(),
		})
		return
	}

	userID := h.sessions.Get(r.Context(), constants.SessionKeyUserId).(uuid.UUID)
	if userID == uuid.Nil {
		jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{})
		return
	}

	isDuplicated, err := h.service.queries.IsDuplicateFlashcard(r.Context(), pgstore.IsDuplicateFlashcardParams{
		UserID:   userID,
		Question: data.Question,
	})
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}
	if isDuplicated {
		jsonutils.EncodeJson(w, r, http.StatusConflict, map[string]any{
			"error": "Flashcard with the same question already exists",
		})
		return
	}

	flashcardID, err := h.service.CreateFlashcard(r.Context(), userID, data.Question, data.Answer)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"flashcardID": flashcardID,
	})
}
