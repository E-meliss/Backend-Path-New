package middleware

import (
	"log/slog"
	"net/http"

	"github.com/E-meliss/wallet-service/internal/http/response"
)

func Recover(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic", "recover", rec)
					response.Error(w, r, http.StatusInternalServerError, "internal_error", "unexpected error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
