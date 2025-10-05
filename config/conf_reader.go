package config

import (
	"context"
	"fmt"

	"github.com/fmotalleb/go-tools/config"
	"github.com/fmotalleb/go-tools/decoder"
	"github.com/fmotalleb/go-tools/decoder/hooks"
	"github.com/fmotalleb/go-tools/log"
	"github.com/fmotalleb/go-tools/template"
)

func Parse(ctx context.Context, dst *Config, path string) error {
	ctx = log.WithLogger(ctx, log.Of(ctx).Named("config.Parse"))
	cfg, err := config.ReadAndMergeConfig(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to read and merge configs: %w", err)
	}
	hooks.RegisterHook(template.StringTemplateEvaluate())
	decoder, err := decoder.Build(dst)
	if err != nil {
		return fmt.Errorf("create decoder: %w", err)
	}

	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
