package handlers

import (
	"net/http"

	apphttp "github.com/E-meliss/wallet-service/internal/http"
	"github.com/E-meliss/wallet-service/internal/http/response"
)

type TransactionsHandler struct {
	deps apphttp.Deps
}

func NewTransactionsHandler(deps apphttp.Deps) *TransactionsHandler {
	return &TransactionsHandler{deps: deps}
}

func (h *TransactionsHandler) Credit(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "credit (stub)"})
}

func (h *TransactionsHandler) Debit(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "debit (stub)"})
}

func (h *TransactionsHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "transfer (stub)"})
}

func (h *TransactionsHandler) History(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "history (stub)"})
}

func (h *TransactionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"message": "transaction get (stub)"})
}
