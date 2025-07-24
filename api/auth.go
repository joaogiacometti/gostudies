package api

import (
	"net/http"

	"github.com/joaogiacometti/gostudies/constants"
	"github.com/joaogiacometti/gostudies/jsonutils"
)

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), constants.SessionKeyUserId) {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]string{
				"error": "No active session found",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
