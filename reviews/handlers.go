package reviews

import (
	"math"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joaogiacometti/gocards/constants"
	"github.com/joaogiacometti/gocards/flashcards"
	"github.com/joaogiacometti/gocards/jsonutils"
)

type ReviewHandler struct {
	service          *ReviewService
	flashCardService *flashcards.FlashcardService
	sessions         *scs.SessionManager
}

func NewReviewHandler(service *ReviewService, flashCardService *flashcards.FlashcardService, sessions *scs.SessionManager) *ReviewHandler {
	return &ReviewHandler{
		service:          service,
		flashCardService: flashCardService,
		sessions:         sessions,
	}
}

func (h *ReviewHandler) HandleReviewFlashcard(w http.ResponseWriter, r *http.Request) {
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

	data, err := jsonutils.DecodeValidJson[RequestReviewFlashcard](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": err.Error(),
		})
		return
	}

	flashcard, err := h.flashCardService.GetFlashcardByID(r.Context(), flashcardID, userID)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	interval := calculateNextReview(*data.Remembered, flashcard.SuccessCount.Int32)
	successCount := calculateSuccessCount(*data.Remembered, flashcard.SuccessCount.Int32)

	err = h.service.ReviewFlashcard(r.Context(), flashcardID, userID, interval, successCount, *data.Remembered)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{})
}

func calculateNextReview(remembered bool, successCount int32) time.Time {
	now := time.Now()

	if !remembered {
		return now
	}

	if successCount <= 0 {
		return now.Add(time.Hour * 24)
	}
	days := math.Pow(2, float64(successCount))
	return now.Add(time.Hour * 24 * time.Duration(days))
}

func calculateSuccessCount(remembered bool, successCount int32) int32 {
	if remembered {
		return successCount + 1
	}
	return 0
}
