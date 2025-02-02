package registerer

import (
	"context"
	"net/http"
)

type ModifierFunc = func(map[string]any) func(any) (any, error)

type modifiersMap map[string]func(map[string]any) func(any) (any, error)

type ModifierRegisterer struct {
	pluginName string
	modifiers  modifiersMap
}

func NewModifierRegisterer(pluginName string) *ModifierRegisterer {
	reg := &ModifierRegisterer{pluginName: pluginName, modifiers: make(modifiersMap)}
	return reg
}

func (reg *ModifierRegisterer) AddModifier(name string, f ModifierFunc) {
	reg.modifiers[name] = f
}

func (reg *ModifierRegisterer) RegisterModifiers(f func(
	name string,
	factoryFunc func(map[string]any) func(any) (any, error),
	appliesToRequest bool,
	appliesToResponse bool,
)) {
	request, ok := reg.modifiers["request"]
	if !ok {
		f(reg.pluginName+"-request", request, true, false)
		delete(reg.modifiers, "request")
	}
	response, ok := reg.modifiers["response"]
	if !ok {
		f(reg.pluginName+"-response", response, false, true)
		delete(reg.modifiers, "response")
	}
	for name, modifier := range reg.modifiers {
		f(reg.pluginName+"-"+name, modifier, true, true)
	}
}

// Stubs

func (reg *ModifierRegisterer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]any, http.Handler) (http.Handler, error),
)) {
}

func (reg *ModifierRegisterer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]any) (http.Handler, error),
)) {
}
