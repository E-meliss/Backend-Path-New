package handlers

import (
	"net/http"

	apphttp "github.com/E-meliss/wallet-service/internal/http"
	"github.com/E-meliss/wallet-service/internal/http/response"
)

type AuthHandler struct {
	deps apphttp.Deps
}

func NewAuthHandler(deps apphttp.Deps) *AuthHandler {
	return &AuthHandler{deps: deps}
}

type registerReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"message": "register ok (stub)"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "login ok (stub)"})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "refresh ok (stub)"})
}
