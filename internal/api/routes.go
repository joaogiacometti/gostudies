package api

import "github.com/go-chi/chi/v5"

func (api *API) BindRoutes() {
	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/signup", api.handleSignup)
			r.Post("/login", api.handleLogin)
			r.Post("/logout", api.handleLogout)
		})
	})
}
