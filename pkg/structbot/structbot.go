package structbot

import (
	"github.com/tjpnz/structbot/internal/transformer"
	"github.com/tjpnz/structbot/pkg/config"
	"github.com/tjpnz/structbot/pkg/factory"
)

// StructBot implements functionality for registration and retrieval of factories.
type StructBot struct {
	config    *config.Config
	factories map[string]factory.Factory
}

// New initializes StructBot.
func New() *StructBot {
	return &StructBot{
		config: &config.Config{
			Transformers: &transformer.TimeTransformer{},
		},
		factories: make(map[string]factory.Factory),
	}
}

// RegisterFactory registers a new factory function by name.
// It returns a pointer to StructBot to allow for chained calls.
func (sb *StructBot) RegisterFactory(name string, fn func() (interface{}, error)) *StructBot {
	sb.factories[name] = factory.New(sb.config, fn)
	return sb
}

// Factory returns a factory function by name or nil if it hasn't been registered.
func (sb *StructBot) Factory(name string) factory.Factory {
	return sb.factories[name]
}
