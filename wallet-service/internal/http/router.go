package http

import (
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type route struct {
	method  string
	pattern string
	handler http.Handler
}

type Router struct {
	routes     []route
	middleware []Middleware
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Use(mw ...Middleware) {
	r.middleware = append(r.middleware, mw...)
}

func (r *Router) Handle(method, pattern string, h http.Handler) {
	r.routes = append(r.routes, route{method: method, pattern: pattern, handler: h})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, rt := range r.routes {
		if rt.method != req.Method {
			continue
		}
		params, ok := match(rt.pattern, req.URL.Path)
		if !ok {
			continue
		}
		ctx := withParams(req.Context(), params)
		h := rt.handler
		// global middleware chain
		for i := len(r.middleware) - 1; i >= 0; i-- {
			h = r.middleware[i](h)
		}
		h.ServeHTTP(w, req.WithContext(ctx))
		return
	}
	http.NotFound(w, req)
}

func match(pattern, path string) (map[string]string, bool) {
	pSeg := splitPath(pattern)
	uSeg := splitPath(path)
	if len(pSeg) != len(uSeg) {
		return nil, false
	}
	params := make(map[string]string)
	for i := 0; i < len(pSeg); i++ {
		ps := pSeg[i]
		us := uSeg[i]
		if strings.HasPrefix(ps, "{") && strings.HasSuffix(ps, "}") {
			key := strings.TrimSuffix(strings.TrimPrefix(ps, "{"), "}")
			params[key] = us
			continue
		}
		if ps != us {
			return nil, false
		}
	}
	return params, true
}

func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}
