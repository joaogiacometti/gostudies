package api

import (
	"net/http"

	"github.com/joaogiacometti/gostudies/internal/jsonutils"
)

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "user_id") {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]string{
				"error": "No active session found",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
