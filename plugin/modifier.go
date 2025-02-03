package plugin

import (
	"strings"

	"github.com/dumb-tech/kragin/registerer"
)

type ModifierPlugin struct {
	name    string
	version string

	modifier registerer.ModifierRegisterer
}

func NewModifier(name string, version string) *ModifierPlugin {
	plug := &ModifierPlugin{
		name:    name,
		version: version,
	}

	plug.modifier = registerer.NewModifierRegisterer(plug.name)

	return plug
}

func (plug *ModifierPlugin) Name() string    { return plug.name }
func (plug *ModifierPlugin) Version() string { return plug.version }

func (plug *ModifierPlugin) ModifiersRegisterer() registerer.ModifierRegisterer {
	return plug.modifier
}

func (plug *ModifierPlugin) RequestHandler(f func() registerer.ModifierFunc) {
	plug.modifier.AddModifier("request", f())
}

func (plug *ModifierPlugin) ResponseHandler(f func() registerer.ModifierFunc) {
	plug.modifier.AddModifier("response", f())
}

func (plug *ModifierPlugin) Debug() string {
	sb := strings.Builder{}

	sb.WriteString("Name: " + plug.name + "\n")
	sb.WriteString("Version: " + plug.version + "\n")

	return sb.String()
}
