package api

import (
	"net/http"

	"github.com/joaogiacometti/gostudies/internal/exceptions"
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
	data, err := jsonutils.DecodeValidJson[requests.LoginRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": err.Error(),
		})
		return
	}

	userId, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
			"error": err.Error(),
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": exceptions.ErrUnexpected,
		})
		return
	}

	api.Sessions.Put(r.Context(), "user_id", userId)

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{})
}

func (api *API) handleLogout(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": exceptions.ErrUnexpected,
		})
		return
	}

	api.Sessions.Remove(r.Context(), "user_id")

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{})
}
