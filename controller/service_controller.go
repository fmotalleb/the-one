package controller

import (
	"context"
	"errors"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/renderer"
)

func Boot(ctx context.Context, cfg *config.Config) error {
	// Compile Templates
	for _, t := range cfg.Templates {
		if err := renderer.RenderTemplates(&t); err != nil && t.GetIsFatal() {
			return errors.Join(
				EngineBootError,
				err,
			)
		}
	}

	return nil
}
