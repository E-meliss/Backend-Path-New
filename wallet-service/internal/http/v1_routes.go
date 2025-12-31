package http

import (
	nethttp "net/http"

	"github.com/ezgicosar/wallet-service/internal/http/handlers"
	"github.com/ezgicosar/wallet-service/internal/http/middleware"
)

type V1Routes struct {
	deps Deps
}

func NewV1Routes(deps Deps) *V1Routes {
	return &V1Routes{deps: deps}
}

func (v *V1Routes) Register(r *Router) {
	authH := handlers.NewAuthHandler(v.deps)
	usersH := handlers.NewUsersHandler(v.deps)
	txH := handlers.NewTransactionsHandler(v.deps)
	balH := handlers.NewBalancesHandler(v.deps)

	// Auth
	r.Handle(nethttp.MethodPost, "/api/v1/auth/register", nethttp.HandlerFunc(authH.Register))
	r.Handle(nethttp.MethodPost, "/api/v1/auth/login", nethttp.HandlerFunc(authH.Login))
	r.Handle(nethttp.MethodPost, "/api/v1/auth/refresh", nethttp.HandlerFunc(authH.Refresh))

	// Users
	r.Handle(nethttp.MethodGet, "/api/v1/users", middleware.AuthStub()(nethttp.HandlerFunc(usersH.List)))
	r.Handle(nethttp.MethodGet, "/api/v1/users/{id}", middleware.AuthStub()(nethttp.HandlerFunc(usersH.Get)))
	r.Handle(nethttp.MethodPut, "/api/v1/users/{id}", middleware.AuthStub()(nethttp.HandlerFunc(usersH.Update)))
	r.Handle(nethttp.MethodDelete, "/api/v1/users/{id}", middleware.AuthStub()(middleware.RequireRole("admin")(nethttp.HandlerFunc(usersH.Delete))))

	// Transactions
	r.Handle(nethttp.MethodPost, "/api/v1/transactions/credit", middleware.AuthStub()(nethttp.HandlerFunc(txH.Credit)))
	r.Handle(nethttp.MethodPost, "/api/v1/transactions/debit", middleware.AuthStub()(nethttp.HandlerFunc(txH.Debit)))
	r.Handle(nethttp.MethodPost, "/api/v1/transactions/transfer", middleware.AuthStub()(nethttp.HandlerFunc(txH.Transfer)))
	r.Handle(nethttp.MethodGet, "/api/v1/transactions/history", middleware.AuthStub()(nethttp.HandlerFunc(txH.History)))
	r.Handle(nethttp.MethodGet, "/api/v1/transactions/{id}", middleware.AuthStub()(nethttp.HandlerFunc(txH.Get)))

	// Balances
	r.Handle(nethttp.MethodGet, "/api/v1/balances/current", middleware.AuthStub()(nethttp.HandlerFunc(balH.Current)))
	r.Handle(nethttp.MethodGet, "/api/v1/balances/historical", middleware.AuthStub()(nethttp.HandlerFunc(balH.Historical)))
	r.Handle(nethttp.MethodGet, "/api/v1/balances/at-time", middleware.AuthStub()(nethttp.HandlerFunc(balH.AtTime)))
}
