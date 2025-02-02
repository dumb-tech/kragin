package registerer

import (
	"context"
	"net/http"
)

type Registerer interface {
	RegisterModifiers(f func(
		name string,
		factoryFunc func(map[string]any) func(any) (any, error),
		appliesToRequest bool,
		appliesToResponse bool,
	))
	RegisterHandlers(f func(
		name string,
		handler func(context.Context, map[string]any, http.Handler) (http.Handler, error),
	))
	RegisterClients(f func(
		name string,
		handler func(context.Context, map[string]any) (http.Handler, error),
	))
}
