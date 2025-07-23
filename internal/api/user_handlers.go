package api

import (
	"net/http"

	"github.com/joaogiacometti/gostudies/internal/jsonutils"
	"github.com/joaogiacometti/gostudies/internal/requests"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) handleSignup(w http.ResponseWriter, r *http.Request) {
	data, err := jsonutils.DecodeValidJson[requests.CreateUserRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": err.Error(),
		})
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	userID, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, password_hash)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": userID,
	})
}

func (api *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	panic("NOT IMPLEMENTED YET")
}

func (api *API) handleLogout(w http.ResponseWriter, r *http.Request) {
	panic("NOT IMPLEMENTED YET")
}
