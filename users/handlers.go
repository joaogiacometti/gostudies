package users

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/joaogiacometti/gostudies/constants"
	"github.com/joaogiacometti/gostudies/exceptions"
	"github.com/joaogiacometti/gostudies/jsonutils"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService *UserService
	Sessions    *scs.SessionManager
}

func NewUserHandler(userService *UserService, sessions *scs.SessionManager) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Sessions:    sessions,
	}
}

func (h *UserHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	data, err := jsonutils.DecodeValidJson[CreateUserRequest](r)
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

	userID, err := h.UserService.CreateUser(r.Context(), data.UserName, data.Email, password_hash)
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

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data, err := jsonutils.DecodeValidJson[LoginRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": err.Error(),
		})
		return
	}

	userId, err := h.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
			"error": err.Error(),
		})
		return
	}

	err = h.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": exceptions.ErrUnexpected,
		})
		return
	}

	h.Sessions.Put(r.Context(), constants.SessionKeyUserId, userId)

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{})
}

func (h *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := h.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": exceptions.ErrUnexpected,
		})
		return
	}

	h.Sessions.Remove(r.Context(), constants.SessionKeyUserId)

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{})
}
