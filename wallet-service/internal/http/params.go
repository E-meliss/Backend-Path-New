package http

import "context"

type paramsKey struct{}

func withParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, paramsKey{}, params)
}

func Param(ctx context.Context, key string) string {
	m, _ := ctx.Value(paramsKey{}).(map[string]string)
	if m == nil {
		return ""
	}
	return m[key]
}
