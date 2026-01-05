package handlers

import (
	"net/http"

	apphttp "github.com/E-meliss/wallet-service/internal/http"
	"github.com/E-meliss/wallet-service/internal/http/response"
)

type UsersHandler struct {
	deps apphttp.Deps
}

func NewUsersHandler(deps apphttp.Deps) *UsersHandler {
	return &UsersHandler{deps: deps}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "users list (stub)"})
}

func (h *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "users get (stub)"})
}

func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "users update (stub)"})
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "users delete (stub)"})
}
