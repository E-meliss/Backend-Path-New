package handlers

import (
	"net/http"

	apphttp "github.com/E-meliss/wallet-service/internal/http"
	"github.com/E-meliss/wallet-service/internal/http/response"
)

type BalancesHandler struct {
	deps apphttp.Deps
}

func NewBalancesHandler(deps apphttp.Deps) *BalancesHandler {
	return &BalancesHandler{deps: deps}
}

func (h *BalancesHandler) Current(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "balance current (stub)"})
}

func (h *BalancesHandler) Historical(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "balance historical (stub)"})
}

func (h *BalancesHandler) AtTime(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "balance at-time (stub)"})
}
