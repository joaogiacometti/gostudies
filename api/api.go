package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joaogiacometti/gocards/flashcards"
	"github.com/joaogiacometti/gocards/reviews"
	"github.com/joaogiacometti/gocards/users"
)

type API struct {
	Router            *chi.Mux
	UserHandlers      *users.UserHandler
	FlashcardHandlers *flashcards.FlashcardHandler
	ReviewHandlers    *reviews.ReviewHandler
	Sessions          *scs.SessionManager
}
