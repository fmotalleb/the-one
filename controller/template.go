package controller

import (
	"errors"

	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/renderer"
)

func compileTemplates(cfg *config.Config, l *zap.Logger) error {
	for _, t := range cfg.Templates {
		tl := l.With(zap.String("src", t.SourceDirectory))
		tl.Debug(
			"rendering template directory",
		)
		if err := renderer.RenderTemplates(&t); err != nil && t.IsFatal {
			tl.Error(
				"failed to render templates",
				zap.Error(err),
			)
			return errors.Join(
				ErrEngineBoot,
				err,
			)
		}
	}

	if len(cfg.Templates) != 0 {
		l.Debug("finished template rendering")
	} else {
		l.Debug("no template was found, ignoring the step")
	}
	return nil
}
