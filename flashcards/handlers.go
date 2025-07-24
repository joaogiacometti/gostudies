package flashcards

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joaogiacometti/gocards/constants"
	"github.com/joaogiacometti/gocards/jsonutils"
	"github.com/joaogiacometti/gocards/pgstore"
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

func (h *FlashcardHandler) HandleGetFlashcards(w http.ResponseWriter, r *http.Request) {
	userID := h.sessions.Get(r.Context(), constants.SessionKeyUserId).(uuid.UUID)
	if userID == uuid.Nil {
		jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{})
		return
	}

	flashcards, err := h.service.GetFlashcards(r.Context(), userID)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, flashcards)
}

func (h *FlashcardHandler) HandleGetFlashcardByID(w http.ResponseWriter, r *http.Request) {
	flashcardID, err := uuid.Parse(chi.URLParam(r, "flashcardID"))
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "Invalid flashcard ID",
		})
		return
	}

	userID := h.sessions.Get(r.Context(), constants.SessionKeyUserId).(uuid.UUID)
	if userID == uuid.Nil {
		jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{})
		return
	}
	flashcard, err := h.service.GetFlashcardByID(r.Context(), flashcardID, userID)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if flashcard == (pgstore.GetFlashcardByIDRow{}) {
		jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
			"error": "Flashcard not found",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, flashcard)
}

func (h *FlashcardHandler) HandleGetNextFlashcardToReview(w http.ResponseWriter, r *http.Request) {
	userID := h.sessions.Get(r.Context(), constants.SessionKeyUserId).(uuid.UUID)
	if userID == uuid.Nil {
		jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{})
		return
	}

	flashcard, err := h.service.queries.GetNextFlashcardToReview(r.Context(), userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			jsonutils.EncodeJson(w, r, http.StatusNoContent, map[string]any{})
			return
		}

		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, flashcard)
}
