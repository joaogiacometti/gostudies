package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joaogiacometti/gostudies/users"
)

type API struct {
	Router       *chi.Mux
	UserHandlers *users.UserHandler
	Sessions     *scs.SessionManager
}
