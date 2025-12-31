package middleware

import (
	"context"
	"net/http"
)

type userKey struct{}
type AuthUser struct {
	ID   int64
	Role string
}

func AuthStub() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := AuthUser{ID: 1, Role: "admin"}
			ctx := context.WithValue(r.Context(), userKey{}, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromCtx(ctx context.Context) (AuthUser, bool) {
	u, ok := ctx.Value(userKey{}).(AuthUser)
	return u, ok
}
