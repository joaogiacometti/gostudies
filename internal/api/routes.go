package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *API) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/signup", api.handleSignup)
			r.Post("/login", api.handleLogin)
			r.With(api.AuthMiddleware).Post("/logout", api.handleLogout)
		})
	})
}
