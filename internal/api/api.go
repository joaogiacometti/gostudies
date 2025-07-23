package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joaogiacometti/gostudies/internal/services"
)

type API struct {
	Router      *chi.Mux
	UserService *services.UserService
	Sessions    *scs.SessionManager
}
