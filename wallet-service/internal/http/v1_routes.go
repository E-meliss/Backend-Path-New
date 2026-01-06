package http

import (
	nethttp "net/http"
)

type AuthHandler interface {
	Register(nethttp.ResponseWriter, *nethttp.Request)
	Login(nethttp.ResponseWriter, *nethttp.Request)
	Refresh(nethttp.ResponseWriter, *nethttp.Request)
}

type UsersHandler interface {
	List(nethttp.ResponseWriter, *nethttp.Request)
	Get(nethttp.ResponseWriter, *nethttp.Request)
	Update(nethttp.ResponseWriter, *nethttp.Request)
	Delete(nethttp.ResponseWriter, *nethttp.Request)
}

type TransactionsHandler interface {
	Credit(nethttp.ResponseWriter, *nethttp.Request)
	Debit(nethttp.ResponseWriter, *nethttp.Request)
	Transfer(nethttp.ResponseWriter, *nethttp.Request)
	History(nethttp.ResponseWriter, *nethttp.Request)
	Get(nethttp.ResponseWriter, *nethttp.Request)
}

type BalancesHandler interface {
	Current(nethttp.ResponseWriter, *nethttp.Request)
	Historical(nethttp.ResponseWriter, *nethttp.Request)
	AtTime(nethttp.ResponseWriter, *nethttp.Request)
}

type V1Routes struct {
	deps Deps
}

func NewV1Routes(deps Deps) *V1Routes {
	return &V1Routes{deps: deps}
}

func (v *V1Routes) Register(r *Router,
	authH AuthHandler,
	usersH UsersHandler,
	txH TransactionsHandler,
	balH BalancesHandler,
	authMW func(nethttp.Handler) nethttp.Handler,
	requireRoleGen func(string) func(nethttp.Handler) nethttp.Handler,
) {
	// Auth
	r.Handle(nethttp.MethodPost, "/api/v1/auth/register", nethttp.HandlerFunc(authH.Register))
	r.Handle(nethttp.MethodPost, "/api/v1/auth/login", nethttp.HandlerFunc(authH.Login))
	r.Handle(nethttp.MethodPost, "/api/v1/auth/refresh", nethttp.HandlerFunc(authH.Refresh))

	// Users
	r.Handle(nethttp.MethodGet, "/api/v1/users", authMW(nethttp.HandlerFunc(usersH.List)))
	r.Handle(nethttp.MethodGet, "/api/v1/users/{id}", authMW(nethttp.HandlerFunc(usersH.Get)))
	r.Handle(nethttp.MethodPut, "/api/v1/users/{id}", authMW(nethttp.HandlerFunc(usersH.Update)))
	r.Handle(nethttp.MethodDelete, "/api/v1/users/{id}", authMW(requireRoleGen("admin")(nethttp.HandlerFunc(usersH.Delete))))

	// Transactions
	r.Handle(nethttp.MethodPost, "/api/v1/transactions/credit", authMW(nethttp.HandlerFunc(txH.Credit)))
	r.Handle(nethttp.MethodPost, "/api/v1/transactions/debit", authMW(nethttp.HandlerFunc(txH.Debit)))
	r.Handle(nethttp.MethodPost, "/api/v1/transactions/transfer", authMW(nethttp.HandlerFunc(txH.Transfer)))
	r.Handle(nethttp.MethodGet, "/api/v1/transactions/history", authMW(nethttp.HandlerFunc(txH.History)))
	r.Handle(nethttp.MethodGet, "/api/v1/transactions/{id}", authMW(nethttp.HandlerFunc(txH.Get)))

	// Balances
	r.Handle(nethttp.MethodGet, "/api/v1/balances/current", authMW(nethttp.HandlerFunc(balH.Current)))
	r.Handle(nethttp.MethodGet, "/api/v1/balances/historical", authMW(nethttp.HandlerFunc(balH.Historical)))
	r.Handle(nethttp.MethodGet, "/api/v1/balances/at-time", authMW(nethttp.HandlerFunc(balH.AtTime)))
}
