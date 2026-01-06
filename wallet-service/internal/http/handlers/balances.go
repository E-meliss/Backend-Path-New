package handlers

import (
	"github.com/E-meliss/wallet-service/internal/http/response"
	"net/http"
)

type BalancesHandler struct {
}

func NewBalancesHandler() *BalancesHandler {
	return &BalancesHandler{}
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
